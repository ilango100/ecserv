package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

type CAPHandler struct {
	Etags           map[string]string
	Deps            map[string][]string
	Root            string
	IndexFile       string
	NotFoundHandler http.Handler
}

func (c *CAPHandler) send(rw http.ResponseWriter, f io.Reader, compr bool) {
	if !compr {
		rw.WriteHeader(200)
		io.Copy(rw, f)
	} else {
		rw.Header().Set("Content-Encoding", "gzip")
		rw.WriteHeader(200)
		gz := gzip.NewWriter(rw)
		io.Copy(gz, f)
		gz.Close()
	}
}

func (c *CAPHandler) sendFile(rw http.ResponseWriter, file string, compr bool) {
	if f, err := os.Open(file); err == nil {
		defer f.Close()
		c.send(rw, f, compr)
	}
}

func (c *CAPHandler) pushAllDeps(rw http.ResponseWriter, filename string) {
	if deps, havedeps := c.Deps[filename]; havedeps {
		p, noerr := rw.(http.Pusher)
		if noerr {
			for _, dep := range deps {
				p.Push("/"+dep, nil)
			}
		}
	}
	rw.Header().Set("etag", "\""+c.Etags[filename]+"\"")
}

func (c *CAPHandler) pushModDeps(rw http.ResponseWriter, filename, oldetag string) bool {
	newetag := c.Etags[filename]

	if len(oldetag)%3 != 0 || len(oldetag) != len(newetag) {
		c.pushAllDeps(rw, filename)
		return true
	}

	if deps, havedeps := c.Deps[filename]; havedeps {
		p, noerr := rw.(http.Pusher)
		if noerr {
			tl := len(newetag)
			for i := 3; i < tl; i += 3 {
				if oldetag[i:i+3] != newetag[i:i+3] {
					p.Push("/"+deps[i/3-1], nil)
				} else {
					h := http.Header(make(map[string][]string))
					h.Set("If-None-Match", "\""+oldetag[i:i+3]+"\"")
					p.Push("/"+deps[i/3-1], &http.PushOptions{Header: h})
				}
			}
		}
	}

	rw.Header().Set("etag", "\""+newetag+"\"")
	return oldetag != newetag
}

func (c *CAPHandler) typeAndSendFile(rw http.ResponseWriter, filename string, compr bool) {
	//Set content type
	mime := mime.TypeByExtension(path.Ext(filename))
	if mime != "" {
		rw.Header().Set("content-type", mime)
	}

	//Send file
	c.sendFile(rw, filename, compr)
}

func (c *CAPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//Set correct filename
	filename := req.URL.Path[1:]
	if filename == "" {
		filename = c.IndexFile
	}

	//search for etag
	_, found := c.Etags[filename]

	//not found, send not found
	if !found {
		c.NotFoundHandler.ServeHTTP(rw, req)

	} else /*found*/ {
		oldetags, etf := req.Header["If-None-Match"]

		accenc, compr := req.Header["Accept-Encoding"]
		if compr {
			compr = strings.Contains(accenc[0], "gzip")
		}

		//Fresh send
		if !etf {
			c.pushAllDeps(rw, filename)

			//Write headers
			rw.Header().Set("cache-control", "public, max-age=172800")

			c.typeAndSendFile(rw, path.Join(c.Root, filename), compr)

		} else /*Check for update and send */ {

			//Extract correct etag
			oldetag := oldetags[0]
			oldetag = strings.Trim(oldetag, "\" ")

			//Push modified deps
			if c.pushModDeps(rw, filename, oldetag) {
				rw.Header().Set("cache-control", "public, max-age=172800")
				c.typeAndSendFile(rw, path.Join(c.Root, filename), compr)

			} else {
				rw.WriteHeader(304)
			}
		}
	}
}

func createCAPHandler(root string) http.Handler {
	ets, err := depEtags(root)
	if err != nil {
		fmt.Println("Push disabled: Error generating etags")
		return http.FileServer(http.Dir(root))
	}
	dep, _ := genDeps(root)

	handler := &CAPHandler{
		Etags:           ets,
		Deps:            dep,
		Root:            root,
		IndexFile:       "index.html",
		NotFoundHandler: http.NotFoundHandler(),
	}

	return handler
}
