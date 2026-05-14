package hasher

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/idleberg/go-hashman/internal/algo"
)

func TestHashFileKnownValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		id   string
		want string
	}{
		{"md5", "5d41402abc4b2a76b9719d911017c592"},
		{"sha1", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"sha256", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"crc32", "3610a686"},
	}

	for _, tt := range tests {
		var alg algo.Algorithm
		for _, a := range algo.Registry {
			if a.ID == tt.id {
				alg = a
				break
			}
		}

		results := HashFile(path, []algo.Algorithm{alg}, 1)
		if len(results) != 1 {
			t.Fatalf("%s: expected 1 result, got %d", tt.id, len(results))
		}
		if results[0].Err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.id, results[0].Err)
		}
		if results[0].Hash != tt.want {
			t.Errorf("%s: got %s, want %s", tt.id, results[0].Hash, tt.want)
		}
	}
}

func TestHashFileNotFound(t *testing.T) {
	alg := algo.Registry[0]
	results := HashFile("/nonexistent/file", []algo.Algorithm{alg}, 1)
	if results[0].Err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestHashFileNonRegular(t *testing.T) {
	dir := t.TempDir()
	link := filepath.Join(dir, "link")
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skip("symlinks not supported")
	}

	alg := algo.Registry[0]
	results := HashFile(link, []algo.Algorithm{alg}, 1)
	if results[0].Err == nil {
		t.Error("expected error for symlink")
	}
}

func TestHashFileConcurrency(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	results := HashFile(path, algo.Registry, 4)
	if len(results) != len(algo.Registry) {
		t.Fatalf("expected %d results, got %d", len(algo.Registry), len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("%s: unexpected error: %v", r.Algorithm.ID, r.Err)
		}
		if r.Hash == "" {
			t.Errorf("%s: empty hash", r.Algorithm.ID)
		}
	}
}
