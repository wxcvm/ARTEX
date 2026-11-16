package kb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// PDF extraction uses optional local pdftotext (Poppler or compatible).
// A private temporary input file provides seekable input on Windows builds
// that cannot read a PDF from stdin. It is removed on every exit path.
// No shell or external service is involved.
// Image-only PDFs fall back to optional local OCR.
func extractPDFText(ctx context.Context,raw []byte,opts PDFOptions)(string,error){
	header:=raw
	if len(header)>1024{header=header[:1024]}
	if !bytes.Contains(header,[]byte("%PDF-")){return "",errors.New("invalid PDF header")}
	dir,err:=os.MkdirTemp("","artex-pdf-")
	if err!=nil{return "",fmt.Errorf("create local PDF workspace: %w",err)}
	defer os.RemoveAll(dir)
	input:=filepath.Join(dir,"document.pdf")
	if err=os.WriteFile(input,raw,0600);err!=nil{return "",fmt.Errorf("write local PDF input: %w",err)}
	if opts.ForceOCR{return extractPDFOCR(ctx,dir,input,opts.OCRLanguage)}
	path,err:=exec.LookPath("pdftotext")
	if err!=nil{return "",errors.New("PDF text import requires local pdftotext (Poppler) on PATH; force OCR is also available")}
	textCtx,cancel:=context.WithTimeout(ctx,30*time.Second)
	defer cancel()
	cmd:=exec.CommandContext(textCtx,path,"-enc","UTF-8","-layout",input,"-")
	var output limitedPDFText
	cmd.Stdout=&output
	var diagnostic limitedPDFText
	cmd.Stderr=&diagnostic
	if err:=cmd.Run();err!=nil{
		if output.overflow{return "",errors.New("PDF extracted text exceeds 4 MiB")}
		if textCtx.Err()!=nil{return "",textCtx.Err()}
		return "",fmt.Errorf("PDF text extraction failed: %w (%s)",err,strings.TrimSpace(diagnostic.buf.String()[:min(diagnostic.buf.Len(),512)]))
	}
	if output.overflow{return "",errors.New("PDF extracted text exceeds 4 MiB")}
	text:=strings.TrimSpace(output.buf.String())
	if text==""{return extractPDFOCR(ctx,dir,input,opts.OCRLanguage)}
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
