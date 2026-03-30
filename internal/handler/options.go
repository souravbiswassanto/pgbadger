package handler

import (
	"fmt"
	"strconv"
	"time"
)

// PgbadgerOptions represents the allowed pgbadger command-line options
// Only whitelisted options are included for security
type PgbadgerOptions struct {
	// Basic options
	Format       *string `json:"format,omitempty"`        // -f, --format
	Outfile      *string `json:"outfile,omitempty"`       // -o, --outfile
	Outdir       *string `json:"outdir,omitempty"`        // -O, --outdir
	Title        *string `json:"title,omitempty"`         // -T, --title
	Jobs         *int    `json:"jobs,omitempty"`          // -j, --jobs
	JobsParallel *int    `json:"jobs_parallel,omitempty"` // -J, --Jobs
	Verbose      *bool   `json:"verbose,omitempty"`       // -v, --verbose
	Quiet        *bool   `json:"quiet,omitempty"`         // -q, --quiet

	// Filtering options
	Dbname     *string `json:"dbname,omitempty"`      // -d, --dbname
	Dbuser     *string `json:"dbuser,omitempty"`      // -u, --dbuser
	Appname    *string `json:"appname,omitempty"`     // -N, --appname
	ClientHost *string `json:"client_host,omitempty"` // -c, --dbclient

	// Time options
	Begin *string `json:"begin,omitempty"` // -b, --begin
	End   *string `json:"end,omitempty"`   // -e, --end

	// Report options
	Top       *int `json:"top,omitempty"`       // -t, --top
	Sample    *int `json:"sample,omitempty"`    // -s, --sample
	Maxlength *int `json:"maxlength,omitempty"` // -m, --maxlength

	// Graph options
	Average      *int  `json:"average,omitempty"`       // -a, --average
	HistoAverage *int  `json:"histo_average,omitempty"` // -A, --histo-average
	Nograph      *bool `json:"nograph,omitempty"`       // -G, --nograph

	// Output options
	Extension      *string `json:"extension,omitempty"`       // -x, --extension
	Prettify       *bool   `json:"prettify,omitempty"`        // -P, --no-prettify (inverted)
	QueryNumbering *bool   `json:"query_numbering,omitempty"` // -Q, --query-numbering

	// Special modes
	SelectOnly  *bool `json:"select_only,omitempty"` // -S, --select-only
	WatchMode   *bool `json:"watch_mode,omitempty"`  // -w, --watch-mode
	Incremental *bool `json:"incremental,omitempty"` // -I, --incremental
	Explode     *bool `json:"explode,omitempty"`     // -E, --explode

	// Exclude options
	ExcludeUser    []string `json:"exclude_user,omitempty"`    // -U, --exclude-user
	ExcludeAppname []string `json:"exclude_appname,omitempty"` // --exclude-appname
	ExcludeClient  []string `json:"exclude_client,omitempty"`  // --exclude-client
	ExcludeDb      []string `json:"exclude_db,omitempty"`      // --exclude-db

	// Include options
	IncludeQuery   []string `json:"include_query,omitempty"`   // --include-query
	IncludePid     []string `json:"include_pid,omitempty"`     // --include-pid
	IncludeSession []string `json:"include_session,omitempty"` // --include-session

	// Advanced options
	Prefix      *string `json:"prefix,omitempty"`       // -p, --prefix
	Ident       *string `json:"ident,omitempty"`        // -i, --ident
	Timezone    *string `json:"timezone,omitempty"`     // -Z, --timezone
	LogTimezone *string `json:"log_timezone,omitempty"` // --log-timezone

	// Data directory (default if not provided)
	DataDir *string `json:"data_dir,omitempty"`
}

// Validate checks if the provided options are valid
func (opts *PgbadgerOptions) Validate() error {
	// Validate format
	if opts.Format != nil {
		validFormats := map[string]bool{
			"syslog": true, "syslog2": true, "stderr": true, "jsonlog": true,
			"csv": true, "pgbouncer": true, "logplex": true, "rds": true, "redshift": true,
		}
		if !validFormats[*opts.Format] {
			return fmt.Errorf("invalid format: %s", *opts.Format)
		}
	}

	// Validate extension
	if opts.Extension != nil {
		validExtensions := map[string]bool{
			"text": true, "html": true, "bin": true, "json": true,
		}
		if !validExtensions[*opts.Extension] {
			return fmt.Errorf("invalid extension: %s", *opts.Extension)
		}
	}

	// Validate jobs
	if opts.Jobs != nil && (*opts.Jobs < 1 || *opts.Jobs > 32) {
		return fmt.Errorf("jobs must be between 1 and 32")
	}

	if opts.JobsParallel != nil && (*opts.JobsParallel < 1 || *opts.JobsParallel > 32) {
		return fmt.Errorf("jobs_parallel must be between 1 and 32")
	}

	// Validate top
	if opts.Top != nil && (*opts.Top < 1 || *opts.Top > 1000) {
		return fmt.Errorf("top must be between 1 and 1000")
	}

	// Validate sample
	if opts.Sample != nil && (*opts.Sample < 1 || *opts.Sample > 100) {
		return fmt.Errorf("sample must be between 1 and 100")
	}

	// Validate maxlength
	if opts.Maxlength != nil && (*opts.Maxlength < 100 || *opts.Maxlength > 1000000) {
		return fmt.Errorf("maxlength must be between 100 and 1000000")
	}

	// Validate average
	if opts.Average != nil && (*opts.Average < 1 || *opts.Average > 1440) {
		return fmt.Errorf("average must be between 1 and 1440 minutes")
	}

	// Validate histo_average
	if opts.HistoAverage != nil && (*opts.HistoAverage < 1 || *opts.HistoAverage > 10080) {
		return fmt.Errorf("histo_average must be between 1 and 10080 minutes")
	}

	// Validate time formats if provided
	if opts.Begin != nil {
		if _, err := time.Parse(time.RFC3339, *opts.Begin); err != nil {
			if _, err := time.Parse("2006-01-02 15:04:05", *opts.Begin); err != nil {
				return fmt.Errorf("invalid begin time format, use RFC3339 or '2006-01-02 15:04:05'")
			}
		}
	}

	if opts.End != nil {
		if _, err := time.Parse(time.RFC3339, *opts.End); err != nil {
			if _, err := time.Parse("2006-01-02 15:04:05", *opts.End); err != nil {
				return fmt.Errorf("invalid end time format, use RFC3339 or '2006-01-02 15:04:05'")
			}
		}
	}

	return nil
}

// BuildCommand builds the pgbadger command with the provided options
func (opts *PgbadgerOptions) BuildCommand(logPath string) []string {
	args := []string{}

	// Basic options
	if opts.Format != nil {
		args = append(args, "-f", *opts.Format)
	}
	if opts.Outfile != nil {
		args = append(args, "-o", *opts.Outfile)
	}
	if opts.Outdir != nil {
		args = append(args, "-O", *opts.Outdir)
	}
	if opts.Title != nil {
		args = append(args, "-T", *opts.Title)
	}
	if opts.Jobs != nil {
		args = append(args, "-j", strconv.Itoa(*opts.Jobs))
	}
	if opts.JobsParallel != nil {
		args = append(args, "-J", strconv.Itoa(*opts.JobsParallel))
	}
	if opts.Verbose != nil && *opts.Verbose {
		args = append(args, "-v")
	}
	if opts.Quiet != nil && *opts.Quiet {
		args = append(args, "-q")
	}

	// Filtering options
	if opts.Dbname != nil {
		args = append(args, "-d", *opts.Dbname)
	}
	if opts.Dbuser != nil {
		args = append(args, "-u", *opts.Dbuser)
	}
	if opts.Appname != nil {
		args = append(args, "-N", *opts.Appname)
	}
	if opts.ClientHost != nil {
		args = append(args, "-c", *opts.ClientHost)
	}

	// Time options
	if opts.Begin != nil {
		args = append(args, "-b", *opts.Begin)
	}
	if opts.End != nil {
		args = append(args, "-e", *opts.End)
	}

	// Report options
	if opts.Top != nil {
		args = append(args, "-t", strconv.Itoa(*opts.Top))
	}
	if opts.Sample != nil {
		args = append(args, "-s", strconv.Itoa(*opts.Sample))
	}
	if opts.Maxlength != nil {
		args = append(args, "-m", strconv.Itoa(*opts.Maxlength))
	}

	// Graph options
	if opts.Average != nil {
		args = append(args, "-a", strconv.Itoa(*opts.Average))
	}
	if opts.HistoAverage != nil {
		args = append(args, "-A", strconv.Itoa(*opts.HistoAverage))
	}
	if opts.Nograph != nil && *opts.Nograph {
		args = append(args, "-G")
	}

	// Output options
	if opts.Extension != nil {
		args = append(args, "-x", *opts.Extension)
	}
	if opts.Prettify != nil && !*opts.Prettify {
		args = append(args, "-P")
	}
	if opts.QueryNumbering != nil && *opts.QueryNumbering {
		args = append(args, "-Q")
	}

	// Special modes
	if opts.SelectOnly != nil && *opts.SelectOnly {
		args = append(args, "-S")
	}
	if opts.WatchMode != nil && *opts.WatchMode {
		args = append(args, "-w")
	}
	if opts.Incremental != nil && *opts.Incremental {
		args = append(args, "-I")
	}
	if opts.Explode != nil && *opts.Explode {
		args = append(args, "-E")
	}

	// Exclude options
	for _, user := range opts.ExcludeUser {
		args = append(args, "-U", user)
	}
	for _, app := range opts.ExcludeAppname {
		args = append(args, "--exclude-appname", app)
	}
	for _, client := range opts.ExcludeClient {
		args = append(args, "--exclude-client", client)
	}
	for _, db := range opts.ExcludeDb {
		args = append(args, "--exclude-db", db)
	}

	// Include options
	for _, query := range opts.IncludeQuery {
		args = append(args, "--include-query", query)
	}
	for _, pid := range opts.IncludePid {
		args = append(args, "--include-pid", pid)
	}
	for _, session := range opts.IncludeSession {
		args = append(args, "--include-session", session)
	}

	// Advanced options
	if opts.Prefix != nil {
		args = append(args, "-p", *opts.Prefix)
	}
	if opts.Ident != nil {
		args = append(args, "-i", *opts.Ident)
	}
	if opts.Timezone != nil {
		args = append(args, "-Z", *opts.Timezone)
	}
	if opts.LogTimezone != nil {
		args = append(args, "--log-timezone", *opts.LogTimezone)
	}

	// Add the log file path
	args = append(args, logPath)

	return args
}

// SetDefaults sets default values for options that should have defaults
func (opts *PgbadgerOptions) SetDefaults() {
	if opts.DataDir == nil {
		defaultDir := "/var/pv/data/log"
		opts.DataDir = &defaultDir
	}
}
