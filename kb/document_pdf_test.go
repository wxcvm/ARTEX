package kb

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// A real minimal one-page PDF, including its cross-reference table.
func knowledgePDFTestFixture()[]byte{
	var b bytes.Buffer
	b.WriteString("%PDF-1.4\n")
	stream:="BT /F1 12 Tf 72 720 Td (ARTEX PDF fixture) Tj ET\n"
	objects:=[]string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream",len(stream),stream),
	}
	offsets:=[]int{0}
	for i,obj:=range objects{
		offsets=append(offsets,b.Len())
		fmt.Fprintf(&b,"%d 0 obj\n%s\nendobj\n",i+1,obj)
	}
	xref:=b.Len()
	fmt.Fprintf(&b,"xref\n0 %d\n0000000000 65535 f \n",len(offsets))
	for _,offset:=range offsets[1:]{fmt.Fprintf(&b,"%010d 00000 n \n",offset)}
	fmt.Fprintf(&b,"trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n",len(offsets),xref)
	return b.Bytes()
}

func TestKnowledgePDFExtraction(t *testing.T){
	if _,err:=exec.LookPath("pdftotext");err!=nil{t.Skip("optional local Poppler pdftotext not installed")}
	text,err:=ExtractDocumentText("fixture.pdf",knowledgePDFTestFixture())
	if err!=nil||!strings.Contains(text,"ARTEX PDF fixture"){t.Fatal(text,err)}
	if _,err:=ExtractDocumentText("invalid.pdf",[]byte("%PDF-1.4\nbroken"));err==nil{t.Fatal("accepted malformed PDF")}
}
func TestKnowledgePDFOutputLimit(t *testing.T){
	var w limitedPDFText
	p:=bytes.Repeat([]byte("x"),DocumentTextLimit+1)
	n,err:=w.Write(p)
	if err!=nil||n!=len(p)||!w.overflow||w.buf.Len()!=DocumentTextLimit{t.Fatal(n,err,w.overflow,w.buf.Len())}
}
