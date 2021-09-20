package obfuscator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTrackerHappyPath(t *testing.T) {
	tracker := NewSimpleTracker()
	assert.Equal(t, map[string]string{}, tracker.Report().AsMap())
	tracker.AddReplacement("A", "a", "b")
	assert.Equal(t, map[string]string{"a": "b"}, tracker.Report().AsMap())
}

func TestSimpleTrackerGetReplacement(t *testing.T) {
	tracker := NewSimpleTracker()
	tracker.AddReplacement("a", "a", "b")
	assert.Equal(t, tracker.GenerateIfAbsent("a", nil), "b")
	assert.Equal(t, tracker.GenerateIfAbsent("c", nil), "")
	assert.Equal(t, tracker.GenerateIfAbsent("D", func() string { return strings.ToLower("D") }), "d")
	tracker.AddReplacement("d", "D", "d")
	assert.Equal(t, tracker.GenerateIfAbsent("F", func() string { return strings.ToLower("F") }), "f")
	tracker.AddReplacement("f", "F", "f")
	assert.Equal(t, map[string]string{"D": "d", "a": "b", "F": "f"}, tracker.Report().AsMap())
}

func TestReportLeakingBack(t *testing.T) {
	tracker := NewSimpleTracker()
	tracker.AddReplacement("foo", "foo", "bar")
	mapping := tracker.Report()
	mapping.Replacements = append(mapping.Replacements, Replacement{Canonical: "foo", ReplacedWith: "baz", Occurrences: []Occurrence{{Original: "foo", Count: 1}}})

	assert.Equal(t, "bar", tracker.GenerateIfAbsent("foo", nil))
}

func TestSimpleReporterInitialize(t *testing.T) {
	tracker := NewSimpleTracker()
	tracker.Initialize(map[string]string{"a": "b"})
	assert.Equal(t, "b", tracker.GenerateIfAbsent("a", nil))
	assert.Equal(t, "b", tracker.GenerateIfAbsent("a", func() string { return strings.ToUpper("a") }))
	assert.Equal(t, "", tracker.GenerateIfAbsent("c", nil))
	assert.Equal(t, "C", tracker.GenerateIfAbsent("c", func() string { return strings.ToUpper("c") }))
}
