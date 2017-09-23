// +build capush

package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
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

func (c *CAPHandler) PushAllDeps(rw http.ResponseWriter, filename string) {
	deps, havedeps := c.Deps[filename]
	if havedeps {
		p, noerr := rw.(http.Pusher)
		if noerr {
			for _, dep := range deps {
				p.Push("/"+dep, nil)
			}
		}
	}
	rw.Header().Set("etag", c.Etags[filename])
}

func (c *CAPHandler) PushModDeps(rw http.ResponseWriter, filename, oldetag string) bool {
	newetag := c.Etags[filename]

	if len(oldetag)%3 != 0 || len(oldetag) != len(newetag) {
		c.PushAllDeps(rw, filename)
		return true
	}
	if oldetag == newetag {
		return false
	}
	mainchanged := oldetag[:3] != newetag[:3]

	deps, havedeps := c.Deps[filename]
	if havedeps {
		p, noerr := rw.(http.Pusher)
		if noerr {
			tl := len(newetag)
			for i := 3; i < tl; i += 3 {
				if oldetag[i:i+3] != newetag[i:i+3] {
					p.Push("/"+deps[i/3-1], nil)
				}
			}
		}
	}

	rw.Header().Set("etag", newetag)
	return mainchanged
}

func (c *CAPHandler) mimeType(filename string) string {
	return mime.TypeByExtension(path.Ext(filename))
}

func (c *CAPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//search for etag
	filename := req.URL.Path[1:]
	_, found := c.Etags[filename]

	//not found, send not found
	if !found {
		c.NotFoundHandler.ServeHTTP(rw, req)

	} else /*found*/ {
		oldetag, etf := req.Header["If-None-Match"]

		//Fresh send
		if !etf {
			c.PushAllDeps(rw, filename)

			//Write headers
			rw.Header().Set("cache-control", "public, max-age=172800")

			//Set content type
			mime := c.mimeType(filename)
			if mime != "" {
				rw.Header().Set("content-type", mime)
			}

			//Send file
			rw.WriteHeader(200)
			c.SendFile(rw, path.Join(set.Root, filename))

		} else /*Check for update and send */ {
			rw.WriteHeader(304)
			rw.Write([]byte(oldetag[0]))
		}
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
		NotFoundHandler: http.NotFoundHandler(),
	}

	return handler
}
