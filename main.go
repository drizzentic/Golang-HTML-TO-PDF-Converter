package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pdfconverter/models"
	u "pdfconverter/pdfGenerator"
	"strconv"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/generate", GeneratePdf).Methods("POST")
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("assets/"))))
	n := negroni.Classic()
	n.UseHandler(r)
	srv := &http.Server{
		Addr: "0.0.0.0:3001",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      n, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func GeneratePdf(w http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	var templateData models.TreatmentSummary
	err := json.Unmarshal(body, &templateData)
	if err != nil {
		handleError(w, err)
	}

	r := u.NewRequestPdf("")
	var templatePath = "templates/treatment_summary.html"
	t := time.Now().Unix()
	outputPath := fmt.Sprintf("storage/%s.pdf", GetMD5Hash(templateData.Name+strconv.FormatInt(int64(t), 10)))
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		go r.GeneratePDF(outputPath)

	} else {
		fmt.Println(err)
	}

	streamPDFbytes, _ := ioutil.ReadFile(outputPath)
	b := bytes.NewBuffer(streamPDFbytes)
	w.Header().Set("Content-type", "application/pdf")
	if _, err := b.WriteTo(w); err != nil {
		_, _ = fmt.Fprintf(w, "%s", err)
	}
	_, _ = w.Write([]byte("PDF Generated"))
}
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Fatal(err)
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
