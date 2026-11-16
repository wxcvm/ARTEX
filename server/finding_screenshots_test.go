package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestFindingScreenshotAPI(t *testing.T) {
	s, f, req := trafficEvidenceServer(t)
	seedServerEvidenceFlow(t, s, "screenshot-flow", []byte("local fixture"))
	list, err := s.evidenceStore().Bind(t.Context(), f.FindingID, []db.TrafficRef{{TrafficID: "screenshot-flow", Note: "saved step note"}})
	if err != nil {
		t.Fatal(err)
	}
	base := "/api/exploration/findings/" + i64s(f.FindingID) + "/screenshots"
	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	var img bytes.Buffer
	p := image.NewRGBA(image.Rect(0, 0, 4, 4))
	p.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(&img, p); err != nil {
		t.Fatal(err)
	}
	upload := func(raw []byte, binding string) *httptest.ResponseRecorder {
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		file, err := form.CreateFormFile("image", "fixture.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = file.Write(raw); err != nil {
			t.Fatal(err)
		}
		if err = form.WriteField("binding_id", binding); err != nil {
			t.Fatal(err)
		}
		if err = form.Close(); err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRequest("POST", base, &body)
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", form.FormDataContentType())
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		return w
	}
	if w := upload([]byte("<svg/>"), ""); w.Code != 400 {
		t.Fatal("accepted non-raster payload", w.Code)
	}
	w := upload(img.Bytes(), i64s(list.Bindings[0].ID))
	if w.Code != 201 {
		t.Fatal(w.Code, w.Body)
	}
	var screenshot db.FindingScreenshot
	if err := json.Unmarshal(w.Body.Bytes(), &screenshot); err != nil || screenshot.Caption != "saved step note" {
		t.Fatal(screenshot, err)
	}
	if w = req("GET", base, ""); w.Code != 200 || !strings.Contains(w.Body.String(), screenshot.ID) {
		t.Fatal(w.Code, w.Body)
	}
	path := base + "/" + screenshot.ID
	if w = req("GET", path, ""); w.Code != 200 || w.Header().Get("Content-Type") != "image/png" {
		t.Fatal(w.Code, w.Header())
	}
	if _, err := png.Decode(bytes.NewReader(w.Body.Bytes())); err != nil {
		t.Fatal(err)
	}
	if w = req("PATCH", path, `{"caption":"edited note"}`); w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	if w = req("GET", base, ""); !strings.Contains(w.Body.String(), "edited note") {
		t.Fatal(w.Body)
	}
	reproductionPath := fmt.Sprintf("/api/exploration/findings/%d/reproduction", f.FindingID)
	stepsJSON, err := json.Marshal(map[string]any{"steps": []db.FindingReproductionStep{
		{Title: "Recorded fixture step", Observation: "Observed fixture response", BindingID: i64s(list.Bindings[0].ID)},
		{Title: "Recorded follow-up"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if w = req("PUT", reproductionPath, string(stepsJSON)); w.Code != 200 || !strings.Contains(w.Body.String(), "Recorded fixture step") {
		t.Fatal(w.Code, w.Body)
	}
	if w = req("PUT", reproductionPath, `{"steps":[{"title":"invalid binding","binding_id":"999999999999999"}]}`); w.Code == 200 {
		t.Fatal("accepted nonexistent binding")
	}
	if w = req("GET", reproductionPath, ""); w.Code != 200 || !strings.Contains(w.Body.String(), "Recorded fixture step") {
		t.Fatal("failed replacement changed saved steps", w.Code, w.Body)
	}
	finding, err := s.m.pg.GetFinding(f.FindingID)
	if err != nil || finding.TaskID == nil {
		t.Fatal("missing finding task", err)
	}
	for _, export := range []string{
		fmt.Sprintf("/api/report?task=%d&scope=all&format=docx", *finding.TaskID),
		fmt.Sprintf("/api/exploration/findings/export?scope=selected&ids=%d&format=docx", f.FindingID),
		fmt.Sprintf("/api/exploration/findings/export?scope=selected&ids=%d&format=md-zip", f.FindingID),
	} {
		w = req("GET", export, "")
		if w.Code != 200 {
			t.Fatal(export, w.Code, w.Body)
		}
		z, err := zip.NewReader(bytes.NewReader(w.Body.Bytes()), int64(w.Body.Len()))
		if err != nil {
			t.Fatal(err)
		}
		foundImage, foundCaption, foundStep := false, false, false
		for _, entry := range z.File {
			if strings.HasSuffix(entry.Name, ".png") {
				r, err := entry.Open()
				if err != nil {
					t.Fatal(err)
				}
				_, err = png.Decode(r)
				r.Close()
				if err != nil {
					t.Fatal(err)
				}
				foundImage = true
			}
			if entry.Name == "word/document.xml" || strings.HasSuffix(entry.Name, ".md") {
				r, err := entry.Open()
				if err != nil {
					t.Fatal(err)
				}
				raw, err := io.ReadAll(r)
				r.Close()
				if err != nil {
					t.Fatal(err)
				}
				foundCaption = foundCaption || strings.Contains(string(raw), "edited note")
				foundStep = foundStep || strings.Contains(string(raw), "Recorded fixture step")
			}
		}
		if !foundImage || !foundCaption || !foundStep {
			t.Fatal("export omitted screenshot or caption", export)
		}
	}
	if w = req("DELETE", path, ""); w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	if w = req("GET", path, ""); w.Code != 404 {
		t.Fatal(w.Code, w.Body)
	}
}
