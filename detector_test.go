package main

import (
	"testing"
)

func TestExtractSymbols(t *testing.T) {
	chunk := `+def calculate_price(amount, tax_rate, discount):
-def calculate_price(amount, tax_rate):
+func HandleRequest(ctx context.Context, req *Request) error {
+func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
+  parser.add_argument("--output-format", help="fmt")
+export function fetchUsers(page, limit) {`

	syms := extractSymbols(chunk)
	want := map[string]bool{
		"calculate_price": false,
		"HandleRequest":   false,
		"ServeHTTP":       false,
		"--output-format": false,
		"fetchUsers":      false,
	}
	for _, s := range syms {
		if _, ok := want[s]; ok {
			want[s] = true
		}
	}
	for sym, found := range want {
		if !found {
			t.Errorf("missing expected symbol: %s (got: %v)", sym, syms)
		}
	}
}

func TestDetectDrift(t *testing.T) {
	diff := "diff --git a/api.py b/api.py\n--- a/api.py\n+++ b/api.py\n@@ -10,3 +10,3 @@\n-def get_users(limit):\n+def get_users(limit, offset, fields):\n"
	docs := map[string]string{
		"README.md":    "## API\nCall `get_users` to retrieve user records.",
		"CHANGELOG.md": "## v1.0\nInitial release with billing module.",
	}
	r := Detect(diff, docs, 80.0)
	if len(r.Drifts) != 1 {
		t.Fatalf("expected 1 drift point, got %d: %+v", len(r.Drifts), r.Drifts)
	}
	if r.Drifts[0].Symbol != "get_users" {
		t.Errorf("expected symbol get_users, got %s", r.Drifts[0].Symbol)
	}
	if r.Drifts[0].StaleDocs[0] != "README.md" {
		t.Errorf("expected stale doc README.md, got %v", r.Drifts[0].StaleDocs)
	}
	if r.Score != 0 {
		t.Errorf("expected score 0, got %.1f", r.Score)
	}
	if r.Pass {
		t.Error("expected FAIL but got PASS")
	}
}

func TestNoDriftWhenDocsUpdated(t *testing.T) {
	diff := "diff --git a/api.py b/api.py\n--- a/api.py\n+++ b/api.py\n@@ -10,1 +10,1 @@\n-def get_users(limit):\n+def get_users(limit, offset):\ndiff --git a/README.md b/README.md\n--- a/README.md\n+++ b/README.md\n@@ -5,1 +5,1 @@\n-get_users(limit)\n+get_users(limit, offset)\n"
	docs := map[string]string{
		"README.md": "Call get_users(limit, offset) to fetch users.",
	}
	r := Detect(diff, docs, 80.0)
	if len(r.Drifts) != 0 {
		t.Errorf("expected 0 drifts, got %d: %+v", len(r.Drifts), r.Drifts)
	}
	if r.Score != 100 {
		t.Errorf("expected score 100, got %.1f", r.Score)
	}
	if !r.Pass {
		t.Error("expected PASS but got FAIL")
	}
}

func TestParseDiffSeparatesCodeAndDocs(t *testing.T) {
	diff := "diff --git a/main.go b/main.go\n--- a/main.go\n+++ b/main.go\n@@ -1 +1 @@\n-old\n+new\ndiff --git a/docs/guide.md b/docs/guide.md\n--- a/docs/guide.md\n+++ b/docs/guide.md\n@@ -1 +1 @@\n-old doc\n+new doc\n"
	code, docs := parseDiff(diff)
	if len(code) != 1 {
		t.Errorf("expected 1 code file, got %d", len(code))
	}
	if _, ok := code["main.go"]; !ok {
		t.Error("expected main.go in code files")
	}
	if !docs["docs/guide.md"] {
		t.Error("expected docs/guide.md in changed docs")
	}
}

func TestNoChangesFullScore(t *testing.T) {
	r := Detect("", map[string]string{"README.md": "hello"}, 80.0)
	if r.Score != 100 {
		t.Errorf("expected 100 score for empty diff, got %.1f", r.Score)
	}
	if !r.Pass {
		t.Error("expected PASS for empty diff")
	}
}
