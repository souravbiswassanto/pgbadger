package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
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

	// Serve login page (sets a crypto-random csrf cookie readable by JS for double-submit CSRF)
	r.GET("/login", func(c *gin.Context) {
		// generate crypto-random csrf token and set as cookie (not HttpOnly so client JS can read it)
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			lg.Errorw("failed to generate csrf token", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		csrfVal := base64.RawURLEncoding.EncodeToString(b)
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

		// Render login template (html/template escapes by default)
		tplPath, err := findTemplateFile("web/login.tmpl")
		if err != nil {
			lg.Errorw("failed to find login template", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles(tplPath)
		if err != nil {
			lg.Errorw("failed to parse login template", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			lg.Errorw("failed to execute login template", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	})
}

// findTemplateFile attempts to locate a template file by trying a few relative
// paths upwards from the current working directory. This helps tests running
// from package directories find project-level assets.
func findTemplateFile(rel string) (string, error) {
	// candidates: cwd/rel, parent/rel, parent/parent/rel, etc.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// try up to 4 levels
	p := cwd
	for i := 0; i < 5; i++ {
		cand := filepath.Join(p, rel)
		if _, err := os.Stat(cand); err == nil {
			return cand, nil
		}
		p = filepath.Dir(p)
	}
	return "", fmt.Errorf("template %s not found in cwd or parent dirs", rel)
}

func handleReportGeneration(cfg *config.Config, lg *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var opts PgbadgerOptions

		// Read options from query parameters (double submit CSRF enforced in middleware)
		// Helpers to read params
		qp := func(key string) string { return c.Query(key) }
		qps := func(key string) []string { return c.QueryArray(key) }

		if v := qp("format"); v != "" {
			v = SanitizeString(v, 128)
			opts.Format = &v
		}
		if v := qp("outfile"); v != "" {
			vv := SanitizeString(v, 512)
			opts.Outfile = &vv
		}
		if v := qp("outdir"); v != "" {
			vv := SanitizeString(v, 512)
			opts.Outdir = &vv
		}
		if v := qp("title"); v != "" {
			vv := SanitizeString(v, 256)
			opts.Title = &vv
		}
		if v := qp("jobs"); v != "" {
			if n, err := strconv.Atoi(SanitizeString(v, 10)); err == nil {
				opts.Jobs = &n
			}
		}
		if v := qp("jobs_parallel"); v != "" {
			if n, err := strconv.Atoi(SanitizeString(v, 10)); err == nil {
				opts.JobsParallel = &n
			}
		}
		if v := qp("verbose"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.Verbose = &b
			}
		}
		if v := qp("quiet"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.Quiet = &b
			}
		}

		if v := qp("dbname"); v != "" {
			vv := SanitizeString(v, 128)
			opts.Dbname = &vv
		}
		if v := qp("dbuser"); v != "" {
			vv := SanitizeString(v, 128)
			opts.Dbuser = &vv
		}
		if v := qp("appname"); v != "" {
			vv := SanitizeString(v, 128)
			opts.Appname = &vv
		}
		if v := qp("client_host"); v != "" {
			vv := SanitizeString(v, 128)
			opts.ClientHost = &vv
		}

		if v := qp("begin"); v != "" {
			vv := SanitizeString(v, 64)
			opts.Begin = &vv
		}
		if v := qp("end"); v != "" {
			vv := SanitizeString(v, 64)
			opts.End = &vv
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
			vv := SanitizeString(v, 16)
			opts.Extension = &vv
		}
		if v := qp("prettify"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.Prettify = &b
			}
		}
		if v := qp("query_numbering"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.QueryNumbering = &b
			}
		}

		if v := qp("select_only"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.SelectOnly = &b
			}
		}
		if v := qp("watch_mode"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.WatchMode = &b
			}
		}
		if v := qp("incremental"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.Incremental = &b
			}
		}
		if v := qp("explode"); v != "" {
			if b, err := strconv.ParseBool(SanitizeString(v, 6)); err == nil {
				opts.Explode = &b
			}
		}

		// list-type params (sanitized)
		if s := qps("exclude_user"); len(s) > 0 {
			opts.ExcludeUser = SanitizeStringSlice(s, 128)
		}
		if s := qps("exclude_appname"); len(s) > 0 {
			opts.ExcludeAppname = SanitizeStringSlice(s, 128)
		}
		if s := qps("exclude_client"); len(s) > 0 {
			opts.ExcludeClient = SanitizeStringSlice(s, 128)
		}
		if s := qps("exclude_db"); len(s) > 0 {
			opts.ExcludeDb = SanitizeStringSlice(s, 128)
		}

		if s := qps("include_query"); len(s) > 0 {
			opts.IncludeQuery = SanitizeStringSlice(s, 512)
		}
		if s := qps("include_pid"); len(s) > 0 {
			opts.IncludePid = SanitizeStringSlice(s, 64)
		}
		if s := qps("include_session"); len(s) > 0 {
			opts.IncludeSession = SanitizeStringSlice(s, 128)
		}

		if v := qp("prefix"); v != "" {
			vv := SanitizeString(v, 128)
			opts.Prefix = &vv
		}
		if v := qp("ident"); v != "" {
			vv := SanitizeString(v, 128)
			opts.Ident = &vv
		}
		if v := qp("timezone"); v != "" {
			vv := SanitizeString(v, 64)
			opts.Timezone = &vv
		}
		if v := qp("log_timezone"); v != "" {
			vv := SanitizeString(v, 64)
			opts.LogTimezone = &vv
		}

		if v := qp("data_dir"); v != "" {
			vv := SanitizeString(v, 256)
			opts.DataDir = &vv
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
