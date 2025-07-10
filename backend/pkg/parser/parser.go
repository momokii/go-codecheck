package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseSemgrepReport parses a Semgrep JSON report file into a SemgrepReport struct
func ParseSemgrepReport(filePath string) (*SemgrepReport, error) {
	// Read the semgrep report file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON into our struct
	var report SemgrepReport
	err = json.Unmarshal(data, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func ParseResultSemgrepFromDatabase(jsonString string) (*SemgrepReport, error) {
	// Parse the JSON string directly into our struct
	var report SemgrepReport
	if jsonString == "" {
		return nil, fmt.Errorf("empty JSON string provided")
	}

	err := json.Unmarshal([]byte(jsonString), &report)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Semgrep report: %w", err)
	}

	return &report, nil
}

// GetVulnerabilities returns a filtered list of security vulnerabilities from the report
func (report *SemgrepReport) GetVulnerabilities() []SemgrepResult {
	// Only include results that have severity of WARNING or ERROR
	var vulnerabilities []SemgrepResult
	for _, result := range report.Results {
		severity := result.Extra.Severity
		if severity == "WARNING" || severity == "ERROR" {
			vulnerabilities = append(vulnerabilities, result)
		}
	}
	return vulnerabilities
}

// GetErrorCount returns the total number of errors in the report
func (report *SemgrepReport) GetErrorCount() int {
	return len(report.Errors)
}

// GetResultsSummary returns a summary of findings by severity
func (report *SemgrepReport) GetResultsSummary() map[string]int {
	summary := make(map[string]int)

	for _, result := range report.Results {
		severity := result.Extra.Severity
		summary[severity]++
	}

	return summary
}

// SaveReport saves a modified report back to JSON
func (report *SemgrepReport) SaveReport(filePath string) error {
	data, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func ParseSemgrepResult(result string) (*ScanStatus, *ScanSummary, error) {
	if result == "" {
		return nil, nil, fmt.Errorf("empty semgrep result")
	}

	scanner := bufio.NewScanner(strings.NewReader(result))
	var (
		inStatusTable bool
		inSummary     bool
		inScanSkipped bool
		status        = ScanStatus{
			LanguageStats: make([]LanguageStat, 0),
		}
		summary = ScanSummary{
			ScanSkipped: make([]string, 0),
		}
		// regex for "Scanning X file ... with Y Code rules"
		reScanning = regexp.MustCompile(`Scanning\s+(\d+)\s+file.*with\s+(\d+)\s+Code rules`)
		// bullet‐point summary (Findings, Rules run, dll)
		reFindings    = regexp.MustCompile(`Findings:\s+(\d+)(?:\s+\((\d+)\s+blocking\))?`)
		reRulesRun    = regexp.MustCompile(`Rules run:\s+(\d+)`)
		reTargets     = regexp.MustCompile(`Targets scanned:\s+(\d+)`)
		reParsedLines = regexp.MustCompile(`Parsed lines:\s+(.+)%`)
		reIgnoreInfo  = regexp.MustCompile(`•\s+(No ignore information available|.+ ignore information available)`)
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 1) detection header Scan Status
		if strings.Contains(line, "┌─────────────┐") {
			inStatusTable = true
			continue
		}
		// 2) detection header Scan Summary
		if strings.Contains(line, "┌──────────────┐") {
			inStatusTable = false
			inSummary = true
			continue
		}
		// --- Parsing Scan Status ---
		if inStatusTable {
			if m := reScanning.FindStringSubmatch(line); m != nil {
				scannedFiles, _ := strconv.Atoi(m[1])
				totalRules, _ := strconv.Atoi(m[2])
				status.ScannedFiles = scannedFiles
				status.TotalRulesToRun = totalRules
				continue
			}

			// Match both formats of language lines:
			// 1. With origin and total rules: "<multilang>      61      57          Community    1062"
			// 2. Just language stats: "go               83      28"
			langLineWithOrigin := regexp.MustCompile(`^(\S+)\s+(\d+)\s+(\d+)\s+(\S+)\s+(\d+)$`)
			langLineSimple := regexp.MustCompile(`^(\S+)\s+(\d+)\s+(\d+)$`)

			if m := langLineWithOrigin.FindStringSubmatch(line); m != nil {
				rules, _ := strconv.Atoi(m[2])
				files, _ := strconv.Atoi(m[3])
				totalR, _ := strconv.Atoi(m[5])
				status.LanguageStats = append(status.LanguageStats, LanguageStat{
					Language:   m[1],
					Rules:      rules,
					Files:      files,
					Origin:     m[4],
					TotalRules: totalR,
				})
			} else if m := langLineSimple.FindStringSubmatch(line); m != nil {
				// For lines without origin and total rules, like "go 83 28"
				rules, _ := strconv.Atoi(m[2])
				files, _ := strconv.Atoi(m[3])
				status.LanguageStats = append(status.LanguageStats, LanguageStat{
					Language: m[1],
					Rules:    rules,
					Files:    files,
					// These will be empty or zero since they're not in this line
					Origin:     "",
					TotalRules: 0,
				})
			}
		}
		// --- Parsing Scan Summary ---
		if inSummary {
			switch {
			case strings.HasPrefix(line, "✅"):
				// if first line of summary contains checkmark → success
				summary.Success = true
			case reFindings.MatchString(line):
				m := reFindings.FindStringSubmatch(line)
				findings, err := strconv.Atoi(m[1])
				if err == nil {
					summary.Findings = findings
				}
				if len(m) > 2 && m[2] != "" {
					blockingFindings, err := strconv.Atoi(m[2])
					if err == nil {
						summary.BlockingFindings = blockingFindings
					}
				}
			case reRulesRun.MatchString(line):
				m := reRulesRun.FindStringSubmatch(line)
				rulesRun, err := strconv.Atoi(m[1])
				if err == nil {
					summary.RulesRun = rulesRun
				}
			case reTargets.MatchString(line):
				m := reTargets.FindStringSubmatch(line)
				targetsScanned, err := strconv.Atoi(m[1])
				if err == nil {
					summary.TargetsScanned = targetsScanned
				}
			case reParsedLines.MatchString(line):
				m := reParsedLines.FindStringSubmatch(line)
				if len(m) > 1 {
					summary.ParsedLinesPercent = m[1] + "%"
				}
			case reIgnoreInfo.MatchString(line):
				m := reIgnoreInfo.FindStringSubmatch(line)
				if len(m) > 1 {
					summary.IgnoreInfo = m[1]
				}
			case strings.Contains(line, "Scan skipped:"):
				// Toggle the scan skipped section flag
				inScanSkipped = true
				continue
			case inScanSkipped && strings.Contains(line, "◦"):
				// This is a skipped item line (with proper indentation)
				// Extract everything after the bullet point
				parts := strings.SplitN(line, "◦", 2)
				if len(parts) > 1 {
					skippedInfo := strings.TrimSpace(parts[1])
					if skippedInfo != "" {
						summary.ScanSkipped = append(summary.ScanSkipped, skippedInfo)
					}
				}
				continue
			case (strings.Contains(line, "•") || strings.Contains(line, "For a detailed list")) && inScanSkipped:
				// New bullet point or end of skipped section, turn off the flag
				inScanSkipped = false
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scan stderr parse error: %w", err)
	}

	// Validation
	if len(status.LanguageStats) == 0 && status.ScannedFiles == 0 {
		return nil, nil, fmt.Errorf("no language statistics found in the scan result")
	}

	return &status, &summary, nil
}
