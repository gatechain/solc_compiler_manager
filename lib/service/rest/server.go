package rest

import (
	"context"
	"fmt"
	"github.com/gatechain/smart_contract_verifier/lib"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

// RestServer represents the Light Client Rest server
type Server struct {
	Mux 		*mux.Router
	Cdc 		*lib.Codec
	CliCtx 		context.Context
	log			log.Logger
	listener 	net.Listener
}

// NewRestServer creates a new rest server instance
func NewRestServer(cdc *lib.Codec) *Server {
	r := mux.NewRouter()
	cliCtx := context.Background()
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	return &Server{
		Mux:    	r,
		Cdc: 		cdc,
		CliCtx: 	cliCtx,
		log:    	logger,
	}
}

// Start starts the rest server
func (rs *Server) Start(listenAddr string, maxOpen int, readTimeout, writeTimeout uint, cors bool) (err error) {
	TrapSignal(func() {
		err := rs.listener.Close()
		_ = rs.log.Log("error closing listener", "err", err)
	})

	cfg := DefaultConfig()
	cfg.MaxOpenConnections = maxOpen
	cfg.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.WriteTimeout = time.Duration(writeTimeout) * time.Second

	rs.listener, err = Listen(listenAddr, cfg)
	if err != nil {
		return
	}
	_ = rs.log.Log(
		"INFO_START_REST_SERVER",
		fmt.Sprintf(
			"Starting application REST service at address: %s",
			listenAddr,
		),
	)

	var h http.Handler = rs.Mux
	if cors {
		allowAllCORS := handlers.CORS(handlers.AllowedHeaders([]string{"Content-Type"}))
		h = allowAllCORS(h)
	}

	return Serve(rs.listener, h, rs.log, cfg)
}

// RegisterRestServerFlags registers the flags required for rest server
func RegisterRestServerFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(FlagListenAddr, "tcp://localhost:1212", "The address for the server to listen on")
	cmd.Flags().Uint(FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Uint(FlagRPCReadTimeout, 10, "The RPC read timeout (in seconds)")
	cmd.Flags().Uint(FlagRPCWriteTimeout, 10, "The RPC write timeout (in seconds)")
	cmd.Flags().Bool(FlagUnsafeCORS, false, "Allows CORS requests from all domains. For development purposes only, use it at your own risk.")

	return cmd
}


// TrapSignal traps SIGINT and SIGTERM and terminates the server correctly.
func TrapSignal(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <- sigs
		if cleanupFunc != nil {
			cleanupFunc()
		}
		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}
		os.Exit(exitCode)
	}()
}

// Listen starts a new net.Listener on the given address.
// It returns an error if the address is invalid or the call to Listen() fails.
func Listen(addr string, config *Config) (listener net.Listener, err error) {
	parts := strings.SplitN(addr, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"invalid listening address %s (use fully formed addresses, including the tcp:// or unix:// prefix)",
			addr,
		)
	}
	proto, addr := parts[0], parts[1]
	listener, err = net.Listen(proto, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %v: %v", addr, err)
	}
	if config.MaxOpenConnections > 0 {
		listener = netutil.LimitListener(listener, config.MaxOpenConnections)
	}

	return listener, nil
}

// Serve creates a http.Server and calls Serve with the given listener. It
// wraps handler with RecoverAndLogHandler and a handler, which limits the max
// body size to config.MaxBodyBytes.
//
// NOTE: This function blocks - you may want to call it in a go-routine.
func Serve(listener net.Listener, handler http.Handler, logger log.Logger, config *Config) error {
	_ = logger.Log(
		"INFO_START_REST_SERVER",
		fmt.Sprintf("Start RPC HTTP server on %s", listener.Addr()),
	)

	s := &http.Server{
		Handler:        RecoverAndLogHandler(maxBytesHandler{h: handler, n: config.MaxBodyBytes}, logger),
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	err := s.Serve(listener)
	_ = logger.Log("RPC HTTP server stopped", "err", err)
	return err
}


// RecoverAndLogHandler wraps an HTTP handler, adding error logging.
// If the inner function panics, the outer function recovers, logs, sends an
// HTTP 500 error response.
func RecoverAndLogHandler(handler http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter to remember the status
		rww := &responseWriterWrapper{-1, w}
		begin := time.Now()

		rww.Header().Set("X-Server-Time", fmt.Sprintf("%v", begin.Unix()))

		defer func() {
			// Handle any panics in the panic handler below. Does not use the logger, since we want
			// to avoid any further panics. However, we try to return a 500, since it otherwise
			// defaults to 200 and there is no other way to terminate the connection. If that
			// should panic for whatever reason then the Go HTTP server will handle it and
			// terminate the connection - panicing is the de-facto and only way to get the Go HTTP
			// server to terminate the request and close the connection/stream:
			// https://github.com/golang/go/issues/17790#issuecomment-258481416
			if e := recover(); e != nil {
				fmt.Fprintf(os.Stderr, "Panic during RPC panic recovery: %v\n%v\n", e, string(debug.Stack()))
				w.WriteHeader(500)
			}
		}()

		defer func() {
			// Send a 500 error if a panic happens during a handler.
			// Without this, Chrome & Firefox were retrying aborted ajax requests,
			// at least to my localhost.
			if e := recover(); e != nil {
				panic(e)
			}

			// Finally, log.
			durationMS := time.Since(begin).Nanoseconds() / 1000000
			if rww.Status == -1 {
				rww.Status = 200
			}
			_ = logger.Log("Served RPC HTTP response",
				"method", r.Method, "url", r.URL,
				"status", rww.Status, "duration", durationMS,
				"remoteAddr", r.RemoteAddr,
			)
		}()

		handler.ServeHTTP(rww, r)
	})
}
