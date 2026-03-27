package handler

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(r *gin.Engine, lg *zap.SugaredLogger) {
	api := r.Group("/api")

	v1 := api.Group("/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1.GET("/report", handleReportGeneration)

	// placeholder for more handlers
	_ = lg
}

func handleReportGeneration(c *gin.Context) {
	// Use shell to expand the glob pattern for log files
	cmd := exec.Command("sh", "-c", "pgbadger /var/pv/data/log/*.log -o /tmp/report.html -f stderr --jobs 2 --verbose")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"output":  string(output),
			"success": false,
		})
		return
	}

	// Read the generated report.html file
	htmlContent, err := os.ReadFile("/tmp/report.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	// Set content type to HTML and return the report
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, string(htmlContent))
}
