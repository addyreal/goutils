package h22p

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	Port            string
	Stderr          io.Writer
	Handler         http.HandlerFunc
	HeaderTimeout   time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:              ":" + s.Port,
		Handler:           s.Handler,
		ReadHeaderTimeout: s.HeaderTimeout,
		ReadTimeout:       s.ReadTimeout,
		WriteTimeout:      s.WriteTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		fmt.Fprintln(s.Stderr, "Attempting shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Fprintln(s.Stderr, "Forcing shutdown")
			os.Exit(1)
		}
	}()

	fmt.Fprintln(s.Stderr, "Serving on :"+s.Port)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	fmt.Fprintln(s.Stderr, "Stopped serving on :"+s.Port)
	return nil
}

func GetRealIp(r *http.Request) string {
	s := r.Header.Get("X-Real-IP")
	if s != "" {
		return s
	}

	s = r.Header.Get("X-Forwarded-For")
	if s != "" {
		return strings.TrimSpace(strings.Split(s, ",")[0])
	}

	s, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return s
}
