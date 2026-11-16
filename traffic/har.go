package traffic

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
)

// HARImportResult describes recorded data ingestion, not target security.
type HARImportResult struct {
	Entries int `json:"entries"`
	Imported int `json:"imported"`
	Duplicates int `json:"duplicates"`
	Skipped int `json:"skipped"`
	Failed int `json:"failed"`
	Unprocessed int `json:"unprocessed"`
	Reasons map[string]int `json:"reasons"`
}

// HARPageAttribution is file-local provenance, not a browser session or identity.
// Ordinal is one-based in log.pages and set only for a unique reference.
type HARPageAttribution struct {
	State string `json:"state"`
	Ordinal int `json:"ordinal,omitempty"`
}

func (p HARPageAttribution) valid() bool {
	switch p.State {
	case "unique":
		return p.Ordinal > 0
	case "missing", "unknown", "ambiguous":
		return p.Ordinal == 0
	default:
		return false
	}
}

// Duplicate page IDs are ambiguous. Malformed optional metadata never prevents
// importing the independently valid HTTP exchange.
func harPageIndex(raw json.RawMessage) map[string]int {
	var pages []json.RawMessage
	index := map[string]int{}
	if json.Unmarshal(raw, &pages) != nil { return index }
	for i, item := range pages {
		var page struct { ID string `json:"id"` }
		if json.Unmarshal(item, &page) != nil || strings.TrimSpace(page.ID) == "" { continue }
		if _, exists := index[page.ID]; exists { index[page.ID] = -1 } else { index[page.ID] = i+1 }
	}
	return index
}

func harPageReference(raw json.RawMessage, index map[string]int) HARPageAttribution {
	if len(raw) == 0 || string(raw) == "null" { return HARPageAttribution{State:"missing"} }
	var id string
	if json.Unmarshal(raw, &id) != nil { return HARPageAttribution{State:"unknown"} }
	if id == "" { return HARPageAttribution{State:"missing"} }
	switch ordinal := index[id]; {
	case ordinal > 0: return HARPageAttribution{State:"unique", Ordinal:ordinal}
	case ordinal < 0: return HARPageAttribution{State:"ambiguous"}
	default: return HARPageAttribution{State:"unknown"}
	}
}

type harHeader struct {
	Name string `json:"name"`
	Value string `json:"value"`
}
type harContent struct {
	Text *string `json:"text"`
	Size *int64 `json:"size"`
	Encoding string `json:"encoding"`
	MIMEType string `json:"mimeType"`
}
type harEntry struct {
	StartedDateTime string `json:"startedDateTime"`
	PageRef json.RawMessage `json:"pageref"`
	Request *struct {
		Method string `json:"method"`
		URL string `json:"url"`
		HTTPVersion string `json:"httpVersion"`
		Headers []harHeader `json:"headers"`
		BodySize *int64 `json:"bodySize"`
		PostData *harContent `json:"postData"`
	} `json:"request"`
	Response *struct {
		Status int `json:"status"`
		HTTPVersion string `json:"httpVersion"`
		Headers []harHeader `json:"headers"`
		Content *harContent `json:"content"`
	} `json:"response"`
}
type harCapture struct {
	url, host, method, path string
	reqHead, respHead string
	reqType, respType string
	request, response []byte
	status int
	ts int64
	pageRef json.RawMessage
	page HARPageAttribution
}

// ImportHAR consumes saved records only. It never requests recorded URLs.
// An exact file digest and entry position provide idempotent import IDs.
func (t *Traffic) ImportHAR(ctx context.Context, raw []byte) (HARImportResult, error) {
	var doc struct {
		Log *struct {
			Entries []json.RawMessage `json:"entries"`
			Pages json.RawMessage `json:"pages"`
		} `json:"log"`
	}
	if err := json.Unmarshal(bytes.TrimPrefix(raw, []byte{0xef, 0xbb, 0xbf}), &doc); err != nil || doc.Log == nil || doc.Log.Entries == nil {
		return HARImportResult{}, errors.New("invalid HAR: log.entries must be an array")
	}
	result := HARImportResult{Entries:len(doc.Log.Entries), Reasons:map[string]int{}}
	digest := sha256.Sum256(raw)
	pages := harPageIndex(doc.Log.Pages)
	for i, item := range doc.Log.Entries {
		if ctx.Err() != nil {
			result.Unprocessed = result.Entries-i
			break
		}
		c, reason := parseHARCapture(item)
		if reason != "" {
			result.Skipped++
			result.Reasons[reason]++
			continue
		}
		c.page = harPageReference(c.pageRef, pages)
		duplicate, err := t.storeHARCapture(fmt.Sprintf("har-%x-%d", digest, i), c)
		switch {
		case err != nil:
			result.Failed++
			result.Reasons["storage_error"]++
		case duplicate:
			result.Duplicates++
		default:
			result.Imported++
		}
	}
	return result,nil
}

func parseHARCapture(raw json.RawMessage) (harCapture,string) {
	var e harEntry
	if err:=json.Unmarshal(raw,&e); err!=nil || e.Request==nil {
		return harCapture{},"invalid_entry"
	}
	r:=e.Request
	u,err:=url.Parse(r.URL)
	if err!=nil || u.Hostname()=="" || (u.Scheme!="http" && u.Scheme!="https") || !harToken(r.Method) {
		return harCapture{},"invalid_request"
	}
	at,err:=time.Parse(time.RFC3339Nano,e.StartedDateTime)
	if err!=nil { return harCapture{},"unknown_capture_time" }
	if e.Response==nil || e.Response.Status==0 { return harCapture{},"unobserved_status" }
	p:=e.Response
	if p.Status<100 || p.Status>599 { return harCapture{},"invalid_status" }
	reqHeaders,ok:=harHeaderLines(r.Headers)
	if !ok { return harCapture{},"invalid_headers" }
	respHeaders,ok:=harHeaderLines(p.Headers)
	if !ok || !harLine(r.HTTPVersion) || !harLine(p.HTTPVersion) { return harCapture{},"invalid_headers" }

	var reqBody []byte
	if r.PostData!=nil && r.PostData.Text!=nil {
		if r.PostData.Encoding!="" { return harCapture{},"unsupported_encoding" }
		reqBody=[]byte(*r.PostData.Text)
		if r.BodySize!=nil && *r.BodySize>=0 && int64(len(reqBody))!=*r.BodySize {
			return harCapture{},"request_size_mismatch"
		}
	} else if r.PostData!=nil || r.BodySize==nil || *r.BodySize!=0 {
		return harCapture{},"request_body_missing"
	}
	c:=p.Content
	if c==nil { return harCapture{},"response_body_missing" }
	var respBody []byte
	if c.Text==nil {
		if c.Size==nil || *c.Size!=0 { return harCapture{},"response_body_missing" }
	} else {
		switch strings.ToLower(c.Encoding) {
		case "":
			respBody=[]byte(*c.Text)
		case "base64":
			respBody,err=base64.StdEncoding.DecodeString(*c.Text)
			if err!=nil { return harCapture{},"invalid_encoding" }
		default:
			return harCapture{},"unsupported_encoding"
		}
	}
	if c.Size!=nil && *c.Size>=0 && int64(len(respBody))!=*c.Size { return harCapture{},"response_size_mismatch" }
	reqVersion,respVersion:=r.HTTPVersion,p.HTTPVersion
	if reqVersion=="" { reqVersion="[HTTP version unavailable]" }
	if respVersion=="" { respVersion="HTTP" }
	reqType,respType:=harContentType(r.Headers),harContentType(p.Headers)
	if reqType=="" && r.PostData!=nil { reqType=r.PostData.MIMEType }
	if respType=="" { respType=c.MIMEType }
	return harCapture{
		url:r.URL,host:u.Hostname(),method:r.Method,path:db.TemplatePath(u.EscapedPath()),
		reqHead:fmt.Sprintf("%s %s %s\n%s",r.Method,u.RequestURI(),reqVersion,reqHeaders),
		respHead:fmt.Sprintf("%s %d\n%s",respVersion,p.Status,respHeaders),
		reqType:reqType,respType:respType,request:reqBody,response:respBody,status:p.Status,ts:at.Unix(),pageRef:e.PageRef,
	},""
}

// Keep the recorded header array, including duplicate names and order.
func harHeaderLines(headers []harHeader) (string,bool) {
	if headers==nil { return "",false }
	var b strings.Builder
	for _,h:=range headers {
		if !harToken(strings.TrimPrefix(h.Name,":")) || !harLine(h.Value) { return "",false }
		b.WriteString(h.Name)
		b.WriteString(": ")
		b.WriteString(h.Value)
		b.WriteByte('\n')
	}
	return b.String(),true
}
func harToken(s string) bool {
	if s=="" { return false }
	for _,r:=range s {
		if (r>='a' && r<='z') || (r>='A' && r<='Z') || (r>='0' && r<='9') || strings.ContainsRune("!#$%&'*+-.^_~|"+string(rune(96)),r) { continue }
		return false
	}
	return true
}
func harLine(s string) bool {
	for _,r:=range s { if (r<32 && r!='\t') || r==127 { return false } }
	return true
}
func harContentType(headers []harHeader) string {
	for _,h:=range headers { if strings.EqualFold(h.Name,"content-type") { return h.Value } }
	return ""
}
func (t *Traffic) storeHARCapture(id string,c harCapture) (bool,error) {
	t.wmu.Lock()
	defer t.wmu.Unlock()
	var exists bool
	if err:=t.db.QueryRow("SELECT EXISTS(SELECT 1 FROM exchanges WHERE id=?)",id).Scan(&exists); err!=nil || exists { return exists,err }
	req:=t.spill(c.request,c.reqType)
	resp:=t.spill(c.response,c.respType)
	// Unlike live capture, an import must not silently retain only a preview.
	if (len(c.request)>maxInlineBody && req.hash=="") || (len(c.response)>maxInlineBody && resp.hash=="") {
		return false,errors.New("HAR body storage incomplete")
	}
	return false,t.insertExchangeWithPage(id,c.ts,c.host,c.method,c.path,c.url,c.status,c.respType,len(c.request),len(c.response),c.reqHead,req,c.respHead,resp,&c.page)
}
