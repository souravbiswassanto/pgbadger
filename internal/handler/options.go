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
	Format       *string  `json:"format,omitempty"`        // -f, --format
	Outfiles     []string `json:"outfiles,omitempty"`      // -o, --outfile (multiple)
	Outdir       *string  `json:"outdir,omitempty"`        // -O, --outdir
	Title        *string  `json:"title,omitempty"`         // -T, --title
	Jobs         *int     `json:"jobs,omitempty"`          // -j, --jobs
	JobsParallel *int     `json:"jobs_parallel,omitempty"` // -J, --Jobs
	Verbose      *bool    `json:"verbose,omitempty"`       // -v, --verbose
	Quiet        *bool    `json:"quiet,omitempty"`         // -q, --quiet

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
	NoPrettify     *bool   `json:"no_prettify,omitempty"`     // -P, --no-prettify (presence disables prettify)
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
	ExcludeFile    []string `json:"exclude_file,omitempty"`    // --exclude-file
	ExcludeLine    []string `json:"exclude_line,omitempty"`    // --exclude-line
	ExcludeQuery   []string `json:"exclude_query,omitempty"`   // --exclude-query

	// Include options
	IncludeQuery   []string `json:"include_query,omitempty"`   // --include-query
	IncludePid     []string `json:"include_pid,omitempty"`     // --include-pid
	IncludeSession []string `json:"include_session,omitempty"` // --include-session

	// Advanced options
	Prefix      *string `json:"prefix,omitempty"`       // -p, --prefix
	Ident       *string `json:"ident,omitempty"`        // -i, --ident
	Timezone    *string `json:"timezone,omitempty"`     // -Z, --timezone
	LogTimezone *string `json:"log_timezone,omitempty"` // --log-timezone

	// Additional options implemented
	LogfileList *string `json:"logfile_list,omitempty"` // -L, --logfile-list
	LastParsed  *string `json:"last_parsed,omitempty"`  // -l, --last-parsed
	NoComment   *bool   `json:"no_comment,omitempty"`   // -C, --nocomment
	DNSResolv   *bool   `json:"dns_resolv,omitempty"`   // -D, --dns-resolv
	HTMLOutdir  *string `json:"html_outdir,omitempty"`  // -H, --html-outdir
	NoMultiline *bool   `json:"no_multiline,omitempty"` // -M, --no-multiline

	RemoteHost  *string  `json:"remote_host,omitempty"`  // -r, --remote-host
	SSHIdentity *string  `json:"ssh_identity,omitempty"` // --ssh-identity
	SSHOption   []string `json:"ssh_option,omitempty"`   // --ssh-option
	SSHPort     *int     `json:"ssh_port,omitempty"`     // --ssh-port
	SSHProgram  *string  `json:"ssh_program,omitempty"`  // --ssh-program
	SSHTimeout  *int     `json:"ssh_timeout,omitempty"`  // --ssh-timeout
	SSHUser     *string  `json:"ssh_user,omitempty"`     // --ssh-user

	Retention    *int    `json:"retention,omitempty"`     // -R, --retention
	ExtraFiles   *bool   `json:"extra_files,omitempty"`   // -X, --extra-files
	Zcat         *string `json:"zcat,omitempty"`          // -z, --zcat
	Command      *string `json:"command,omitempty"`       // --command
	CSVSeparator *string `json:"csv_separator,omitempty"` // --csv-separator

	DayReport   *string `json:"day_report,omitempty"`   // --day-report
	MonthReport *string `json:"month_report,omitempty"` // --month-report

	// Disable specific report sections
	DisableAutovacuum *bool `json:"disable_autovacuum,omitempty"`
	DisableCheckpoint *bool `json:"disable_checkpoint,omitempty"`
	DisableConnection *bool `json:"disable_connection,omitempty"`
	DisableError      *bool `json:"disable_error,omitempty"`
	DisableHourly     *bool `json:"disable_hourly,omitempty"`
	DisableLock       *bool `json:"disable_lock,omitempty"`
	DisableQuery      *bool `json:"disable_query,omitempty"`
	DisableSession    *bool `json:"disable_session,omitempty"`
	DisableTemporary  *bool `json:"disable_temporary,omitempty"`
	DisableType       *bool `json:"disable_type,omitempty"`

	DumpAllQueries *bool `json:"dump_all_queries,omitempty"` // --dump-all-queries
	DumpRawCSV     *bool `json:"dump_raw_csv,omitempty"`     // --dump-raw-csv
	EnableChecksum *bool `json:"enable_checksum,omitempty"`  // --enable-checksum

	IncludeFile []string `json:"include_file,omitempty"` // --include-file
	IncludeTime []string `json:"include_time,omitempty"` // --include-time

	IsoWeekNumber *bool `json:"iso_week_number,omitempty"` // --iso-week-number
	KeepComments  *bool `json:"keep_comments,omitempty"`   // --keep-comments
	StartMonday   *bool `json:"start_monday,omitempty"`    // --start-monday

	Tempdir *string `json:"tempdir,omitempty"`  // --tempdir
	PIDDir  *string `json:"pid_dir,omitempty"`  // --pid-dir
	PIDFile *string `json:"pid_file,omitempty"` // --pid-file

	NoFork         *bool `json:"no_fork,omitempty"`         // --no-fork
	NoProcessInfo  *bool `json:"no_process_info,omitempty"` // --no-process-info
	NoProgressbar  *bool `json:"no_progressbar,omitempty"`  // --no-progressbar
	NoReport       *bool `json:"no_report,omitempty"`       // --noreport
	NoWeek         *bool `json:"no_week,omitempty"`         // --no-week
	NormalizedOnly *bool `json:"normalized_only,omitempty"` // --normalized-only
	PgbouncerOnly  *bool `json:"pgbouncer_only,omitempty"`  // --pgbouncer-only

	PieLimit     *int  `json:"pie_limit,omitempty"`     // --pie-limit
	PrettifyJSON *bool `json:"prettify_json,omitempty"` // --prettify-json
	Rebuild      *bool `json:"rebuild,omitempty"`       // --rebuild

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
	// support multiple -o outfiles
	for _, of := range opts.Outfiles {
		args = append(args, "-o", of)
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
	// multiple outfiles support
	for _, of := range opts.Outfiles {
		args = append(args, "-o", of)
	}
	if opts.Outdir != nil {
		args = append(args, "-O", *opts.Outdir)
	}
	if opts.NoPrettify != nil && *opts.NoPrettify {
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
	for _, f := range opts.ExcludeFile {
		args = append(args, "--exclude-file", f)
	}
	for _, l := range opts.ExcludeLine {
		args = append(args, "--exclude-line", l)
	}
	for _, q := range opts.ExcludeQuery {
		args = append(args, "--exclude-query", q)
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
	for _, f := range opts.IncludeFile {
		args = append(args, "--include-file", f)
	}
	for _, t := range opts.IncludeTime {
		args = append(args, "--include-time", t)
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

	// Additional options
	if opts.LogfileList != nil {
		args = append(args, "-L", *opts.LogfileList)
	}
	if opts.LastParsed != nil {
		args = append(args, "-l", *opts.LastParsed)
	}
	if opts.NoComment != nil && *opts.NoComment {
		args = append(args, "-C")
	}
	if opts.DNSResolv != nil && *opts.DNSResolv {
		args = append(args, "-D")
	}
	if opts.HTMLOutdir != nil {
		args = append(args, "-H", *opts.HTMLOutdir)
	}
	if opts.NoMultiline != nil && *opts.NoMultiline {
		args = append(args, "-M")
	}

	if opts.RemoteHost != nil {
		args = append(args, "-r", *opts.RemoteHost)
	}
	if opts.SSHIdentity != nil {
		args = append(args, "--ssh-identity", *opts.SSHIdentity)
	}
	for _, o := range opts.SSHOption {
		args = append(args, "--ssh-option", o)
	}
	if opts.SSHPort != nil {
		args = append(args, "--ssh-port", strconv.Itoa(*opts.SSHPort))
	}
	if opts.SSHProgram != nil {
		args = append(args, "--ssh-program", *opts.SSHProgram)
	}
	if opts.SSHTimeout != nil {
		args = append(args, "--ssh-timeout", strconv.Itoa(*opts.SSHTimeout))
	}
	if opts.SSHUser != nil {
		args = append(args, "--ssh-user", *opts.SSHUser)
	}

	if opts.Retention != nil {
		args = append(args, "-R", strconv.Itoa(*opts.Retention))
	}
	if opts.ExtraFiles != nil && *opts.ExtraFiles {
		args = append(args, "-X")
	}
	if opts.Zcat != nil {
		args = append(args, "-z", *opts.Zcat)
	}
	if opts.Command != nil {
		args = append(args, "--command", *opts.Command)
	}
	if opts.CSVSeparator != nil {
		args = append(args, "--csv-separator", *opts.CSVSeparator)
	}

	if opts.DayReport != nil {
		args = append(args, "--day-report", *opts.DayReport)
	}
	if opts.MonthReport != nil {
		args = append(args, "--month-report", *opts.MonthReport)
	}

	// disable specific reports
	if opts.DisableAutovacuum != nil && *opts.DisableAutovacuum {
		args = append(args, "--disable-autovacuum")
	}
	if opts.DisableCheckpoint != nil && *opts.DisableCheckpoint {
		args = append(args, "--disable-checkpoint")
	}
	if opts.DisableConnection != nil && *opts.DisableConnection {
		args = append(args, "--disable-connection")
	}
	if opts.DisableError != nil && *opts.DisableError {
		args = append(args, "--disable-error")
	}
	if opts.DisableHourly != nil && *opts.DisableHourly {
		args = append(args, "--disable-hourly")
	}
	if opts.DisableLock != nil && *opts.DisableLock {
		args = append(args, "--disable-lock")
	}
	if opts.DisableQuery != nil && *opts.DisableQuery {
		args = append(args, "--disable-query")
	}
	if opts.DisableSession != nil && *opts.DisableSession {
		args = append(args, "--disable-session")
	}
	if opts.DisableTemporary != nil && *opts.DisableTemporary {
		args = append(args, "--disable-temporary")
	}
	if opts.DisableType != nil && *opts.DisableType {
		args = append(args, "--disable-type")
	}

	if opts.DumpAllQueries != nil && *opts.DumpAllQueries {
		args = append(args, "--dump-all-queries")
	}
	if opts.DumpRawCSV != nil && *opts.DumpRawCSV {
		args = append(args, "--dump-raw-csv")
	}
	if opts.EnableChecksum != nil && *opts.EnableChecksum {
		args = append(args, "--enable-checksum")
	}

	if opts.IsoWeekNumber != nil && *opts.IsoWeekNumber {
		args = append(args, "--iso-week-number")
	}
	if opts.KeepComments != nil && *opts.KeepComments {
		args = append(args, "--keep-comments")
	}
	if opts.StartMonday != nil && *opts.StartMonday {
		args = append(args, "--start-monday")
	}

	if opts.Tempdir != nil {
		args = append(args, "--tempdir", *opts.Tempdir)
	}
	if opts.PIDDir != nil {
		args = append(args, "--pid-dir", *opts.PIDDir)
	}
	if opts.PIDFile != nil {
		args = append(args, "--pid-file", *opts.PIDFile)
	}

	if opts.NoFork != nil && *opts.NoFork {
		args = append(args, "--no-fork")
	}
	if opts.NoProcessInfo != nil && *opts.NoProcessInfo {
		args = append(args, "--no-process-info")
	}
	if opts.NoProgressbar != nil && *opts.NoProgressbar {
		args = append(args, "--no-progressbar")
	}
	if opts.NoReport != nil && *opts.NoReport {
		args = append(args, "--noreport")
	}
	if opts.NoWeek != nil && *opts.NoWeek {
		args = append(args, "--no-week")
	}
	if opts.NormalizedOnly != nil && *opts.NormalizedOnly {
		args = append(args, "--normalized-only")
	}
	if opts.PgbouncerOnly != nil && *opts.PgbouncerOnly {
		args = append(args, "--pgbouncer-only")
	}

	if opts.PieLimit != nil {
		args = append(args, "--pie-limit", strconv.Itoa(*opts.PieLimit))
	}
	if opts.PrettifyJSON != nil && *opts.PrettifyJSON {
		args = append(args, "--prettify-json")
	}
	if opts.Rebuild != nil && *opts.Rebuild {
		args = append(args, "--rebuild")
	}

	// Add positional log file(s) unless logfile-list is used
	if opts.LogfileList == nil {
		args = append(args, logPath)
	}

	return args
}

// SetDefaults sets default values for options that should have defaults
func (opts *PgbadgerOptions) SetDefaults() {
	if opts.DataDir == nil {
		defaultDir := "/var/pv/data/log"
		opts.DataDir = &defaultDir
	}
}
