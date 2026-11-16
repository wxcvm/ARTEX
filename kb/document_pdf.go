package kb

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
	"time"
)

// PDF extraction uses the optional local Poppler pdftotext executable. No
// document is written to a temporary file, uploaded, or passed to a shell.
// Scanned/image-only PDFs need OCR, which is not performed here.
func extractPDFText(ctx context.Context,raw []byte)(string,error){
	header:=raw
	if len(header)>1024{header=header[:1024]}
	if !bytes.Contains(header,[]byte("%PDF-")){return "",errors.New("invalid PDF header")}
	path,err:=exec.LookPath("pdftotext")
	if err!=nil{return "",errors.New("PDF text import requires local pdftotext (Poppler) on PATH")}
	ctx,cancel:=context.WithTimeout(ctx,30*time.Second)
	defer cancel()
	cmd:=exec.CommandContext(ctx,path,"-enc","UTF-8","-layout","-","-")
	cmd.Stdin=bytes.NewReader(raw)
	var output limitedPDFText
	cmd.Stdout=&output
	cmd.Stderr=io.Discard
	if err:=cmd.Run();err!=nil{
		if output.overflow{return "",errors.New("PDF extracted text exceeds 4 MiB")}
		if ctx.Err()!=nil{return "",ctx.Err()}
		return "",errors.New("PDF text extraction failed (invalid or encrypted PDF)")
	}
	if output.overflow{return "",errors.New("PDF extracted text exceeds 4 MiB")}
	text:=strings.TrimSpace(output.buf.String())
	if text==""{return "",errors.New("PDF contains no extractable text; scanned documents require OCR")}
	return text,nil
}

type limitedPDFText struct{
	buf bytes.Buffer
	overflow bool
}

func (w *limitedPDFText)Write(p []byte)(int,error){
	remaining:=DocumentTextLimit-w.buf.Len()
	if len(p)>remaining{
		w.overflow=true
		// Drain stdout without retaining excessive output.
		if remaining>0{_,_=w.buf.Write(p[:remaining])}
		return len(p),nil
	}
	return w.buf.Write(p)
}
