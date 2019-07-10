package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"pdfconverter/models"
	u "pdfconverter/pdfGenerator"
)

func main() {
	r := mux.NewRouter()
	// On the default page we will simply serve our static index page.
	r.HandleFunc("/generate", GeneratePdf).Methods("POST")
	http.ListenAndServe(":3001", r)
}

func GeneratePdf(w http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	//html template data
	var templateData models.TreatmentSummary
	json.Unmarshal(body, &templateData)

	r := u.NewRequestPdf("")
	//html template path
	templatePath := "templates/treatment_summary.html"
	//path for download pdf
	outputPath := "storage/treatment_summary.pdf"

	response := models.GeneratedPdf{}
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		ok, _ := r.GeneratePDF(outputPath)
		response.Message = "pdf generated successfully"
		response.Url = outputPath
		fmt.Println(ok, "pdf generated successfully")
	} else {
		response.Message = "pdf generated unsuccessfully"
		response.Url = err.Error()
		fmt.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templateData)
}
