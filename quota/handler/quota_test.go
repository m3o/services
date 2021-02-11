package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsTimeForReset(t *testing.T) {
	tcs := []struct {
		name     string
		freq     resetFrequency
		timeStr  string
		expected bool
	}{
		{
			name:     "never",
			freq:     Never,
			timeStr:  "20210201 15:04:05",
			expected: false,
		},
		{
			name:     "neverAgain",
			freq:     Never,
			timeStr:  "20210205 15:04:05",
			expected: false,
		},
		{
			name:     "daily",
			freq:     Daily,
			timeStr:  "20210201 15:04:05",
			expected: true,
		},
		{
			name:     "dailyAgain",
			freq:     Daily,
			timeStr:  "20210205 15:04:05",
			expected: true,
		},
		{
			name:     "monthlyYes",
			freq:     Monthly,
			timeStr:  "20210201 15:04:05",
			expected: true,
		},
		{
			name:     "monthlyNo",
			freq:     Monthly,
			timeStr:  "20210205 15:04:05",
			expected: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tObj, _ := time.Parse("20060102 15:04:05", tc.timeStr)
			assert.Equal(t, tc.expected, isTimeForReset(tc.freq, tObj))
		})
	}
}
