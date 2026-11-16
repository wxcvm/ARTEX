package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

const maxFindingScreenshotBytes = 8 << 20

func (s *Server) registerFindingScreenshots(mux *http.ServeMux) {
	base := "/api/exploration/findings/{id}/screenshots"
	mux.HandleFunc("GET "+base, s.listFindingScreenshots)
	mux.HandleFunc("POST "+base, s.createFindingScreenshot)
	mux.HandleFunc("GET "+base+"/{screenshot_id}", s.readFindingScreenshot)
	mux.HandleFunc("PATCH "+base+"/{screenshot_id}", s.changeFindingScreenshot)
	mux.HandleFunc("DELETE "+base+"/{screenshot_id}", s.changeFindingScreenshot)
}

// Decode and re-encode raster images: retain the visible evidence, but do not
// serve arbitrary uploaded payloads or embedded metadata as authenticated files.
func sanitizeFindingScreenshot(raw []byte) (string, []byte, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil || (format != "png" && format != "jpeg") {
		return "", nil, errors.New("only valid PNG and JPEG screenshots are supported")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || int64(cfg.Width)*int64(cfg.Height) > 16_000_000 {
		return "", nil, errors.New("screenshot dimensions exceed 16 megapixels")
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return "", nil, errors.New("invalid screenshot image")
	}
	var out bytes.Buffer
	mime := "image/png"
	if format == "jpeg" {
		mime = "image/jpeg"
		err = jpeg.Encode(&out, img, &jpeg.Options{Quality: 95})
	} else {
		encoder := png.Encoder{CompressionLevel: png.BestSpeed}
		err = encoder.Encode(&out, img)
	}
	if err != nil {
		return "", nil, err
	}
	if out.Len() > maxFindingScreenshotBytes {
		return "", nil, errors.New("encoded screenshot exceeds 8 MiB")
	}
	return mime, out.Bytes(), nil
}

func (s *Server) listFindingScreenshots(w http.ResponseWriter, r *http.Request) {
	id, ok := s.findingTrafficAccess(w, r, false)
	if !ok {
		return
	}
	out, err := s.m.pg.ListFindingScreenshots(r.Context(), id)
	if err != nil {
		evidenceError(w, err)
		return
	}
	writeJSON(w, 200, out)
}

func (s *Server) createFindingScreenshot(w http.ResponseWriter, r *http.Request) {
	id, ok := s.findingTrafficAccess(w, r, true)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxFindingScreenshotBytes+(64<<10))
	if err := r.ParseMultipartForm(maxFindingScreenshotBytes); err != nil {
		writeErr(w, 400, "invalid screenshot upload or file too large")
		return
	}
	defer r.MultipartForm.RemoveAll()
	file, _, err := r.FormFile("image")
	if err != nil {
		writeErr(w, 400, "image is required")
		return
	}
	defer file.Close()
	raw, err := io.ReadAll(io.LimitReader(file, maxFindingScreenshotBytes+1))
	if err != nil || len(raw) > maxFindingScreenshotBytes {
		writeErr(w, 400, "screenshot exceeds 8 MiB")
		return
	}
	mime, data, err := sanitizeFindingScreenshot(raw)
	if err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	caption := strings.TrimSpace(r.FormValue("caption"))
	if utf8.RuneCountInString(caption) > 4000 {
		writeErr(w, 400, "caption exceeds 4000 characters")
		return
	}
	var bindingID *int64
	if rawID := r.FormValue("binding_id"); rawID != "" {
		v, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil || v <= 0 {
			writeErr(w, 400, "invalid binding id")
			return
		}
		bindingID = &v
	}
	out, err := s.m.pg.CreateFindingScreenshot(r.Context(), id, bindingID, caption, mime, data)
	if err != nil {
		evidenceError(w, err)
		return
	}
	writeJSON(w, 201, out)
}

func (s *Server) readFindingScreenshot(w http.ResponseWriter, r *http.Request) {
	id, ok := s.findingTrafficAccess(w, r, false)
	if !ok {
		return
	}
	mime, data, err := s.m.pg.ReadFindingScreenshot(r.Context(), id, r.PathValue("screenshot_id"))
	if err != nil {
		evidenceError(w, err)
		return
	}
	w.Header().Set("Content-Type", mime)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(data)
}

func (s *Server) changeFindingScreenshot(w http.ResponseWriter, r *http.Request) {
	id, ok := s.findingTrafficAccess(w, r, true)
	if !ok {
		return
	}
	var caption *string
	if r.Method == http.MethodPatch {
		var body struct {
			Caption *string `json:"caption"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10)).Decode(&body); err != nil || body.Caption == nil {
			writeErr(w, 400, "caption is required")
			return
		}
		value := strings.TrimSpace(*body.Caption)
		if utf8.RuneCountInString(value) > 4000 {
			writeErr(w, 400, "caption exceeds 4000 characters")
			return
		}
		caption = &value
	}
	if err := s.m.pg.ChangeFindingScreenshot(r.Context(), id, r.PathValue("screenshot_id"), caption); err != nil {
		evidenceError(w, err)
		return
	}
	writeJSON(w, 200, map[string]bool{"ok": true})
}
