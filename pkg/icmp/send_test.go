package icmp

import (
	"strings"
	"testing"
)

func TestChunks_Empty(t *testing.T) {
	result := Chunks("", 10)
	if result != nil {
		t.Fatalf("expected nil for empty input, got %v", result)
	}
}

func TestChunks_ShorterThanChunkSize(t *testing.T) {
	result := Chunks("hello", 100)
	if len(result) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(result))
	}
	if result[0] != "hello" {
		t.Fatalf("expected %q, got %q", "hello", result[0])
	}
}

func TestChunks_ExactMultiple(t *testing.T) {
	result := Chunks("abcdef", 2)
	expected := []string{"ab", "cd", "ef"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d chunks, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Fatalf("chunk %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestChunks_NonMultiple(t *testing.T) {
	result := Chunks("abcde", 2)
	expected := []string{"ab", "cd", "e"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d chunks, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Fatalf("chunk %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestChunks_ChunkSizeOne(t *testing.T) {
	result := Chunks("abc", 1)
	expected := []string{"a", "b", "c"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d chunks, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Fatalf("chunk %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestChunks_BinaryData(t *testing.T) {
	// 'é' encodes as two bytes: 0xc3 0xa9.
	// With chunkSize=1 the result must be 4 single-byte strings, not 3 runes.
	input := "a\xc3\xa9b" // bytes: 'a', 0xc3, 0xa9, 'b'
	result := Chunks(input, 1)
	if len(result) != 4 {
		t.Fatalf("expected 4 byte-sized chunks, got %d: %v", len(result), result)
	}
	if strings.Join(result, "") != input {
		t.Fatalf("round-trip mismatch: expected %q, got %q", input, strings.Join(result, ""))
	}
}
