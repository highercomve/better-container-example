package server

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/highercomve/better-container-example/utils"
)

type key int

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
)

type Interface struct {
	net.Interface
	Addr []net.Addr
}

type IndexData struct {
	GOARCH     string
	GOOS       string
	OS         string
	Host       string
	Kernel     string
	Machine    string
	Interfaces []Interface
}

func Start() {
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/healthz", healthz())
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/", serveTemplate(logger))

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)

	atomic.StoreInt32(&healthy, 1)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func serveTemplate(logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templatePath := filepath.Clean(r.URL.Path)

		if templatePath == "/" {
			templatePath = "/index.html"
		}

		if templatePath == "/favicon.ico" {
			return
		}

		lp := filepath.Join("templates", "layout.html")
		fp := filepath.Join("templates", templatePath)

		tmpl, _ := template.ParseFiles(lp, fp)

		data := getData()
		logger.Printf("%+v", data)
		err := tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			logger.Println(err.Error())
		}
	})
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() { logRequest(logger, r) }()
			next.ServeHTTP(w, r)
		})
	}
}

func logRequest(logger *log.Logger, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDKey).(string)
	if !ok {
		requestID = "unknown"
	}
	logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getData() *IndexData {
	ainterfaces, err := net.Interfaces()
	if err != nil {
		ainterfaces = []net.Interface{}
	}
	var interfaces []Interface

	for _, val := range ainterfaces {
		if val.HardwareAddr != nil {
			i := Interface{}
			i.Interface = val
			i.Addr, _ = val.Addrs()
			interfaces = append(interfaces, i)
		}
	}
	uname := &syscall.Utsname{}

	data := &IndexData{
		GOARCH:     runtime.GOARCH,
		GOOS:       runtime.GOOS,
		OS:         "",
		Kernel:     "",
		Machine:    "",
		Interfaces: interfaces,
	}
	if err := syscall.Uname(uname); err == nil {
		data.OS = utils.Int8ToStr(uname.Sysname[:])
		data.Kernel = utils.Int8ToStr(uname.Release[:])
		data.Machine = utils.Int8ToStr(uname.Machine[:])
	}

	return data
}
