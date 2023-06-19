package helpers

import (
	"testing"

	"github.com/prplx/wordy/internal/types"
)

type testWithText struct {
	text string
}

func (t testWithText) GetText() string {
	return t.text
}

func TestBuildMessageFromSliceOfTexted(t *testing.T) {
	slice := []types.WithText{
		testWithText{"test1"},
		testWithText{"test2"},
		testWithText{"test3"},
	}
	result := BuildMessageFromSliceOfTexted(slice)
	expected := "- test1\n- test2\n- test3"
	if result != expected {
		t.Fatalf("Expected %s, got %s", expected, result)
	}

	slice = []types.WithText{
		testWithText{"test1"},
	}
	result = BuildMessageFromSliceOfTexted(slice)
	expected = "- test1"
	if result != expected {
		t.Fatalf("Expected %s, got %s", expected, result)
	}

	slice = []types.WithText{}
	result = BuildMessageFromSliceOfTexted(slice)
	expected = ""
	if result != expected {
		t.Fatalf("Expected %s, got %s", expected, result)
	}
}
