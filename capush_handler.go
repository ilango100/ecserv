// +build capush

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type CAPHandler struct {
	Etags           map[string]string
	Deps            map[string][]string
	NotFoundHandler http.Handler
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

func (c *CAPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//search for etag
	filename := req.URL.Path[1:]
	etag, found := c.Etags[filename]

	//not found, send not found
	if !found {
		c.NotFoundHandler.ServeHTTP(rw, req)

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
		Etags:           ets,
		Deps:            dep,
		NotFoundHandler: http.NotFoundHandler,
	}

	return handler
}
