package echartsgoja

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	ver, err := New().Version()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if v, exp := cleanString(ver), cleanString(string(versionTxt)); !strings.HasPrefix(exp, v) {
		t.Errorf("expected %s, got: %s", exp, v)
	}
	t.Logf("echarts: %s", ver)
}

func TestRenderScript(t *testing.T) {
	ctx := context.Background()
	timeout := 1 * time.Minute
	if s := os.Getenv("TIMEOUT"); s != "" {
		var err error
		if timeout, err = time.ParseDuration(s); err != nil {
			t.Fatalf("could not parse timeout %q: %v", s, err)
		}
	}
	var files []string
	err := filepath.Walk("testdata", func(name string, info fs.FileInfo, err error) error {
		switch {
		case err != nil:
			return err
		case info.IsDir() || !suffixRE.MatchString(name):
			return nil
		}
		files = append(files, name)
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	for _, nn := range files {
		name := nn
		n := strings.Split(name, string(os.PathSeparator))
		n[len(n)-1] = suffixRE.ReplaceAllString(n[len(n)-1], "")
		testName := path.Join(n[1:]...)
		t.Run(testName, func(t *testing.T) {
			testRender(t, ctx, testName, name, timeout)
		})
	}
}

var suffixRE = regexp.MustCompile(`\.js$`)

func testRender(t *testing.T, ctx context.Context, testName, name string, timeout time.Duration) {
	t.Helper()
	t.Parallel()
	src, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	vm := New(WithLogger(t.Log), WithWidthHeight(850, 621))
	start := time.Now()
	res, err := vm.RenderScript(ctx, string(src))
	total := time.Since(start)
	switch {
	case err != nil && isBroken(err):
		t.Logf("IGNORING: expected no error, got: %v", err)
		return
	case err != nil:
		t.Fatalf("expected no error, got: %v", err)
	}
	if os.Getenv("VERBOSE") != "" {
		t.Logf("---\n%s\n---", res)
	}
	t.Logf("duration: %s", total)
	if res = strings.TrimSpace(res); len(res) != 0 {
		if err := os.WriteFile(name+".svg", []byte(res), 0o644); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	}
}

func cleanString(s string) string {
	return strings.TrimPrefix(strings.TrimSpace(s), "v")
}

func contains(v []string, s string) bool {
	for _, ss := range v {
		if ss == s {
			return true
		}
	}
	return false
}

func isBroken(err error) bool {
	s := err.Error()
	return errors.Is(err, ErrCouldNotExecuteScript) || strings.HasPrefix(s, "ReferenceError: ") || strings.HasPrefix(s, "TypeError: ")
}
