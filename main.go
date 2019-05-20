package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type dumpedData struct {
	rq *http.Request
	content []byte
}

type byteReadCloser struct {
	rc io.ReadCloser
	content io.Reader
}

func (dr *byteReadCloser) Read(buf []byte) (int, error) {
	return dr.content.Read(buf)
}

func (dr *byteReadCloser) Close() error {
	return dr.rc.Close()
}

type dumpingHandler struct {
	targetHost string
	filename string
	dump chan dumpedData
}

func newDumpingHandler(targetHost, filename string) *dumpingHandler {
	d := dumpingHandler{
		targetHost: targetHost,
		filename: filename,
		dump: make(chan  dumpedData, 1000),
	}
	go d.runDumper(filename)
	return &d
}

func (d *dumpingHandler) handle(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	d.dump <-dumpedData{ rq: r, content: content }
	r.Body = &byteReadCloser{ r.Body, bytes.NewReader(content)}
	r.Host = d.targetHost
	r.URL.Host = d.targetHost
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (d *dumpingHandler) runDumper(filename string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for dd := range d.dump {
		writeStr(f, dd.rq.Method + " " + dd.rq.URL.Path)
		writeStr(f, "\n")
		for k, vv := range dd.rq.Header {
			for _, v := range vv {
				writeStr(f, fmt.Sprintf("%s: %s\n", k, v))
			}
		}
		writeStr(f, "\n")
		write(f, dd.content)
		f.Sync()
		if err != nil {
			log.Fatal(err)
		}
		writeStr(f, "\n---\n")
	}
}

func write(f *os.File, buf []byte) {
	_, err := f.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
}

func writeStr(f *os.File, s string) {
	write(f, []byte(s))
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func serve(port int, targetHost, filename string) {
	dumper := newDumpingHandler(targetHost, filename)
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dumper.handle(w, r)
		}),
	}
	log.Fatal(server.ListenAndServe())
}

func main() {
	targetPtr := flag.String("target", "", "The target to send traffic to as hostname:port")
	filePtr := flag.String("file", "", "The file to dump to")
	portPtr := flag.Int("port", 80, "The port to listen to")
	flag.Parse()
	if *targetPtr == "" || *filePtr == "" {
		fmt.Fprintf(os.Stderr,"Both -file and -target must be specified\n")
		os.Exit(1)
	}
	serve(*portPtr, *targetPtr, *filePtr)
}
