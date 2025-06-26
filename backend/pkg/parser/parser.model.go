package parser

// ============================= SEMGREP REPORT STRUCTURES

// SemgrepReport represents the full structure of a Semgrep report
type SemgrepReport struct {
	Version string          `json:"version"`
	Results []SemgrepResult `json:"results"`
	Errors  []SemgrepError  `json:"errors"`
	Paths   SemgrepPaths    `json:"paths"`
	Time    SemgrepTime     `json:"time"`
}

// SemgrepResult represents an individual finding in the report
type SemgrepResult struct {
	CheckID string       `json:"check_id"`
	Path    string       `json:"path"`
	Start   Position     `json:"start"`
	End     Position     `json:"end"`
	Extra   ResultExtras `json:"extra"`
}

// Position represents a location in code
type Position struct {
	Line   int `json:"line"`
	Col    int `json:"col"`
	Offset int `json:"offset"`
}

// ResultExtras contains additional information about a result
type ResultExtras struct {
	Message         string                 `json:"message"`
	Metadata        map[string]interface{} `json:"metadata"`
	Severity        string                 `json:"severity"`
	Fingerprint     string                 `json:"fingerprint"`
	Lines           string                 `json:"lines"`
	ValidationState string                 `json:"validation_state"`
	EngineKind      string                 `json:"engine_kind"`
}

// SemgrepError represents an error that occurred during scanning
type SemgrepError struct {
	Code    int           `json:"code"`
	Level   string        `json:"level"`
	Type    []interface{} `json:"type"`
	Message string        `json:"message"`
	Path    string        `json:"path"`
	Spans   []Span        `json:"spans"`
}

// Span represents a code span referenced in an error
type Span struct {
	File  string   `json:"file"`
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// SemgrepPaths contains information about scanned paths
type SemgrepPaths struct {
	Scanned []string `json:"scanned"`
}

// SemgrepTime contains timing information
type SemgrepTime struct {
	Rules          []interface{}  `json:"rules"`
	RulesParseTime float64        `json:"rules_parse_time"`
	ProfilingTimes ProfilingTimes `json:"profiling_times"`
	ParsingTime    ParsingTime    `json:"parsing_time"`
	Targets        []interface{}  `json:"targets"`
	TotalBytes     int            `json:"total_bytes"`
	MaxMemoryBytes int64          `json:"max_memory_bytes"`
}

// ProfilingTimes contains timing metrics for the scan
type ProfilingTimes struct {
	ConfigTime  float64 `json:"config_time"`
	CoreTime    float64 `json:"core_time"`
	IgnoresTime float64 `json:"ignores_time"`
	TotalTime   float64 `json:"total_time"`
}

// ParsingTime contains parsing timing information
type ParsingTime struct {
	TotalTime     float64       `json:"total_time"`
	PerFileTime   PerFileTime   `json:"per_file_time"`
	VerySlowFiles []interface{} `json:"very_slow_files"`
}

// PerFileTime contains per-file parsing metrics
type PerFileTime struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
}

// ============================== SEMGREP SCAN RESULT CLI
type ScanStatus struct {
	ScannedFiles    int            `json:"scanned_files"`
	TotalRulesToRun int            `json:"total_rules_to_run"`
	LanguageStats   []LanguageStat `json:"language_stats"`
}

type LanguageStat struct {
	Language   string `json:"language"`
	Rules      int    `json:"rules"`
	Files      int    `json:"files"`
	Origin     string `json:"origin"`
	TotalRules int    `json:"total_rules"`
}

type ScanSummary struct {
	Success            bool     `json:"success"`
	Findings           int      `json:"findings"`
	BlockingFindings   int      `json:"blocking_findings"`
	RulesRun           int      `json:"rules_run"`
	TargetsScanned     int      `json:"targets_scanned"`
	ParsedLinesPercent string   `json:"parsed_lines_percent"`
	IgnoreInfo         string   `json:"ignore_info"`
	ScanSkipped        []string `json:"scan_skipped"`
}
