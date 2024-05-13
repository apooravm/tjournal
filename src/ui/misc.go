package ui

import "fmt"

func timeStrParser(timestr string) string {
	// 2024-02-04T16:17:54.361333+00:00
	y, m, d, time := timestr[0:4], timestr[5:7], timestr[8:10], timestr[11:16]
	return fmt.Sprintf("\n\n%s, %s-%s-%s", time, d, m, y)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
