package app

import (
	"Taurus/pkg/router" // 替换为实际的包路径
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DefaultHost and DefaultPort are the default server address and port
var (
	host = "0.0.0.0"
	port = 8080
)

// Default initializes and starts the HTTP server with default settings
func Default() {
	Start(host, port)
}

// Start initializes and starts the HTTP server with graceful shutdown
func Start(host string, port int) {
	// Load routes
	r := router.LoadRoutes()

	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		log.Printf("Server is running on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", addr, err)
		}
	}()

	log.Printf("Server is ready to handle requests at %s\n", addr)

	// Block until a signal is received
	<-stop

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped successfully.")

}

/*
init 函数的特点
自动调用：init 函数会在程序运行时自动调用，无需显式调用。
执行顺序：
每个包可以有多个 init 函数（包括在同一个文件或不同文件中）。
init 函数的执行顺序是按照包的依赖关系，从被依赖的包到依赖包依次执行。
在同一个包中，init 函数按照文件的编译顺序执行。
无参数和返回值：init 函数不能接受参数，也不能返回值。
*/

// init is automatically called before the main function
func init() {
	// Parse command-line arguments
	configPath := flag.String("config", "config", "Path to the configuration file or directory")
	flag.StringVar(&host, "host", "0.0.0.0", "Host to listen on")
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()

	// Initialize all modules
	Initialize(*configPath)
}
