package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

var funcRe = regexp.MustCompile(`(?m)^[+-].*?(?:def |func |function |class )\s*(?:\([^)]*\)\s*)?([A-Za-z_]\w{1,})`)
var flagRe = regexp.MustCompile(`"(--[\w][\w-]*)"`) 

// Drift represents a single code-doc synchronization gap.
type Drift struct {
	Symbol    string   `json:"symbol"`
	CodeFile  string   `json:"code_file"`
	StaleDocs []string `json:"stale_docs"`
}

// Report is the full analysis output.
type Report struct {
	Drifts    []Drift `json:"drifts"`
	Score     float64 `json:"freshness_score"`
	Pass      bool    `json:"pass"`
	Threshold float64 `json:"threshold"`
}

// Detect runs the full drift analysis on a unified diff and doc file contents.
func Detect(diff string, docContents map[string]string, threshold float64) Report {
	codeChunks, changedDocs := parseDiff(diff)
	type symEntry struct {
		name, file string
	}
	seen := map[string]bool{}
	var syms []symEntry
	for f, chunk := range codeChunks {
		for _, s := range extractSymbols(chunk) {
			key := s + "::" + f
			if !seen[key] {
				seen[key] = true
				syms = append(syms, symEntry{s, f})
			}
		}
	}
	var drifts []Drift
	totalRefs, staleRefs := 0, 0
	for _, si := range syms {
		var stale []string
		lower := strings.ToLower(si.name)
		for docPath, content := range docContents {
			if strings.Contains(strings.ToLower(content), lower) {
				totalRefs++
				if !changedDocs[docPath] {
					stale = append(stale, docPath)
					staleRefs++
				}
			}
		}
		if len(stale) > 0 {
			drifts = append(drifts, Drift{Symbol: si.name, CodeFile: si.file, StaleDocs: stale})
		}
	}
	score := 100.0
	if totalRefs > 0 {
		score = float64(totalRefs-staleRefs) / float64(totalRefs) * 100
	}
	return Report{Drifts: drifts, Score: score, Pass: score >= threshold, Threshold: threshold}
}

func parseDiff(diff string) (code map[string]string, docs map[string]bool) {
	code = make(map[string]string)
	docs = make(map[string]bool)
	for _, part := range strings.Split(diff, "diff --git ") {
		var fname string
		for _, line := range strings.Split(part, "\n") {
			if strings.HasPrefix(line, "+++ b/") {
				fname = line[6:]
				break
			}
		}
		if fname == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(fname))
		if ext == ".md" || ext == ".rst" || ext == ".adoc" {
			docs[fname] = true
		} else {
			code[fname] = part
		}
	}
	return
}

func extractSymbols(chunk string) []string {
	seen := map[string]bool{}
	var out []string
	for _, re := range []*regexp.Regexp{funcRe, flagRe} {
		for _, m := range re.FindAllStringSubmatch(chunk, -1) {
			s := m[1]
			if !seen[s] {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return out
}
