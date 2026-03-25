package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/souravbiswassanto/pgbadger/config"
	"github.com/souravbiswassanto/pgbadger/internal/server"
	"github.com/souravbiswassanto/pgbadger/pkg/logger"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		// init global logger
		lg := logger.New(cfg.Log.Level)
		defer lg.Sync()

		srv := server.New(cfg, lg)

		// start server in goroutine
		go func() {
			if err := srv.Start(); err != nil {
				lg.Fatal("server failed to start", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			lg.Fatal("server forced to shutdown", err)
		}

		fmt.Println("server exiting")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
