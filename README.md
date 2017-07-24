# EcServ

**EcServ** is a very simple, but very flexible webserver using Go.

## Features of EcServ

- Automatic SSL certificate aquisition from Let's Encrypt
- Full HTTP/2.0. Optional support for HTTP/1.1 for redirection only.

## How to use

First install Go tools from [Golang](http://golang.org) and setup the configurations.
Then,

```
go get github.com/ilango100/ecserv
cd $GOPATH/src/github.com/iango100/ecserv
make
./ecserv
```
It will automatically start the setup process in interactive mode, i.e. it asks for you to set the intial values.

After setup if you need to edit any settings, just edit the .ecserv file in the source file directory:
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
- `email` is used for setting up account at Let's Encrypt
- `cert` is the folder in which certificates are stored. Default recommended.
- `domain` is your domain in which you want to set your website
- `errlog` is the file in which error logs are written //Not implemented

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

