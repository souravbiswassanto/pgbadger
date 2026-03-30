package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/souravbiswassanto/pgbadger/config"
	"go.uber.org/zap"
)

func Register(r *gin.Engine, cfg *config.Config, lg *zap.SugaredLogger) {
	authHandler := NewAuthHandler(cfg, lg)

	// Apply auth middleware to all API routes
	api := r.Group("/api")
	api.Use(authHandler.AuthMiddleware())

	v1 := api.Group("/v1")
	v1.POST("/login", authHandler.Login)
	v1.POST("/report", handleReportGeneration(cfg, lg))

	// Health check (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func handleReportGeneration(cfg *config.Config, lg *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var opts PgbadgerOptions

		// Bind JSON request body to options struct
		if err := c.ShouldBindJSON(&opts); err != nil {
			lg.Errorw("Invalid request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Set defaults
		opts.SetDefaults()

		// Validate options
		if err := opts.Validate(); err != nil {
			lg.Errorw("Invalid pgbadger options", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid options",
				"details": err.Error(),
			})
			return
		}

		// Determine log file path
		logPath := filepath.Join(*opts.DataDir, "*.log")

		// Build the pgbadger command
		args := opts.BuildCommand(logPath)
		cmd := exec.Command("pgbadger", args...)

		// Set a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		cmd = exec.CommandContext(ctx, "pgbadger", args...)

		// Log the command being executed (for debugging)
		lg.Infow("Executing pgbadger command",
			"command", fmt.Sprintf("pgbadger %v", args),
			"user", c.GetString("username"))

		// Execute the command
		output, err := cmd.CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			lg.Errorw("pgbadger command timed out", "command", fmt.Sprintf("pgbadger %v", args))
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Command timed out",
				"timeout": "5 minutes",
			})
			return
		}

		if err != nil {
			lg.Errorw("pgbadger command failed",
				"error", err,
				"output", string(output),
				"command", fmt.Sprintf("pgbadger %v", args))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate report",
				"details": err.Error(),
				"output":  string(output),
			})
			return
		}

		// Determine output file path
		outputFile := "/tmp/report.html"
		if opts.Outfile != nil {
			outputFile = *opts.Outfile
		} else if opts.Extension != nil && *opts.Extension == "json" {
			outputFile = "/tmp/report.json"
		} else if opts.Extension != nil && *opts.Extension == "text" {
			outputFile = "/tmp/report.txt"
		}

		// Read the generated report file
		htmlContent, err := os.ReadFile(outputFile)
		if err != nil {
			lg.Errorw("Failed to read report file", "error", err, "file", outputFile)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read generated report",
				"details": err.Error(),
			})
			return
		}

		// Set appropriate content type based on extension
		contentType := "text/html"
		if opts.Extension != nil {
			switch *opts.Extension {
			case "json":
				contentType = "application/json"
			case "text":
				contentType = "text/plain"
			case "bin":
				contentType = "application/octet-stream"
			}
		}

		// Return the report
		c.Header("Content-Type", contentType)
		c.String(http.StatusOK, string(htmlContent))

		lg.Infow("Report generated successfully",
			"output_file", outputFile,
			"content_type", contentType,
			"user", c.GetString("username"))
	}
}
