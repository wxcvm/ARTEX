package traffic

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"sort"
	"strings"
)

// CompareRecorded reports differences between stored representations. It does
// not replay traffic, infer authorization, or classify vulnerability outcomes.
type RecordedComparison struct {
	Baseline RecordedMeta `json:"baseline"`
	Candidate RecordedMeta `json:"candidate"`
	SameURL bool `json:"same_url"`
	StatusRelation string `json:"status_relation"`
	RequestHeaders NameChanges `json:"request_headers"`
	ResponseHeaders NameChanges `json:"response_headers"`
	RequestBody BodyComparison `json:"request_body"`
	ResponseBody BodyComparison `json:"response_body"`
}

type RecordedMeta struct {
	ID string `json:"id"`
	Method string `json:"method"`
	TS int64 `json:"ts"`
	Status int `json:"status"`
	ContentType string `json:"content_type"`
}

type NameChanges struct {
	Available bool `json:"available"`
	Added []string `json:"added"`
	Removed []string `json:"removed"`
	Changed []string `json:"changed"`
	AddedCount int `json:"added_count"`
	RemovedCount int `json:"removed_count"`
	ChangedCount int `json:"changed_count"`
}

type BodyComparison struct {
	BaselineBytes int64 `json:"baseline_bytes"`
	CandidateBytes int64 `json:"candidate_bytes"`
	Relation string `json:"relation"` // equal, different, unknown
	JSONFields NameChanges `json:"json_fields"`
}

const compareJSONBudget = 256 << 10
const compareNamePreview = 50

type recordedBody struct {
	size int64
	digest string
	valid bool
	object map[string]any
}
type recordedSide struct {
	meta RecordedMeta
	url string
	reqHead, respHead map[string][]string
	request, response recordedBody
}

func (t *Traffic) CompareRecorded(ctx context.Context, baseline, candidate string) (RecordedComparison, error) {
	if baseline == "" || candidate == "" || baseline == candidate {
		return RecordedComparison{}, errors.New("two distinct traffic IDs are required")
	}
	sides := make([]recordedSide, 0, 2)
	err := t.ReadEvidence(ctx, []string{baseline, candidate}, func(e EvidenceExchange) error {
		req, err := inspectRecordedBody(ctx, e.Request, e.ReqLen, e.ReqHash)
		if err != nil { return err }
		resp, err := inspectRecordedBody(ctx, e.Response, e.RespLen, e.RespHash)
		if err != nil { return err }
		sides = append(sides, recordedSide{
			meta: RecordedMeta{ID:e.ID, Method:e.Method, TS:e.TS, Status:e.Status, ContentType:e.ContentType},
			url:e.URL, reqHead:recordedHeaderMap(e.ReqHead), respHead:recordedHeaderMap(e.RespHead),
			request:req, response:resp,
		})
		return nil
	})
	if err != nil { return RecordedComparison{}, err }
	a,b := sides[0],sides[1]
	status := "unknown"
	if a.meta.Status >= 100 && a.meta.Status <= 599 && b.meta.Status >= 100 && b.meta.Status <= 599 {
		status = "different"
		if a.meta.Status == b.meta.Status { status = "equal" }
	}
	return RecordedComparison{
		Baseline:a.meta, Candidate:b.meta, SameURL:a.url==b.url, StatusRelation:status,
		RequestHeaders:compareNames(a.reqHead,b.reqHead), ResponseHeaders:compareNames(a.respHead,b.respHead),
		RequestBody:compareRecordedBodies(a.request,b.request), ResponseBody:compareRecordedBodies(a.response,b.response),
	},nil
}

// Parse the recorded header representation, not a reconstructed HTTP message.
// Names are case-insensitive; repeated values remain ordered. Values never
// leave the comparison response, which reports changed names only.
func recordedHeaderMap(raw string) map[string][]string {
	raw = strings.ReplaceAll(raw,"\r\n","\n")
	lines := strings.Split(raw,"\n")
	if len(lines)==0 || strings.TrimSpace(lines[0])=="" { return nil }
	out := map[string][]string{}
	for _, line := range lines[1:] {
		if line == "" { break }
		// HTTP/2 pseudo-header names start with a colon.
		offset:=0
		if strings.HasPrefix(line,":") { offset=1 }
		i:=strings.IndexByte(line[offset:],':')
		if i<0 { return nil }
		i+=offset
		name:=strings.ToLower(strings.TrimSpace(line[:i]))
		if name=="" { return nil }
		out[name]=append(out[name],strings.TrimSpace(line[i+1:]))
	}
	return out
}

func compareNames[V any](a,b map[string]V) NameChanges {
	out:=NameChanges{Available:a!=nil && b!=nil, Added:[]string{}, Removed:[]string{}, Changed:[]string{}}
	if !out.Available { return out }
	for k,v:=range a {
		other,ok:=b[k]
		if !ok { out.Removed=append(out.Removed,k) } else if !reflect.DeepEqual(v,other) { out.Changed=append(out.Changed,k) }
	}
	for k:=range b { if _,ok:=a[k]; !ok { out.Added=append(out.Added,k) } }
	sort.Strings(out.Added); sort.Strings(out.Removed); sort.Strings(out.Changed)
	out.AddedCount=len(out.Added); out.RemovedCount=len(out.Removed); out.ChangedCount=len(out.Changed)
	if len(out.Added)>compareNamePreview { out.Added=out.Added[:compareNamePreview] }
	if len(out.Removed)>compareNamePreview { out.Removed=out.Removed[:compareNamePreview] }
	if len(out.Changed)>compareNamePreview { out.Changed=out.Changed[:compareNamePreview] }
	return out
}

type compareContextReader struct { ctx context.Context; r io.Reader }
func (r compareContextReader) Read(p []byte) (int,error) {
	if err:=r.ctx.Err(); err!=nil { return 0,err }
	return r.r.Read(p)
}

func inspectRecordedBody(ctx context.Context, r io.Reader, expected int64, blobHash string) (recordedBody,error) {
	h:=sha256.New()
	reader:=io.TeeReader(compareContextReader{ctx:ctx,r:r},h)
	// Stream the full body for byte comparison; retain only a small prefix for
	// optional JSON analysis. Large bodies are never loaded in full.
	prefix,err:=io.ReadAll(io.LimitReader(reader,compareJSONBudget+1))
	if err!=nil { return recordedBody{},err }
	rest,err:=io.Copy(io.Discard,reader)
	if err!=nil { return recordedBody{},err }
	size:=int64(len(prefix))+rest
	digest:=hex.EncodeToString(h.Sum(nil))
	out:=recordedBody{size:size,digest:digest,valid:expected>=0 && expected==size && (blobHash=="" || blobHash==digest)}
	if !out.valid || size>compareJSONBudget { return out,nil }
	// Compare only parsed top-level object members, with JSON numbers retained
	// as their recorded lexical representation. Never return member values.
	if trimmed:=bytes.TrimSpace(prefix); len(trimmed)>0 && trimmed[0]=='{' && json.Valid(trimmed) {
		dec:=json.NewDecoder(bytes.NewReader(trimmed))
		dec.UseNumber()
		var obj map[string]any
		if dec.Decode(&obj)==nil { out.object=obj }
	}
	return out,nil
}

func compareRecordedBodies(a,b recordedBody) BodyComparison {
	out:=BodyComparison{BaselineBytes:a.size,CandidateBytes:b.size,Relation:"unknown",JSONFields:compareNames(a.object,b.object)}
	if a.valid && b.valid {
		out.Relation="different"
		if a.size==b.size && a.digest==b.digest { out.Relation="equal" }
	}
	return out
}
