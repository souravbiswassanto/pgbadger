package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

	// Serve a basic login page which sets an initial csrf cookie
	r.GET("/login", func(c *gin.Context) {
		// generate csrf token and set as cookie (not HttpOnly so client JS can read it)
		// TODO: use a more secure random token generator
		csrfVal := fmt.Sprintf("%d", time.Now().UnixNano())
		secure := cfg.Security.EnableTLS
		csrf := &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfVal,
			Path:     "/",
			HttpOnly: false,
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(time.Hour.Seconds()),
			Expires:  time.Now().Add(time.Hour),
		}
		http.SetCookie(c.Writer, csrf)

		// serve the static login page
		c.File("./web/login.html")
	})
}

func handleReportGeneration(cfg *config.Config, lg *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var opts PgbadgerOptions

		// Read options from query parameters (double submit CSRF enforced in middleware)
		// Helpers to read params
		qp := func(key string) string { return c.Query(key) }
		qps := func(key string) []string { return c.QueryArray(key) }

		if v := qp("format"); v != "" {
			opts.Format = &v
		}
		if v := qp("outfile"); v != "" {
			opts.Outfile = &v
		}
		if v := qp("outdir"); v != "" {
			opts.Outdir = &v
		}
		if v := qp("title"); v != "" {
			opts.Title = &v
		}
		if v := qp("jobs"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				opts.Jobs = &n
			}
		}
		if v := qp("jobs_parallel"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				opts.JobsParallel = &n
			}
		}
		if v := qp("verbose"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.Verbose = &b
			}
		}
		if v := qp("quiet"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.Quiet = &b
			}
		}

		if v := qp("dbname"); v != "" {
			opts.Dbname = &v
		}
		if v := qp("dbuser"); v != "" {
			opts.Dbuser = &v
		}
		if v := qp("appname"); v != "" {
			opts.Appname = &v
		}
		if v := qp("client_host"); v != "" {
			opts.ClientHost = &v
		}

		if v := qp("begin"); v != "" {
			opts.Begin = &v
		}
		if v := qp("end"); v != "" {
			opts.End = &v
		}

		if v := qp("top"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				opts.Top = &n
			}
		}
		if v := qp("sample"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				opts.Sample = &n
			}
		}
		if v := qp("maxlength"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				opts.Maxlength = &n
			}
		}

		if v := qp("average"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				opts.Average = &n
			}
		}
		if v := qp("histo_average"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				opts.HistoAverage = &n
			}
		}
		if v := qp("nograph"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.Nograph = &b
			}
		}

		if v := qp("extension"); v != "" {
			opts.Extension = &v
		}
		if v := qp("prettify"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.Prettify = &b
			}
		}
		if v := qp("query_numbering"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.QueryNumbering = &b
			}
		}

		if v := qp("select_only"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.SelectOnly = &b
			}
		}
		if v := qp("watch_mode"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.WatchMode = &b
			}
		}
		if v := qp("incremental"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.Incremental = &b
			}
		}
		if v := qp("explode"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				opts.Explode = &b
			}
		}

		// list-type params
		if s := qps("exclude_user"); len(s) > 0 {
			opts.ExcludeUser = s
		}
		if s := qps("exclude_appname"); len(s) > 0 {
			opts.ExcludeAppname = s
		}
		if s := qps("exclude_client"); len(s) > 0 {
			opts.ExcludeClient = s
		}
		if s := qps("exclude_db"); len(s) > 0 {
			opts.ExcludeDb = s
		}

		if s := qps("include_query"); len(s) > 0 {
			opts.IncludeQuery = s
		}
		if s := qps("include_pid"); len(s) > 0 {
			opts.IncludePid = s
		}
		if s := qps("include_session"); len(s) > 0 {
			opts.IncludeSession = s
		}

		if v := qp("prefix"); v != "" {
			opts.Prefix = &v
		}
		if v := qp("ident"); v != "" {
			opts.Ident = &v
		}
		if v := qp("timezone"); v != "" {
			opts.Timezone = &v
		}
		if v := qp("log_timezone"); v != "" {
			opts.LogTimezone = &v
		}

		if v := qp("data_dir"); v != "" {
			opts.DataDir = &v
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

		// Determine log file path and expand glob
		logGlob := filepath.Join(*opts.DataDir, "*.log")
		files, err := filepath.Glob(logGlob)
		if err != nil || len(files) == 0 {
			lg.Errorw("No log files found for pattern", "pattern", logGlob, "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "no log files found"})
			return
		}

		// Build the pgbadger command; pass all file paths explicitly
		args := opts.BuildCommand(files[0])
		if len(files) > 1 {
			// append remaining files as positional args
			for _, f := range files[1:] {
				args = append(args, f)
			}
		}

		// Set a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, "pgbadger", args...)

		// Log the command being executed (for debugging)
		lg.Infow("Executing pgbadger command",
			"command", fmt.Sprintf("pgbadger %v", args),
			"user", c.GetString("username"))

		// Execute the command and capture output
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

		// Write output to a temp file and serve it
		tmpFile, err := os.CreateTemp("/tmp", "pgbadger-*.out")
		if err != nil {
			lg.Errorw("Failed to create temp file", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp file"})
			return
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(output); err != nil {
			lg.Errorw("Failed to write report to temp file", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write report"})
			return
		}
		_ = tmpFile.Close()

		// Determine content type
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

		data, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			lg.Errorw("Failed to read temp report file", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read report"})
			return
		}

		c.Header("Content-Type", contentType)
		c.String(http.StatusOK, string(data))

		lg.Infow("Report generated successfully",
			"temp_file", tmpFile.Name(),
			"content_type", contentType,
			"user", c.GetString("username"))
	}
}
