package kb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ocrLanguagePattern = regexp.MustCompile(`^[A-Za-z0-9_]+(\+[A-Za-z0-9_]+)*$`)

func extractPDFOCR(ctx context.Context,dir,input,language string)(string,error){
	if language==""{language="eng"}
	if len(language)>128||!ocrLanguagePattern.MatchString(language){return "",errors.New("invalid OCR language")}
	paths:=map[string]string{}
	for _,name:=range []string{"pdfinfo","pdftoppm","tesseract"}{
		p,err:=exec.LookPath(name)
		if err!=nil{return "",fmt.Errorf("scanned PDF OCR requires local %s on PATH (Poppler and Tesseract)",name)}
		paths[name]=p
	}
	ctx,cancel:=context.WithTimeout(ctx,5*time.Minute)
	defer cancel()
	run:=func(name string,args ...string)(string,error){
		cmd:=exec.CommandContext(ctx,paths[name],args...)
		cmd.Env=append(os.Environ(),"LC_ALL=C")
		var stdout,stderr limitedPDFText
		cmd.Stdout=&stdout;cmd.Stderr=&stderr
		err:=cmd.Run()
		if ctx.Err()!=nil{return "",ctx.Err()}
		if stdout.overflow{return "",errors.New("OCR output exceeds 4 MiB")}
		if err!=nil{return "",fmt.Errorf("%s failed: %w (%s)",name,err,strings.TrimSpace(stderr.buf.String()[:min(stderr.buf.Len(),512)]))}
		return stdout.buf.String(),nil
	}
	info,err:=run("pdfinfo",input)
	if err!=nil{return "",err}
	pages:=0
	for _,line:=range strings.Split(info,"\n"){
		if strings.HasPrefix(line,"Pages:"){
			pages,err=strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line,"Pages:")))
			if err!=nil{return "",errors.New("could not read PDF page count")}
			break
		}
	}
	if pages<=0{return "",errors.New("could not read PDF page count")}
	var out limitedPDFText
	prefix:=filepath.Join(dir,"ocr-page")
	imagePath:=prefix+".png"
	for page:=1;page<=pages;page++{
		n:=strconv.Itoa(page)
		// Render and remove one bounded-size page at a time. Never accumulate
		// a directory of page images or invoke a shell.
		if _,err=run("pdftoppm","-f",n,"-l",n,"-singlefile","-scale-to","2400","-png",input,prefix);err!=nil{return "",fmt.Errorf("render page %d: %w",page,err)}
		value,err:=run("tesseract",imagePath,"stdout","-l",language)
		_ = os.Remove(imagePath)
		if err!=nil{return "",fmt.Errorf("OCR page %d: %w",page,err)}
		if value=strings.TrimSpace(value);value!=""{
			if out.buf.Len()>0{_,_=out.Write([]byte("\n\n"))}
			_,_=out.Write([]byte(value))
			if out.overflow{return "",errors.New("OCR extracted text exceeds 4 MiB")}
		}
	}
	if out.buf.Len()==0{return "",errors.New("OCR found no text in PDF")}
	return out.buf.String(),nil
}
