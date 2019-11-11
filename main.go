package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Echo struct {
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	Headers    Header                 `json:"headers"`
	Body       string                 `json:"body"`
	HostName   string                 `json:"hostname"`
	SubDomains []string               `json:"subdomains"`
	Query      map[string][]string    `json:"query"`
	Protocol   string                 `json:"protocol"`
	RemoteAddr string                 `json:"remoteaddr"`
	Os         map[string]interface{} `json:"os"'`
}

type Header map[string][]string

type server struct {
	logger *log.Logger
	router *http.ServeMux
}

type bodyRecorder struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *bodyRecorder) Write(b []byte) (int, error) {
	count, err := w.ResponseWriter.Write(b)
	w.Writer.Write(b)
	return count, err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run() error {

	s := server{
		logger: log.New(os.Stdout, "", 0),
		router: http.NewServeMux(),
	}

	fmt.Println("Server is listening on :8080")

	return http.ListenAndServe(":8080", handlers.CombinedLoggingHandler(os.Stdout, s.routes()))
}

func (s *server) routes() *http.ServeMux {
	s.router.Handle("/", s.withLogging(s.all()))
	return s.router
}

func (s *server) all() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(buildEcho(r)); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recorder := bodyRecorder{ResponseWriter: w, Writer: s.logger.Writer()}
		h(&recorder, r)
	}
}
func (h Header) MarshalJSON() ([]byte, error) {
	simple := make(map[string]string)
	for k, v := range h {
		simple[strings.ToLower(k)] = strings.Join(v, ", ")
	}
	return json.Marshal(simple)
}

func buildEcho(r *http.Request) Echo {
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)

	osInfo := make(map[string]interface{})
	osInfo["hostname"], _ = os.Hostname()

	e := Echo{
		Path:       r.RequestURI,
		Body:       bodyString,
		Method:     r.Method,
		Headers:    Header(r.Header),
		SubDomains: strings.Split(strings.Split(r.Host, ":")[0], "."),
		HostName:   r.Host,
		Os:         osInfo,
		Protocol:   strings.ToLower(r.Proto),
		RemoteAddr: r.RemoteAddr,
		Query:      r.URL.Query(),
	}

	return e
}
