package traffic

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
)

// ExternalExchange is one captured request/response pair handed in from
// outside the MITM proxy (currently the browser-extension ingest endpoint).
// Header maps are lowercased-name oriented on the extension side but any
// casing is accepted; bodies are raw text (binary capture is skipped by the
// extension, not handled here).
type ExternalExchange struct {
	CaptureID string `json:"captureId,omitempty"`
	ExchangeID string `json:"exchangeId,omitempty"` // source request ID, scoped by capture ID
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	Status          int               `json:"status"`
	RequestHeaders  map[string]string `json:"requestHeaders"`
	RequestBody     string            `json:"requestBody"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	ResponseBody    string            `json:"responseBody"`
}

// RecordExternal stores an externally captured exchange through the same
// spill/blob/index path as proxy traffic, so ingested entries are
// indistinguishable from proxy-captured ones on the Traffic page and in
// full-text search.
func (t *Traffic) RecordExternal(e ExternalExchange) error {
	_, err := t.RecordExternalWithReceipt(e)
	return err
}

// RecordExternalWithReceipt acknowledges a previously persisted browser request
// without overwriting its bodies, status or attribution after an upload retry.
func (t *Traffic) RecordExternalWithReceipt(e ExternalExchange) (duplicate bool, err error) {
	u, err := url.Parse(e.URL)
	if err != nil || u.Host == "" {
		return false, fmt.Errorf("invalid url %q", e.URL)
	}
	if len(e.CaptureID) > 128 { return false, fmt.Errorf("capture ID too long") }
	method := strings.ToUpper(strings.TrimSpace(e.Method))
	if method == "" {
		return false, fmt.Errorf("empty method for %q", e.URL)
	}
	// Status 0 means the external recorder did not observe an HTTP status.
	// Preserve that uncertainty rather than fabricating a successful response.
	status := e.Status

	host := u.Hostname()
	reqHead := fmt.Sprintf("%s %s %s\n%s", method, u.RequestURI(), "HTTP/1.1", headerLines(toHeaderMap(e.RequestHeaders)))
	respHead := fmt.Sprintf("HTTP %d\n%s", status, headerLines(toHeaderMap(e.ResponseHeaders)))
	if status == 0 {
		respHead = "HTTP status unavailable\n" + headerLines(toHeaderMap(e.ResponseHeaders))
	}
	respCT := e.ResponseHeaders["content-type"]

	t.wmu.Lock()
	defer t.wmu.Unlock()

	// Network.requestId is stable within a browser capture. Replayed batches
	// must not overwrite a previously committed exchange when the HTTP
	// acknowledgement was lost. Old clients without IDs retain old behavior.
	id := ""
	if e.CaptureID != "" && e.ExchangeID != "" {
		sum := sha256.Sum256([]byte("artex-browser-v1\x00" + e.CaptureID + "\x00" + e.ExchangeID))
		id = "ext-" + hex.EncodeToString(sum[:])
		var present int
		err := t.db.QueryRow(`SELECT 1 FROM exchanges WHERE id=?`, id).Scan(&present)
		if err == nil { return true, nil }
		if !errors.Is(err, sql.ErrNoRows) { return false, fmt.Errorf("read captured exchange id: %w", err) }
	}

	// Bodies are spilled before the transaction opens, matching record().
	reqB := t.spill([]byte(e.RequestBody), e.RequestHeaders["content-type"])
	respB := t.spill([]byte(e.ResponseBody), respCT)

	now := time.Now()
	if id == "" {
		n := t.seq.Add(1)
		id = fmt.Sprintf("%d-%04d", now.Unix(), n%10000)
	}

	if err := t.insertExchangeWithContext(id, now.Unix(), host, method, db.TemplatePath(u.EscapedPath()), u.String(), status, respCT, len(e.RequestBody), len(e.ResponseBody), reqHead, reqB, respHead, respB, nil, e.CaptureID); err != nil {
		log.Printf("[traffic] failed to ingest %s: %v", u.String(), err)
		return false, err
	}
	return false, nil
}

// insertExchange writes the index row, the body row, blob references and the
// FTS entry inside one transaction. record() and RecordExternal share it so
// both capture paths stay byte-identical in storage layout.
func (t *Traffic) insertExchange(id string, ts int64, host, method, urlTemplate, rawURL string, status int, contentType string, reqLen, respLen int, reqHead string, reqB storedBody, respHead string, respB storedBody) error {
	return t.insertExchangeWithPage(id, ts, host, method, urlTemplate, rawURL, status, contentType, reqLen, respLen, reqHead, reqB, respHead, respB, nil)
}

func (t *Traffic) insertExchangeWithPage(id string, ts int64, host, method, urlTemplate, rawURL string, status int, contentType string, reqLen, respLen int, reqHead string, reqB storedBody, respHead string, respB storedBody, page *HARPageAttribution) error {
	return t.insertExchangeWithContext(id,ts,host,method,urlTemplate,rawURL,status,contentType,reqLen,respLen,reqHead,reqB,respHead,respB,page,"")
}

func (t *Traffic) insertExchangeWithContext(id string, ts int64, host, method, urlTemplate, rawURL string, status int, contentType string, reqLen, respLen int, reqHead string, reqB storedBody, respHead string, respB storedBody, page *HARPageAttribution, captureID string) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // no-op once committed

	res, err := tx.Exec(`INSERT OR REPLACE INTO exchanges(id,ts,host,method,url_template,url,status,content_type,req_len,resp_len,path)
VALUES(?,?,?,?,?,?,?,?,?,?,'')`,
		id, ts, host, method, urlTemplate, rawURL, status, contentType, reqLen, respLen)
	if err != nil {
		return fmt.Errorf("writing index: %w", err)
	}
	rowid, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fetching rowid: %w", err)
	}

	if _, err := tx.Exec(`INSERT OR REPLACE INTO exchange_bodies(id,req_head,req_body,req_blob,resp_head,resp_body,resp_blob)
VALUES(?,?,?,?,?,?,?)`,
		id, reqHead, reqB.inline, nullIfEmpty(reqB.hash), respHead, respB.inline, nullIfEmpty(respB.hash)); err != nil {
		return fmt.Errorf("writing body: %w", err)
	}

	for _, h := range []string{reqB.hash, respB.hash} {
		if h == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO blob_refs(hash,exchange_id) VALUES(?,?)`, h, id); err != nil {
			return fmt.Errorf("registering blob reference: %w", err)
		}
	}

	if t.fts {
		idx := strings.Join([]string{rawURL, reqHead, reqB.index, respHead, respB.index}, "\n")
		if _, err := tx.Exec(`INSERT INTO ex_fts(rowid,content) VALUES(?,?)`, rowid, idx); err != nil {
			return fmt.Errorf("writing full-text index: %w", err)
		}
	}

	if page != nil {
		if !page.valid() { return fmt.Errorf("invalid HAR page attribution") }
		if _, err := tx.Exec(`INSERT INTO har_page_attribution(exchange_id,state,ordinal) VALUES(?,?,?)`, id, page.State, page.Ordinal); err != nil {
			return fmt.Errorf("writing HAR page attribution: %w", err)
		}
	}

	if captureID != "" {
		if _,err:=tx.Exec(`INSERT INTO capture_attribution(exchange_id,capture_id) VALUES(?,?)`,id,captureID);err!=nil{
			return fmt.Errorf("writing capture attribution: %w",err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing: %w", err)
	}
	return nil
}

func toHeaderMap(flat map[string]string) map[string][]string {
	m := make(map[string][]string, len(flat))
	for k, v := range flat {
		m[k] = []string{v}
	}
	return m
}
