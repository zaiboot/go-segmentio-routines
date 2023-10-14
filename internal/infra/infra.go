package infra

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"zaiboot/segmentIO.tests/internal/configs"

	"github.com/rs/zerolog"
)

func MonitorProcesses(l *zerolog.Logger, wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for run := true; run; {
		select {
		case s := <-sigchan:
			l.Info().
				Str("signal", s.String()).
				Msg("Terminated due to syscall")
			run = false
			cancel()
			// app was terminated either via  SIGTERM or docker stopped.
		case <-ctx.Done():
			l.Error().
				Msg("Terminated due to error in ctx")
			run = false

		}
	}

}

func StartReadynessWebServer(c *configs.Config, l *zerolog.Logger, wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/kill_me", func(rw http.ResponseWriter, _ *http.Request) {
		cancel()
		rw.WriteHeader(http.StatusOK)
	})

	s := &http.Server{
		Addr:              ":" + c.PORT,
		Handler:           mux,
		ReadHeaderTimeout: 15 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func(s1 *http.Server) {

		for run := true; run; {
			// TODO: Make this better
			select {
			case <-ctx.Done():
				if err := s1.Close(); err != nil {
					l.Err(err).Msg("Unable to close gracefully")
				}
				l.Error().
					Msg("Terminated due to error in ctx")
				run = false

			}
		}
	}(s)

	err := s.ListenAndServe()
	if err != nil && err.Error() != "http: Server closed" {
		l.Err(err).Msg("Unable to start metrics server")
	}
}
