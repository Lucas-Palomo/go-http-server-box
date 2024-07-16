package box

import (
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go/http3"
	"net/http"
)

type proto int

const (
	HTTP1 proto = iota + 1
	HTTP2
	HTTP3
)

type Server struct {
	addr      string
	TLSConfig *tls.Config
	proto     proto
}

func New(addr string, proto proto) *Server {
	return &Server{
		addr:  addr,
		proto: proto,
	}
}

func (srv *Server) SetTLS(config *tls.Config) {
	srv.TLSConfig = config
}

func (srv *Server) LoadTLSCert(cert string, key string) error {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return err
	}

	if srv.TLSConfig == nil {
		srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}
		return nil
	}

	srv.TLSConfig.Certificates = append(srv.TLSConfig.Certificates, certificate)
	return nil
}

func Launch(srv *Server, handler http.Handler) error {

	switch srv.proto {
	case HTTP1, HTTP2:
		server := &http.Server{
			Addr:      srv.addr,
			Handler:   handler,
			TLSConfig: srv.TLSConfig,
		}

		if srv.proto == HTTP2 {
			return server.ListenAndServeTLS("", "")
		}
		return server.ListenAndServe()
	case HTTP3:
		quicSrv := &http3.Server{
			Addr:      srv.addr,
			Handler:   handler,
			TLSConfig: srv.TLSConfig,
		}

		hErr := make(chan error, 1)
		qErr := make(chan error, 1)

		go func() {
			app := New(quicSrv.Addr, HTTP2)
			app.SetTLS(quicSrv.TLSConfig)

			quicHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				quicSrv.SetQUICHeaders(res.Header())
				handler.ServeHTTP(res, req)
			})

			hErr <- Launch(app, quicHandler)
		}()
		go func() {
			qErr <- quicSrv.ListenAndServe()
		}()

		select {
		case err := <-hErr:
			quicSrv.Close()
			return err
		case err := <-qErr:
			return err
		}
	}

	return fmt.Errorf("unknown proto: %d", srv.proto)
}
