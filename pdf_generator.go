package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type RequestPdf struct {
	header      string
	body        string
	footer      string
	size        string
	orientation string
}

// new request to pdf function
func NewRequestPdf(size, orientation string) *RequestPdf {
	return &RequestPdf{
		size:        size,
		orientation: orientation,
	}
}

//parsing template function

func (r *RequestPdf) parseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Panic(err)
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.Panic(err)
		return "", err
	}
	fmt.Println(buf.String())
	return buf.String(), nil
}

func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {
	b, err := r.parseTemplate(templateFileName, data)
	r.body = b
	if err != nil {
		return err
	}
	return nil
}

func (r *RequestPdf) ParseTemplateHeader(templateFileName string, data interface{}) error {
	b, err := r.parseTemplate(templateFileName, data)
	r.header = b
	if err != nil {
		return err
	}
	return nil
}
func (r *RequestPdf) ParseTemplateFooter(templateFileName string, data interface{}) error {
	b, err := r.parseTemplate(templateFileName, data)
	r.footer = b
	if err != nil {
		return err
	}
	return nil
}

// generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {
	t := time.Now().Unix()
	// write whole the body

	if _, err := os.Stat("storage/cloneTemplate/"); os.IsNotExist(err) {
		errDir := os.Mkdir("storage/cloneTemplate/", 0777)
		if errDir != nil {
			log.Fatal(errDir)
		}
	}
	err := os.WriteFile("storage/cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err != nil {
		panic(err)
	}

	f, err := os.Open("storage/cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	page := wkhtmltopdf.NewPageReader(f)
	// page.PageOptions.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)
	pdfg.PageSize.Set(r.size)
	pdfg.Orientation.Set(r.orientation)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll("storage/cloneTemplate/")

	return true, nil
}
