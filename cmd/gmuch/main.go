package main

import (
	"net/http"
	"os"
	"time"

	stdlog "log"

	"golang.org/x/net/context"

	"github.com/gmuch/gmuch"
	"github.com/gmuch/gmuch/server"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/namsral/flag"

	shttp "github.com/gmuch/gmuch/server/http"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	listen = flag.String("listen", ":7080", "HTTP listen address")
	dbPath = flag.String("db-path", "", "The path to the notmuch database")
)

func main() {
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)
	stdlog.SetFlags(0)                             // flags are handled by Go kit's logger
	stdlog.SetOutput(log.NewStdlibAdapter(logger)) // redirect anything using stdlib log to us

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounter(stdprometheus.CounterOpts{
		Namespace: "gmuch",
		Subsystem: "api",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := metrics.NewTimeHistogram(time.Microsecond, kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "gmuch",
		Subsystem: "api",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys))

	var g server.GmuchService
	g = gmuch.New(*dbPath, logger)
	g = server.LoggingMiddleware(logger)(g)
	g = server.InstrumentingMiddleware(requestCount, requestLatency)(g)

	root := context.Background()
	mux := http.NewServeMux()
	mux.Handle("/query", httptransport.NewServer(
		root,
		shttp.EndpointenizeQuery(g),
		shttp.DecodeQueryRequest,
		shttp.EncodeQueryResponse,
	))
	mux.Handle("/thread", httptransport.NewServer(
		root,
		shttp.EndpointenizeThread(g),
		shttp.DecodeThreadRequest,
		shttp.EncodeThreadResponse,
	))
	mux.Handle("/metrics", stdprometheus.Handler())

	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, mux))
}
