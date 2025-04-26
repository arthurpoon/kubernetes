package metrics

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

// TestDurationRecorder records the duration of tests and helps identify which tests
// should be marked as [Slow] based on the 5-minute threshold.
type TestDurationRecorder struct {
	mu        sync.Mutex
	durations map[string]time.Duration
	outputDir string
}

// NewTestDurationRecorder creates a new TestDurationRecorder
func NewTestDurationRecorder(outputDir string) *TestDurationRecorder {
	return &TestDurationRecorder{
		durations: make(map[string]time.Duration),
		outputDir: outputDir,
	}
}

// RecordTestStart records the start time of a test
func (r *TestDurationRecorder) RecordTestStart(testName string) time.Time {
	return time.Now()
}

// RecordTestEnd records the end time and calculates duration
func (r *TestDurationRecorder) RecordTestEnd(testName string, startTime time.Time) {
	duration := time.Since(startTime)
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.durations[testName] = duration
}

// WriteReport writes a report of test durations and slow tag recommendations
func (r *TestDurationRecorder) WriteReport() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	reportPath := filepath.Join(r.outputDir, "test_durations.txt")
	f, err := os.Create(reportPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %v", err)
	}
	defer f.Close()

	slowThreshold := 5 * time.Minute
	
	fmt.Fprintf(f, "Test Duration Report\n")
	fmt.Fprintf(f, "==================\n\n")
	fmt.Fprintf(f, "Tests that took longer than %v (should be marked [Slow]):\n\n", slowThreshold)

	for testName, duration := range r.durations {
		if duration >= slowThreshold {
			fmt.Fprintf(f, "- %s: %.2f minutes\n", testName, duration.Minutes())
		}
	}

	fmt.Fprintf(f, "\nTests that took less than %v:\n\n", slowThreshold)
	for testName, duration := range r.durations {
		if duration < slowThreshold {
			fmt.Fprintf(f, "- %s: %.2f minutes\n", testName, duration.Minutes())
		}
	}

	klog.Infof("Test duration report written to: %s", reportPath)
	return nil
}
