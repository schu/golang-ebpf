// vim: set ts=8:sts=8:sw=8:noet

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/iovisor/iomodules/hover/daemon"
	"github.com/iovisor/iomodules/hover/util"
)

var listenSocket string
var helpFlag bool

func init() {
	const (
		listenSocketDefault = "127.0.0.1:5000"
		listenSocketHelp    = "address:port to serve up the api"
	)
	flag.StringVar(&listenSocket, "listen", listenSocketDefault, listenSocketHelp)
	flag.BoolVar(&helpFlag, "h", false, "print this help")
	flag.Usage = func() {
		fmt.Printf("Usage: %s -listen 0.0.0.0:5000\n", filepath.Base(os.Args[0]))
		fmt.Printf(" -listen ADDR:PORT  %s (default=%s)\n", listenSocketHelp, listenSocketDefault)
	}
}

func main() {
	flag.Parse()
	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	s := daemon.NewServer()
	go func() {
		<-c
		if s != nil {
			s.Close()
		}
		os.Exit(1)
	}()
	if s == nil {
		util.Warn.Println("Failed to start Hover Server")
		os.Exit(1)
	}
	util.Info.Printf("Hover Server listening on %s\n", listenSocket)
	http.ListenAndServe(listenSocket, s.Handler())
}
