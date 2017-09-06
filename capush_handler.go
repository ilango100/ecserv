// +build capush

package main

import (
	"fmt"
	"net/http"
)

type CAPHandler struct {
	etags map[string]string
	deps  map[string][]string
}

func (c *CAPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//search for etag
	etag, found := c.etags[req.URL.Path[1:]]

	//not found, send not found
	if !found {
		rw.WriteHeader(404)
		rw.Write([]byte("File not found!"))

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
		etags: ets,
		deps:  dep,
	}

	return handler
}
