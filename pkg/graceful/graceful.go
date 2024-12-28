package graceful

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// Followings are the known constants.
const (
	defaultShutdownTimeout = 10 * time.Second
)

// Followings are the known errors from go-lib/graceful.
var (
	// errShutdownTimeoutExceeded happens when the process of
	// gracefully shutdown the server exceeding timeout limit.
	errShutdownTimeoutExceeded = errors.New("graceful: shutdown timeout exceeded")
)

// ServeHTTP starts the given HTTP server on the given address
// and add graceful shutdown handler.
//
// timeout specify how long to wait for the graceful shutdown
// handler to run. if timeout = 0, default value of 10 second
// will be used.
func ServeHTTP(server *http.Server, address string, timeout time.Duration) error {
	// start listener
	lis, err := listen(address)
	if err != nil {
		return err
	}

	// wait for and handle termination signal
	exit := wait(func() error {
		if timeout == 0 {
			timeout = defaultShutdownTimeout
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// graceful shutdown
		done := make(chan struct{})
		go func() {
			server.Shutdown(ctx)
			close(done)
		}()

		select {
		case <-ctx.Done():
			return errShutdownTimeoutExceeded
		case <-done:
		}

		return nil
	})

	// start serving
	log.Printf("graceful: http server running on address: %v\n", address)
	err = server.Serve(lis) // always return a non-nil error
	if err != http.ErrServerClosed {
		return err
	}

	// wait until exit channel is unblocked
	<-exit

	log.Println("graceful: http server stopped")
	return nil
}

// ServeHTTP starts the given gRPC server on the given address
// and add graceful shutdown handler.
//
// timeout specify how long to wait for the graceful shutdown
// handler to run. if timeout = 0, default value of 10 second
// will be used.
func ServeGRPC(server *grpc.Server, address string, timeout time.Duration) error {
	// start listener
	lis, err := listen(address)
	if err != nil {
		return err
	}

	// wait for and handle termination signal
	exit := wait(func() error {
		if timeout == 0 {
			timeout = defaultShutdownTimeout
		}

		// graceful shutdown
		done := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(done)
		}()

		select {
		case <-time.After(timeout):
			// force shutdown
			server.Stop()
			return errShutdownTimeoutExceeded
		case <-done:
		}

		return nil
	})

	// start serving
	log.Printf("graceful: grpc server running on address: %v\n", address)
	err = server.Serve(lis)
	if err != nil {
		return err
	}

	// wait until exit channel is unblocked
	<-exit

	log.Println("graceful: grpc server stopped")
	return nil
}

// listen listens to the given address.
func listen(address string) (net.Listener, error) {
	return net.Listen("tcp4", address)
}

// wait waits for termination signal.
//
// wait executes the given onSignal after receiving signal.
// onSignal should define how to shutdown the server, thus
// server can properly shutdown upon receiving termination
// signal.
//
// wait returns a read-only channel which will be closed
// after the signal received and onSignal executed.
func wait(onSignal func() error) <-chan struct{} {
	exit := make(chan struct{})
	go func() {
		// wait for termination signal
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-sig

		// execute onSignal
		err := onSignal()
		if err != nil {
			log.Printf("graceful: shutdown failed. err: %v\n", err)
		} else {
			log.Println("graceful: shutdown succeed")
		}

		// unblock exit channel
		close(exit)
	}()
	return exit
}
