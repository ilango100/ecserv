package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

//settings variables
var set struct {
	Root   string `json:"root"`
	Email  string `json:"email"`
	Cert   string `json:"cert"`
	Domain string `json:"domain"`
	ErrLog string `json:"errlog"`
}

//settings file
var isfile string

func main() {

	fmt.Println("Starting server...")
	go httpServ()

	handler := createHandler()

	crtmgr := autocert.Manager{
		Cache:      autocert.DirCache(set.Cert),
		HostPolicy: autocert.HostWhitelist(set.Domain),
		Email:      set.Email,
		ForceRSA:   true,
		Prompt:     tosPrompt,
	}

	tconf := &tls.Config{
		GetCertificate: crtmgr.GetCertificate,
		NextProtos:     []string{"h2"},
	}

	server := &http.Server{
		Addr:      ":https",
		Handler:   handler,
		TLSConfig: tconf,
	}

	log.Fatal(server.ListenAndServeTLS("", ""))

}

func httpServ() {
	serv := &http.Server{
		Addr: ":http",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loc := "https://" + set.Domain + r.URL.String()
			w.Header().Set("Location", loc)
			w.WriteHeader(301)
			w.Write([]byte(loc))
		}),
	}

	serv.ListenAndServe()
}

func tosPrompt(url string) bool {
	fmt.Printf("\nDo you agree to Terms Of Service at %s? (y/n) ", url)
	var in string
	fmt.Scan(&in)
	if in == "y" || in == "yes" {
		return true
	}
	fmt.Println("Not accepted...")
	return false
}
