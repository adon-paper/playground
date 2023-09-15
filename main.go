package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PaymentDetails struct {
	BuyerUserName              string
	BuyerCompanyName           string
	BuyerPhone                 string
	SellerCompanyName          string
	SellerPhone                string
	Amount                     string
	ExternalID                 string
	ReferenceNote              string
	PaymentDate                string
	DisbursementEstimationTime string
	BankName                   string
	BankAccountName            string
	BankAccountNumber          string
}

func main() {
	var (
		rootdir, _ = os.Getwd()
	)
	// This is the only way I have found to be able to serve files requested in the templates
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/",
		http.FileServer(http.Dir(path.Join(rootdir, "/assets/css/")))))

	http.Handle("/assets/img/", http.StripPrefix("/assets/img/",
		http.FileServer(http.Dir(path.Join(rootdir, "/assets/img/")))))

	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/",
		http.FileServer(http.Dir(path.Join(rootdir, "/assets/fonts/")))))

	http.HandleFunc("/", index)
	http.ListenAndServe("localhost:8080", nil)
}
func index(w http.ResponseWriter, r *http.Request) {
	pdf := NewRequestPdf(wkhtmltopdf.PageSizeA4, wkhtmltopdf.OrientationPortrait)
	templatePath := `template/payment-completed.html`
	//path for download pdf
	filename := "payment-completed.pdf"
	outputPath := `storage/` + filename
	data := PaymentDetails{
		BuyerUserName:              "Abdian Rizky",
		BuyerCompanyName:           "ADNS",
		BuyerPhone:                 "085157572725",
		SellerCompanyName:          "No Invoice",
		SellerPhone:                "085157530733",
		Amount:                     "10,000,000.00",
		ExternalID:                 "1694684908FZ8N2",
		ReferenceNote:              "aoiefj",
		PaymentDate:                "2023-09-14 16:48:35",
		DisbursementEstimationTime: "2023-09-15 14:00:00",
		BankName:                   "Bank Tabungan Pensiunan Nasional (BTPN)",
		BankAccountName:            "ABDIAN RIZKY RAMADAN",
		BankAccountNumber:          "90200175670",
	}

	if err := pdf.ParseTemplate(templatePath, data); err == nil {
		fmt.Println(err)
	}

	ok, err := pdf.GeneratePDF(outputPath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ok)
}
