package fs

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Ensure path traversal is rejected for write action.
func TestWriteFileTraversalDenied(t *testing.T) {
	root := t.TempDir()
	act := NewWriteFile(Config{Roots: []string{root}, AllowWrite: true})

	target := filepath.Join(root, "ok.txt")
	args, _ := json.Marshal(map[string]string{"path": target, "content": "hi"})
	if _, err := act.Invoke(context.Background(), args); err != nil {
		t.Fatalf("expected write in root to succeed: %v", err)
	}

	outside := filepath.Join(root, "..", "evil.txt")
	args2, _ := json.Marshal(map[string]string{"path": outside, "content": "nope"})
	if _, err := act.Invoke(context.Background(), args2); err == nil {
		t.Fatalf("expected traversal write to fail")
	}
	if _, err := os.Stat(filepath.Clean(outside)); err == nil {
		t.Fatalf("traversal file unexpectedly created")
	}
}
