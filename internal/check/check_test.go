package check

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvChecker(t *testing.T) {
	const key = "CHECK_TEST_ENV"
	os.Unsetenv(key)
	c := EnvChecker{}
	res := c.Check(DepInput{Name: key, Type: "env"})
	if res.Status != "MISSING" {
		t.Fatalf("expected missing when unset, got %s", res.Status)
	}
	os.Setenv(key, "ok")
	defer os.Unsetenv(key)
	res = c.Check(DepInput{Name: key, Type: "env"})
	if res.Status != "OK" {
		t.Fatalf("expected OK when set, got %s", res.Status)
	}
}

func TestFileChecker(t *testing.T) {
	c := FileChecker{}
	res := c.Check(DepInput{Name: filepath.Join("this", "does", "not", "exist"), Type: "file"})
	if res.Status != "MISSING" {
		t.Fatalf("expected missing for absent file, got %s", res.Status)
	}
	res = c.Check(DepInput{Name: "check.go", Type: "file"})
	if res.Status != "OK" {
		t.Fatalf("expected OK for existing file, got %s", res.Status)
	}
}
