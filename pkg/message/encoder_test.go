package message

import (
	"testing"
)

func TestQueenSonoMarshall_Empty(t *testing.T) {
	result := QueenSonoMarshall(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestQueenSonoMarshall_Single(t *testing.T) {
	result := QueenSonoMarshall([]string{"hello"})
	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
	if result[0] != "0,hello" {
		t.Fatalf("expected %q, got %q", "0,hello", result[0])
	}
}

func TestQueenSonoMarshall_Multiple(t *testing.T) {
	result := QueenSonoMarshall([]string{"a", "b", "c"})
	expected := []string{"0,a", "1,b", "2,c"}
	for i, v := range expected {
		if result[i] != v {
			t.Fatalf("index %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestQueenSonoUnmarshall_RoundTrip(t *testing.T) {
	chunks := QueenSonoMarshall([]string{"foo", "bar", "baz"})
	for wantIdx, chunk := range chunks {
		_, idx, err := QueenSonoUnmarshall(chunk)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if idx != wantIdx {
			t.Fatalf("expected index %d, got %d", wantIdx, idx)
		}
	}
}

func TestQueenSonoUnmarshall_DataWithComma(t *testing.T) {
	encoded := "3,hello,world"
	msg, idx, err := QueenSonoUnmarshall(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 3 {
		t.Fatalf("expected index 3, got %d", idx)
	}
	if msg != "hello,world" {
		t.Fatalf("expected %q, got %q", "hello,world", msg)
	}
}

func TestQueenSonoUnmarshall_MissingComma(t *testing.T) {
	_, _, err := QueenSonoUnmarshall("nocomma")
	if err == nil {
		t.Fatal("expected error for missing comma, got nil")
	}
}

func TestQueenSonoUnmarshall_Empty(t *testing.T) {
	_, _, err := QueenSonoUnmarshall("")
	if err == nil {
		t.Fatal("expected error for empty string, got nil")
	}
}
