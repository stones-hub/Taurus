package app

import (
	"Taurus/config"
	"Taurus/pkg/mcp"
	"Taurus/pkg/router"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ANSI escape sequences define colors
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

// DefaultHost and DefaultPort are the default server address and port
var (
	env        = ".env.local"
	configPath = "./config"
)

// Default initializes and starts the HTTP server with default settings
func Default() {
	Start(config.Core.AppHost, config.Core.AppPort)
}

// Start initializes and starts the HTTP server with graceful shutdown
func Start(host string, port int) {
	// Load routes
	r := router.LoadRoutes()

	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr:        addr,
		Handler:     r,
		IdleTimeout: 1 * time.Minute,
	}

	// use errChan to receive http server startup error
	errChan := make(chan error, 1)

	// Run server in a goroutine
	go func() {
		log.Printf("%s🔗 -> Server is running on %s %s \n", Green, addr, Reset)
		// when server startup failed, write error to errChan.
		// But http.ErrServerClosed is not an error,,because it is expected when the server is closed.
		// ListenAndServe is a blocking call
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// Block until a signal is received or an error is returned.
	// If an error is returned, it is a fatal error and the program will exit.
	if err := signalWaiter(errChan); err != nil {
		log.Fatalf("%sServer startup failed: %v %s\n", Red, err, Reset)
	}

	// If signalWaiter returns nil, it means the server is running. But received a signal, so we need to shutdown the server.

	// Create a deadline to wait for, 5 seconds or cancel() are all called ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if mcp.GlobalMCPServer != nil {
		srv.RegisterOnShutdown(func() {
			if err := mcp.GlobalMCPServer.Shutdown(ctx); err != nil {
				log.Printf("%sMCP server shutdown failed: %v %s\n", Red, err, Reset)
			}
		})
	}

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("%sServer forced to shutdown: %v %s\n", Red, err, Reset)
	}

	log.Printf("%s🔗 -> Server shutdown successfully. %s\n", Green, Reset)

	gracefulCleanup(ctx)
}

// signalWaiter waits for a signal or an error, then return
func signalWaiter(errCh chan error) error {
	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)

	// Block until a signal is received or an error is returned
	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			log.Printf("%s🔗 -> Received signal: %s, graceful shutdown... %s\n", Yellow, sig, Reset)
			// graceful shutdown
			return nil
		}
	case err := <-errCh:
		return err
	}

	return nil
}

// gracefulCleanup is called when the server is shutting down. we can do some cleanup work here.
func gracefulCleanup(ctx context.Context) {

	log.Printf("%s🔗 -> Waiting for all requests to be processed... %s\n", Yellow, Reset)
	done := make(chan struct{})

	go func() {
		Cleanup()
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Printf("%s🔗 -> Server stopped successfully. %s\n", Green, Reset)
	case <-ctx.Done():
		// If 5 seconds have passed and the server has not stopped, it means the server is not responding, so we need to force it to stop.
		log.Printf("%s🔗 -> Server stopped forcefully. %s\n", Red, Reset)
	}
}

// init is automatically called before the main function
// --env .env.local --config ./config
func init() {
	// custom usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s\n", Cyan+"==================== Usage ===================="+Reset)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s-e, --env <file>%s      Specify the environment file (default \".env.local\")\n", Green, Reset)
		fmt.Fprintf(os.Stderr, "  %s-c, --config <path>%s   Specify the configuration file or directory (default \"config\")\n", Green, Reset)
		fmt.Fprintf(os.Stderr, "  %s-h, --help%s            Show this help message\n", Green, Reset)
		fmt.Fprintf(os.Stderr, "%s\n", Cyan+"==============================================="+Reset)
	}

	// set command line arguments and their aliases
	flag.StringVar(&env, "env", ".env.local", "Environment file")
	flag.StringVar(&env, "e", ".env.local", "Environment file (alias)")
	flag.StringVar(&configPath, "config", "config", "Path to the configuration file or directory")
	flag.StringVar(&configPath, "c", "config", "Path to the configuration file or directory (alias)")

	// parse command line arguments
	flag.Parse()

	// initialize all modules.
	// the env file is not needed, because the makefile has already written the environment variables into the env file, but for the sake of rigor, we still pass the env file to the initialize function
	initialize(configPath, env)
}
