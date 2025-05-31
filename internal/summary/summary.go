package summary

import "fmt"

type (
	// GrandTotal represents a breakdown of total time spent
	GrandTotal struct {
		Hours   uint32 `json:"hours"`   // Total hours component
		Minutes uint32 `json:"minutes"` // Total minutes component
		Seconds uint64 `json:"seconds"` // Total seconds component
		Text    string `json:"text"`    // Human-readable time representation
		TotalMs uint64 `json:"totalMs"` // Total time in milliseconds
	}

	TimelineItem struct {
		Start    uint64 `json:"start"`    // The Beginning timestamp of heartbeat
		Duration uint64 `json:"duration"` // The duration between heartbeat range
		Project  string `json:"project"`  // Project name
	}

	Summary struct {
		GrandTotal GrandTotal     `json:"grandTotal"`
		Timeline   []TimelineItem `json:"timeline"`
	}
)

const (
	millisecondsPerSecond = 1000
	millisecondsPerMinute = 60 * millisecondsPerSecond
	millisecondsPerHour   = 60 * millisecondsPerMinute
)

// NewGrandTotal creates and initializes a new GrandTotal instance
// Parameters:
//   - totalMs: total time in milliseconds (must be >= 0)
// Returns:
//   - *GrandTotal: initialized pointer to GrandTotal
//   - error: if totalMs is negative
func NewGrandTotal(totalMs uint64) (*GrandTotal, error) {
	if totalMs < 0 {
		return nil, fmt.Errorf("totalMs must be non-negative")
	}

	gt := &GrandTotal{
		TotalMs: totalMs,
	}

	// Calculate all time components at once
	seconds := totalMs / millisecondsPerSecond
	gt.Seconds = seconds % 60

	minutes := seconds / 60
	gt.Minutes = uint32(minutes % 60)

	gt.Hours = uint32(minutes / 60)

	gt.Text = gt.FormatDurationText(gt.Hours, gt.Minutes)

	return gt, nil
}

// getUnit returns the appropriate unit string (singular or plural)
func getUnit(value uint32, singular string) string {
	if value < 2 {
		return singular
	}
	return singular + "s"
}

// FormatDuration formats the duration in a human-readable way
// Examples:
// - 0 min
// - 5 mins
// - 1 hr 30 mins
func (gt *GrandTotal) FormatDurationText(hours uint32, minutes uint32) string {
	if gt.TotalMs <= 0 {
		return ""
	}

	switch {
	case hours > 0:
		return fmt.Sprintf("%d %s %d %s",
			hours, getUnit(hours, "hr"),
			minutes, getUnit(minutes, "min"))
	default:
		return fmt.Sprintf("%d %s", minutes, getUnit(minutes, "min"))
	}
}
