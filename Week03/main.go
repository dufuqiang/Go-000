package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

type appHandler struct {
	content string
}

func (aH *appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from a HandleFunc #1!\n")
}

type debugHandler struct{}

func (dH *debugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pprof.Index(w, r)
}

// 业务逻辑 API
func newAppServer() *http.Server {
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: &appHandler{"hello"},
	}
	return server
}



// debug API
func newDebugServer() *http.Server {
	server := &http.Server{
		Addr:    "localhost:6060",
		Handler: &debugHandler{},
	}
	return server
}

func app(shutdownCh chan bool, closed chan struct{}) {
	g := new(errgroup.Group)
	appServer := newAppServer()
	debugServer := newDebugServer()

	g.Go(func() error {
		if err := appServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := debugServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	go func() {
		<-shutdownCh
		appServer.Shutdown(context.Background())
		debugServer.Shutdown(context.Background())
		close(closed)
	}()
	if err := g.Wait(); err != nil {
		appServer.Shutdown(context.Background())
		debugServer.Shutdown(context.Background())
		close(closed)
	}
}
func main()  {
	quit := make(chan os.Signal, 1)
	shutdown := make(chan bool)
	closed := make(chan struct{})

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go app(shutdown, closed)

	select {
		case <-quit:
			log.Println("receive your info, shutting down")
			shutdown <- true
	}

	<-closed
	log.Println("end")

}


