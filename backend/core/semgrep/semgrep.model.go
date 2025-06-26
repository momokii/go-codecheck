package semgrep

type ScanResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}
