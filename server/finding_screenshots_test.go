package server

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
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
	if w = req("DELETE", path, ""); w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	if w = req("GET", path, ""); w.Code != 404 {
		t.Fatal(w.Code, w.Body)
	}
}
