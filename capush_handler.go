// +build capush

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CAPHandler struct {
	etags           map[string]string
	deps            map[string][]string
	notFoundHandler http.Handler
}

func (c *CAPHandler) Send(rw http.ResponseWriter, f io.Reader) {
	io.Copy(rw, f)
}

func (c *CAPHandler) SendFile(rw http.ResponseWriter, file string) {
	f, err := os.Open(file)
	if err == nil {
		c.Send(rw, f)
	}
}

func (c *CAPHandler) Etag(req *http.Request) (string, error) {
	//search for etag
	filename := req.URL.Path[1:]
	etag, found := c.etags[filename]

	if found {
		return etag, nil
	} else {
		return nil, errors.New("File not found")
	}
}

func (c *CAPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//search for etag
	etag, found := c.etags[req.URL.Path[1:]]

	//not found, send not found
	if !found {
		c.notFoundHandler.ServeHTTP(rw, req)

	} else /*found*/ {
		rw.Write([]byte(etag))
	}
}

func createHandler() http.Handler {
	ets, err := depEtags(set.Root)
	if err != nil {
		fmt.Println("Error generating etags... Falling back to normal handler...")
		return http.FileServer(http.Dir(set.Root))
	}
	dep, _ := genDeps(set.Root)

	handler := &CAPHandler{
		etags:           ets,
		deps:            dep,
		notFoundHandler: http.NotFoundHandler,
	}

	return handler
}
