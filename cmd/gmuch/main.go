package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	stdlog "log"

	"github.com/gmuch/gmuch"
	"github.com/gmuch/gmuch/server"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	"github.com/namsral/flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	sgrpc "github.com/gmuch/gmuch/server/grpc"
	shttp "github.com/gmuch/gmuch/server/http"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	dbPath    = flag.String("db-path", "", "The path to the folder with a '.notmuch' in it")
	debugAddr = flag.String("debug-addr", ":8000", "Address for HTTP debug/instrumentation server")
	httpAddr  = flag.String("http-addr", ":8001", "Address for HTTP (JSON) server")
	grpcAddr  = flag.String("grpc-addr", ":8002", "Address for gRPC server")
)

func main() {
	flag.Parse()

	// package log
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC).With("caller", log.DefaultCaller)
		stdlog.SetFlags(0)                             // flags are handled by Go kit's logger
		stdlog.SetOutput(log.NewStdlibAdapter(logger)) // redirect anything using stdlib log to us
	}

	// package metrics
	var (
		requestCount   metrics.Counter
		requestLatency metrics.TimeHistogram
	)
	{
		fieldKeys := []string{"method", "error"}
		requestCount = kitprometheus.NewCounter(stdprometheus.CounterOpts{
			Namespace: "gmuch",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
		requestLatency = metrics.NewTimeHistogram(time.Nanosecond, metrics.NewMultiHistogram(
			expvar.NewHistogram("request_duration_ns", 0, 5e9, 1, 50, 95, 99),
			kitprometheus.NewSummary(stdprometheus.SummaryOpts{
				Namespace: "gmuch",
				Subsystem: "api",
				Name:      "duration_ns",
				Help:      "Request duration in nanoseconds.",
			}, fieldKeys),
		))
	}

	// Business domain
	var g server.GmuchService
	{
		g = gmuch.New(*dbPath, logger)
		g = server.LoggingMiddleware(logger)(g)
		g = server.InstrumentingMiddleware(requestCount, requestLatency)(g)
	}

	// Mechanical stuff
	rand.Seed(time.Now().UnixNano())
	root := context.Background()
	errc := make(chan error)

	go func() {
		errc <- interrupt()
	}()

	// Debug/instrumentation
	go func() {
		transportLogger := log.NewContext(logger).With("transport", "debug")
		_ = transportLogger.Log("addr", *debugAddr)
		errc <- http.ListenAndServe(*debugAddr, nil) // DefaultServeMux
	}()

	// Transport: HTTP/JSON
	go func() {
		transportLogger := log.NewContext(logger).With("transport", "HTTP/JSON")
		mux := http.NewServeMux()

		mux.Handle("/query", httptransport.NewServer(
			root,
			shttp.EndpointenizeQuery(g),
			shttp.DecodeQueryRequest,
			shttp.EncodeQueryResponse,
			httptransport.ServerErrorLogger(transportLogger),
		))
		mux.Handle("/thread", httptransport.NewServer(
			root,
			shttp.EndpointenizeThread(g),
			shttp.DecodeThreadRequest,
			shttp.EncodeThreadResponse,
			httptransport.ServerErrorLogger(transportLogger),
		))

		_ = transportLogger.Log("addr", *httpAddr)
		errc <- http.ListenAndServe(*httpAddr, mux)
	}()

	// Transport: gRPC
	go func() {
		transportLogger := log.NewContext(logger).With("transport", "gRPC")
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}
		s := grpc.NewServer() // uses its own, internal context
		sgrpc.RegisterGmuchServer(s, sgrpc.Binding{g})
		_ = transportLogger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	_ = logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
