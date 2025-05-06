package app

import (
	"Taurus/config"
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

// ANSI 转义序列定义颜色
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
	Start(config.AppConfig.AppHost, config.AppConfig.AppPort)
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
		log.Printf("%sServer is running on %s %s \n", Green, addr, Reset)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%sCould not listen on %s: %v %s\n", Red, addr, err, Reset)
		}
	}()

	// Block until a signal is received
	<-stop

	// Create a deadline to wait for, 5 seconds or cancel() are all called ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%sServer forced to shutdown: %v %s\n", Red, err, Reset)
	}

	log.Printf("%sWaiting for all requests to be processed... %s\n", Yellow, Reset)
	done := make(chan struct{})

	go func() {
		Cleanup()
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Printf("%sServer stopped successfully. %s\n", Green, Reset)
	case <-ctx.Done():
		// 如果5秒内没有处理完，则强制关闭
		log.Printf("%sServer stopped forcefully. %s\n", Red, Reset)
	}
}

// init is automatically called before the main function
// --env .env.local --config ./config
func init() {
	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s\n", Cyan+"==================== Usage ===================="+Reset)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s-e, --env <file>%s      Specify the environment file (default \".env.local\")\n", Green, Reset)
		fmt.Fprintf(os.Stderr, "  %s-c, --config <path>%s   Specify the configuration file or directory (default \"config\")\n", Green, Reset)
		fmt.Fprintf(os.Stderr, "  %s-h, --help%s            Show this help message\n", Green, Reset)
		fmt.Fprintf(os.Stderr, "%s\n", Cyan+"==============================================="+Reset)
	}

	// 设置命令行参数及其别名
	flag.StringVar(&env, "env", ".env.local", "Environment file")
	flag.StringVar(&env, "e", ".env.local", "Environment file (alias)")
	flag.StringVar(&configPath, "config", "config", "Path to the configuration file or directory")
	flag.StringVar(&configPath, "c", "config", "Path to the configuration file or directory (alias)")

	// 解析命令行参数
	flag.Parse()

	// Initialize all modules, 其实env传不传无所谓，因为makefile中已经将环境变量写入了， 但是为了严谨还是将envifle传入
	initialize(configPath, env)
}
