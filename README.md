# HTTP Server Box

This package allows you to start an HTTP server that can be extended by other packages.

> ## Disclaimer
> This package is a just wrapper for [http](https://pkg.go.dev/net/http) and [quic](https://github.com/quic-go/quic-go) packages. 

## Motivation

Initially, I looked for some packages that would simplify the implementation of an HTTP server compatible with HTTP3.

## Installation

```shell
go get github.com/Lucas-Palomo/go-http-server-box
```

## Examples of usage

### Simple HTTP3 Server

```go
package main

import (
	"github.com/Lucas-Palomo/go-http-server-box/server"
	"net/http"
)


func main() {
	app := server.New(":8080", server.HTTP3)
	err := app.LoadTLSCert("./certs/test_cert.pem", "./certs/test_cert.key")
	if err != nil {
		panic(err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
		w.WriteHeader(http.StatusOK)
	})

	server.Launch(app, handler)
}
```
### Echo HTTP3 Server

```go
package main

import (
	"github.com/Lucas-Palomo/go-http-server-box/server"
	"github.com/labstack/echo/v4"
	"net/http"
)


func main() {
	app := server.New(":8080", server.HTTP3)
	err := app.LoadTLSCert("./certs/test_cert.pem", "./certs/test_cert.key")
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	server.Launch(app, e)
}
```
