package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	fThreshold = flag.Float64("threshold", 80.0, "Min freshness score to pass (0-100)")
	fBase      = flag.String("base", "HEAD~1", "Base git ref for diff comparison")
	fFormat    = flag.String("format", "text", "Output format: text or json")
	fDocs      = flag.String("docs", ".", "Root directory to scan for doc files")
)

func main() {
	flag.Parse()
	diff, err := exec.Command("git", "diff", *fBase).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: git diff failed: %v\nEnsure you are inside a git repository.\n", err)
		os.Exit(2)
	}
	if len(diff) == 0 {
		fmt.Println("\u2705 No changes detected.")
		return
	}
	docs := loadDocs(*fDocs)
	report := Detect(string(diff), docs, *fThreshold)
	switch *fFormat {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(report)
	default:
		printReport(report)
	}
	if !report.Pass {
		os.Exit(1)
	}
}

func loadDocs(root string) map[string]string {
	docs := make(map[string]string)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" || ext == ".rst" || ext == ".adoc" {
			data, readErr := os.ReadFile(path)
			if readErr == nil {
				docs[path] = string(data)
			}
		}
		return nil
	})
	return docs
}

func printReport(r Report) {
	fmt.Printf("\n\U0001F4CA DocDrift Report \u2014 Freshness: %.0f%% (threshold: %.0f%%)\n\n", r.Score, r.Threshold)
	if len(r.Drifts) == 0 {
		fmt.Println("\u2705 No documentation drift detected!")
		return
	}
	for i, d := range r.Drifts {
		fmt.Printf("  \u26A0\uFE0F  [%d] Symbol '%s' changed in %s\n", i+1, d.Symbol, d.CodeFile)
		for _, doc := range d.StaleDocs {
			fmt.Printf("      \U0001F4C4 Stale: %s\n", doc)
		}
	}
	fmt.Println()
	if r.Pass {
		fmt.Println("\u2705 PASS \u2014 within threshold")
	} else {
		fmt.Println("\u274C FAIL \u2014 doc drift exceeds threshold. Update docs before merging.")
	}
}
