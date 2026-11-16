package kb

import (
	"bytes"
	"compress/zlib"
	"context"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Rasterize the existing text fixture, then wrap only its pixels in a new PDF.
// The resulting PDF has no font resources or text operators.
func scannedKnowledgePDF(t *testing.T) []byte {
	t.Helper()
	dir:=t.TempDir()
	input:=filepath.Join(dir,"text.pdf")
	if err:=os.WriteFile(input,knowledgePDFTestFixture(),0600);err!=nil{t.Fatal(err)}
	prefix:=filepath.Join(dir,"raster")
	if out,err:=exec.Command("pdftoppm","-f","1","-l","1","-singlefile","-scale-to","2400","-png",input,prefix).CombinedOutput();err!=nil{t.Fatal(err,string(out))}
	f,err:=os.Open(prefix+".png");if err!=nil{t.Fatal(err)}
	img,_,err:=image.Decode(f);f.Close();if err!=nil{t.Fatal(err)}
	bounds:=img.Bounds()
	pixels:=make([]byte,bounds.Dx()*bounds.Dy())
	for y:=0;y<bounds.Dy();y++{for x:=0;x<bounds.Dx();x++{pixels[y*bounds.Dx()+x]=color.GrayModel.Convert(img.At(x+bounds.Min.X,y+bounds.Min.Y)).(color.Gray).Y}}
	var compressed bytes.Buffer
	zw:=zlib.NewWriter(&compressed)
	if _,err=zw.Write(pixels);err!=nil{t.Fatal(err)}
	if err=zw.Close();err!=nil{t.Fatal(err)}
	content:="q 612 0 0 792 0 0 cm /Im0 Do Q\n"
	objects:=[][]byte{
		[]byte("<< /Type /Catalog /Pages 2 0 R >>"),
		[]byte("<< /Type /Pages /Kids [3 0 R] /Count 1 >>"),
		[]byte("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /XObject << /Im0 4 0 R >> >> /Contents 5 0 R >>"),
		append(append([]byte(fmt.Sprintf("<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceGray /BitsPerComponent 8 /Filter /FlateDecode /Length %d >>\nstream\n",bounds.Dx(),bounds.Dy(),compressed.Len())),compressed.Bytes()...),[]byte("\nendstream")...),
		[]byte(fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream",len(content),content)),
	}
	var b bytes.Buffer
	b.WriteString("%PDF-1.4\n")
	offsets:=[]int{0}
	for i,obj:=range objects{offsets=append(offsets,b.Len());fmt.Fprintf(&b,"%d 0 obj\n",i+1);b.Write(obj);b.WriteString("\nendobj\n")}
	xref:=b.Len()
	fmt.Fprintf(&b,"xref\n0 %d\n0000000000 65535 f \n",len(offsets))
	for _,offset:=range offsets[1:]{fmt.Fprintf(&b,"%010d 00000 n \n",offset)}
	fmt.Fprintf(&b,"trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n",len(offsets),xref)
	return b.Bytes()
}

func TestKnowledgeScannedPDFOCR(t *testing.T){
	for _,name:=range []string{"pdftotext","pdfinfo","pdftoppm","tesseract"}{if _,err:=exec.LookPath(name);err!=nil{t.Skip("optional local OCR tool missing: "+name)}}
	langs,err:=exec.Command("tesseract","--list-langs").CombinedOutput()
	if err!=nil||!strings.Contains(string(langs),"eng"){t.Skip("English OCR language data missing")}
	scan:=scannedKnowledgePDF(t)
	for _,tc:=range []struct{name string; raw []byte; opts PDFOptions}{
		{"scanned fallback",scan,PDFOptions{}},
		{"forced text PDF",knowledgePDFTestFixture(),PDFOptions{ForceOCR:true,OCRLanguage:"eng"}},
	}{
		t.Run(tc.name,func(t *testing.T){
			text,err:=ExtractDocumentTextWithOptions(context.Background(),"fixture.pdf",tc.raw,tc.opts)
			if err!=nil||!strings.Contains(text,"ARTEX PDF fixture"){t.Fatal(text,err)}
		})
	}
	if _,err:=ExtractDocumentTextWithOptions(context.Background(),"fixture.pdf",scan,PDFOptions{ForceOCR:true,OCRLanguage:"--bad"});err==nil{t.Fatal("invalid language accepted")}
}
