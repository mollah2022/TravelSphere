package utils

import "fmt"

// FormatPopulation converts a population number into a human-readable format
func FormatPopulation(pop int64) string {
	switch {
	case pop >= 1_000_000_000:
		return fmt.Sprintf("%.2fB", float64(pop)/1_000_000_000)
	case pop >= 1_000_000:
		return fmt.Sprintf("%.2fM", float64(pop)/1_000_000)
	case pop >= 1_000:
		return fmt.Sprintf("%.2fK", float64(pop)/1_000)
	default:
		return fmt.Sprintf("%d", pop)
	}
}
