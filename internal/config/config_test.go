package config

import (
	"os"
	"reflect"
	"testing"
)

func TestSplitAndTrim(t *testing.T) {
	in := "a, b, ,c"
	want := []string{"a", "b", "c"}
	if got := splitAndTrim(in); !reflect.DeepEqual(got, want) {
		t.Fatalf("splitAndTrim(%q) = %v, want %v", in, got, want)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	key := "TEST_GETENV"
	_ = os.Unsetenv(key)
	if got := getEnvOrDefault(key, "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}
	_ = os.Setenv(key, "value")
	if got := getEnvOrDefault(key, "fallback"); got != "value" {
		t.Fatalf("expected value, got %q", got)
	}
}

func TestLoadValidation(t *testing.T) {

	os.Clearenv()
	_, err := Load()
	if err == nil {
		t.Fatalf("expected error when no providers configured")
	}
}
