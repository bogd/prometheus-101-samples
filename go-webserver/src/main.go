package main

import (
    "fmt"
    "log"
    "net/http"
	"strconv"
	"time"
	"math/rand"
	
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

var name string = "STUDENT_NAME_HERE"

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	}, 
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}


func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/hello" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }


    fmt.Fprintf(w, "Hello, %s!", name)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/slow" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
	
    delay := float64(rand.Intn(500))/100  // Random delay up to 5 seconds, with 10ms granularity
    time.Sleep(time.Duration(delay)*time.Second)

    fmt.Fprintf(w, "Hello, %s, after %6.2f seconds!", name, delay)
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	// Root endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())

	// Hello handler
	router.HandleFunc("/hello", helloHandler)
	
	// Slow handler
	router.HandleFunc("/slow", slowHandler)
	
	// Serving static files
	fileServer := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/").Handler(fileServer)

	fmt.Println("Serving requests on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}
