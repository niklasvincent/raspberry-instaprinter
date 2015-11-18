package main

import (
  "fmt"
  "io"
  "mime"
  "net/http"
  "crypto/tls"
  "net"
  "github.com/gorilla/mux"
)

type embeddedServer struct {

  /*
    custom struct that embeds golang's standard http.Server type
    another way of looking at this is that embeddedServer "inherits" from http.Server,
    though this is not strictly accurate. Have a look at note below for additional information
  */

  http.Server
  webserverCertificate string
  webserverKey string
}

func (srv *embeddedServer) ListenAndServeTLS(addr string) error {

  /*
    This is where we "hide" or "override" the default "ListenAndServeTLS" method so we modify it to accept
    hardcoded certificates and keys rather than the default filenames
    The default implementation of ListenAndServeTLS was obtained from:
    https://github.com/zenazn/goji/blob/master/graceful/server.go#L33
    and tls.X509KeyPair (http://golang.org/pkg/crypto/tls/#X509KeyPair) is used,
    rather than the default tls.LoadX509KeyPair
   */

  config := &tls.Config{
    MinVersion: tls.VersionTLS10,
  }
  if srv.TLSConfig != nil {
    *config = *srv.TLSConfig
  }
  if config.NextProtos == nil {
    config.NextProtos = []string{"http/1.1"}
  }

  var err error
  config.Certificates = make([]tls.Certificate, 1)
  config.Certificates[0], err = tls.X509KeyPair([]byte(srv.webserverCertificate), []byte(srv.webserverKey))
  if err != nil {
  	log.Error("Could not set up SSL certificate")
    return err
  }

  conn, err := net.Listen("tcp", addr)
  if err != nil {
  	log.Error("Could not bind to " + addr)
    return err
  }

  tlsListener := tls.NewListener(conn, config)
  log.Info("Running embedded webserver with SSL listening on: " + addr)
  return srv.Serve(tlsListener)
}

func startEmbeddedWebServer(listenAddress string) {
	embeddedTLSserver := &embeddedServer{
	  webserverCertificate: WebserverCertificate,
	  webserverKey: WebserverPrivateKey,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", serveStartPage).Methods("GET")
	router.HandleFunc("/asset/{name:[A-Za-z0-9-_]+}.{suffix:[a-z]+}", serveAsset).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./images/")))

	http.Handle("/", router)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))

	log.Fatal(embeddedTLSserver.ListenAndServeTLS(listenAddress))
}

func serveStartPage(w http.ResponseWriter, r *http.Request) {
	data, err := Asset("web/index.html")
	if err != nil {
	    // Asset was not found.
	}
  	io.WriteString(w, string(data))
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	suffix := params["suffix"]
	var filename = fmt.Sprintf("web/%s/%s.%s", suffix, name, suffix)
	log.Debug("Got request for " + filename)
	data, err := Asset(filename)
	if err != nil {
	    log.Error("Could not find file: " + filename)
	}
	var contentType = mime.TypeByExtension(fmt.Sprintf(".%s", suffix))
	w.Header().Set("Content-Type", contentType)
  	io.WriteString(w, string(data))
}
