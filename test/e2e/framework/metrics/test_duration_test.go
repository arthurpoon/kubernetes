package metrics

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTestDurationRecorder(t *testing.T) {
	tmpDir := "./"
	recorder := NewTestDurationRecorder(tmpDir)

	// Test case 1: Fast test
	fastTestName := "FastTest"
	startTime := recorder.RecordTestStart(fastTestName)
	time.Sleep(100 * time.Millisecond) // Simulate short test
	recorder.RecordTestEnd(fastTestName, startTime)

	// Test case 2: "Slow" test (we'll use a shorter duration for testing)
	slowTestName := "SlowTest"
	startTime = recorder.RecordTestStart(slowTestName)
	time.Sleep(200 * time.Millisecond) // Simulate longer test
	recorder.RecordTestEnd(slowTestName, startTime)

	// Write report
	if err := recorder.WriteReport(); err != nil {
		t.Errorf("Failed to write report: %v", err)
	}

	// Verify report exists
	reportPath := filepath.Join(tmpDir, "test_durations.txt")
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Errorf("Report file was not created")
	}

	// Verify durations were recorded
	if len(recorder.durations) != 2 {
		t.Errorf("Expected 2 test durations, got %d", len(recorder.durations))
	}

	// Verify both tests were recorded
	if _, exists := recorder.durations[fastTestName]; !exists {
		t.Errorf("Fast test duration was not recorded")
	}
	if _, exists := recorder.durations[slowTestName]; !exists {
		t.Errorf("Slow test duration was not recorded")
	}
}
