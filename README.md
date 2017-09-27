# EcServ
[![Build Status](https://travis-ci.org/ilango100/ecserv.svg?branch=master)](https://travis-ci.org/ilango100/ecserv) [![Go Report Card](https://goreportcard.com/badge/github.com/ilango100/ecserv)](https://goreportcard.com/report/github.com/ilango100/ecserv)

**EcServ** is a very simple, but very flexible webserver with automated SSL Certificate acquisition using Go.

## Features of EcServ

- Automatic SSL certificate aquisition from Let's Encrypt
- Full HTTP/2.0. Support for HTTP/1.1 for redirection only.
- Cache Aware Server Push (CAPush). Pushes the resources required for a page (stylesheet, script, etc.)
- Gzip compression support

## How to use

First install Go tools from [Golang](http://golang.org) and setup the configurations.
Then,

```
go get github.com/ilango100/ecserv
cd $GOPATH/src/github.com/iango100/ecserv
go get
go build
./ecserv
```
It will automatically start the setup process in interactive mode, i.e. it asks for you to set the intial values.

When the server starts for the first time, make a request to the server. When the certificate is acquired for first time, it asks whether to accept terms; Accept it on first request, and the server is ready to go!

After setup if you need to edit any settings, just edit the ecset file in the source file directory:
```
{
 "root": "C:\\Users\\<username>\\EcServ",
 "email": "username@example.com",
 "cert": "cert",
 "domain": "<your-domain-name>",
 "errlog": "errors.log"
}
```

Where 
- `root` is your root directory in which you have to put your site files.
- `email` is used for setting up account at Let's Encrypt.
- `cert` is the folder in which certificates are stored. Default recommended.
- `domain` is your domain in which you want to set your website.
- `errlog` is the file in which error logs are written //Not implemented

## Configuring CAPush

Cache Aware Server is a feature for the modern HTTP/2 servers. It enables pushing the required resources like stylesheet, javascript, images etc., without the browser making additional round trip.

To enable CAPush, just include `deps.json` file in your root directory.
The format of the `deps.json` should be like:
```
{
	"index.html" : ["style.css", "script.js"],
	"page1.html" : ["style.css", "img1.png"],
	"page2.html" : ["script.js", "img1.png", "img2.jpg"]
}
```
If you want to include your stylesheet and script in all your files, include a special name "global" with dependants:
```
{
	"global" : ["style.css", "script.js"],
	"page1.html" : ["img1.png"],
	"page2.html" : ["img1.png", img2.jpg"]
}
```
Now, style.css and script.js will be pushed with all the main files automatically.

Here index.html is the main file and style.css and script.js are the dependant files, which are to be pushed with index.html. If there is a subdirectory, create a `deps.json` file in the subdirectory separately.

Now start the server, the server will take care of pushing the dependant files along with the main file. If the server detects the browser already has cached copy of style.css, it just pushes 304 Not Modified response, which also avoids the browser revalidating the cache. 

Internally, CAPush uses Etags to check for file updates. Due to CAPush, your site will be very fast. EcServ is one of the few servers that have implemented the Cache Aware Server Push mechanism.

## Extending / CGI

If you wanna make this a forum etc, you can easily do so by defining your handler in the handler.go file.

For example, add this code to the handler.go file (import "math/rand")
```
mux.HandleFunc("/rand",func(w http.ResponsWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(rand.Intn(10))
	})
```
This creates a page that gives a random number within 10 each time the "/rand" page is requested.

Similarly,
```
mux.HandleFunc("/ip",func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(r.RemoteAddr))
	})
```
This creates a page "/ip" that shows the public ip address (with port) of the client that accessed the site.

## Bugs / Contributing

You can contribute by creating a pull request.

If you have found any bugs, please open an issue.

