package dpt

import (
	"strings"
	"testing"
	"time"
)

func TestDPT_19001_Pack(t *testing.T) {
	// Test case 1: Basic date and time
	dp1 := DPT_19001{
		Year:         123, // 2023 (1900 + 123)
		Month:        5,
		DayOfMonth:   15,
		DayOfWeek:    1, // Monday
		HourOfDay:    14,
		Minutes:      30,
		Seconds:      45,
		SummerTime:   true,
		ExternalSync: true,
		ReliableSync: true,
	}

	data1 := dp1.Pack()
	if len(data1) != 9 {
		t.Errorf("Expected packed data length of 9, got %d", len(data1))
	}

	// Check specific bytes
	if data1[0] != 0 { // Padding byte
		t.Errorf("Expected padding byte to be 0, got %d", data1[0])
	}
	if data1[1] != 123 { // Year
		t.Errorf("Expected Year byte to be 123, got %d", data1[1])
	}
	if data1[2] != 5 { // Month
		t.Errorf("Expected Month byte to be 5, got %d", data1[2])
	}
	if data1[3] != 15 { // Day
		t.Errorf("Expected Day byte to be 15, got %d", data1[3])
	}
	if data1[4] != (1<<5 | 14) { // DayOfWeek and Hour
		t.Errorf("Expected DayOfWeek and Hour byte to be %d, got %d", (1<<5 | 14), data1[4])
	}
	if data1[5] != 30 { // Minutes
		t.Errorf("Expected Minutes byte to be 30, got %d", data1[5])
	}
	if data1[6] != 45 { // Seconds
		t.Errorf("Expected Seconds byte to be 45, got %d", data1[6])
	}
	if data1[7] != 0x01 { // Flags (only SummerTime set)
		t.Errorf("Expected Flags byte to be 0x01, got 0x%02x", data1[7])
	}
	if data1[8] != 0xC0 { // Clock Quality (ExternalSync and ReliableSync set)
		t.Errorf("Expected Clock Quality byte to be 0xC0, got 0x%02x", data1[8])
	}

	// Test case 2: With all flags set
	dp2 := DPT_19001{
		Year:         123, // 2023 (1900 + 123)
		Month:        5,
		DayOfMonth:   15,
		DayOfWeek:    1, // Monday
		HourOfDay:    14,
		Minutes:      30,
		Seconds:      45,
		Fault:        true,
		WorkingDay:   true,
		NoWD:         true,
		NoYear:       true,
		NoDate:       true,
		NoDayOfWeek:  true,
		NoTime:       true,
		SummerTime:   true,
		ExternalSync: true,
		ReliableSync: true,
	}

	data2 := dp2.Pack()
	if data2[7] != 0xFF { // All flags set
		t.Errorf("Expected Flags byte to be 0xFF, got 0x%02x", data2[7])
	}
}

func TestDPT_19001_Unpack(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    []byte
		expected DPT_19001
		wantErr  bool
	}{
		{
			name: "Standard date and time",
			input: []byte{
				0,           // Padding byte
				123,         // Year (2023)
				5,           // Month (May)
				15,          // Day (15th)
				(1<<5 | 14), // DayOfWeek (Monday) and Hour (14)
				30,          // Minutes (30)
				45,          // Seconds (45)
				0x01,        // Flags (only SummerTime set)
				0xC0,        // Clock Quality (ExternalSync and ReliableSync set)
			},
			expected: DPT_19001{
				Year:         123,
				Month:        5,
				DayOfMonth:   15,
				DayOfWeek:    1,
				HourOfDay:    14,
				Minutes:      30,
				Seconds:      45,
				SummerTime:   true,
				ExternalSync: true,
				ReliableSync: true,
			},
			wantErr: false,
		},
		{
			name:  "April 2, 2025",
			input: []byte{0, 0x7d, 0x04, 0x02, 0x6e, 0x09, 0x1c, 0x40, 0x80},
			expected: DPT_19001{
				Year:         125,
				Month:        4,
				DayOfMonth:   2,
				DayOfWeek:    3,
				HourOfDay:    14,
				Minutes:      9,
				Seconds:      28,
				Fault:        false,
				WorkingDay:   true,
				NoWD:         false,
				NoYear:       false,
				NoDate:       false,
				NoDayOfWeek:  false,
				NoTime:       false,
				SummerTime:   false,
				ExternalSync: true,
				ReliableSync: false,
			},
			wantErr: false,
		},
		{
			name:     "Invalid data length",
			input:    []byte{1, 2, 3},
			expected: DPT_19001{},
			wantErr:  true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dp DPT_19001
			err := dp.Unpack(tt.input)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Skip further checks if we expected an error
			if tt.wantErr {
				return
			}

			// Check all fields
			if dp.Year != tt.expected.Year {
				t.Errorf("Year = %d, want %d", dp.Year, tt.expected.Year)
			}
			if dp.Month != tt.expected.Month {
				t.Errorf("Month = %d, want %d", dp.Month, tt.expected.Month)
			}
			if dp.DayOfMonth != tt.expected.DayOfMonth {
				t.Errorf("DayOfMonth = %d, want %d", dp.DayOfMonth, tt.expected.DayOfMonth)
			}
			if dp.DayOfWeek != tt.expected.DayOfWeek {
				t.Errorf("DayOfWeek = %d, want %d", dp.DayOfWeek, tt.expected.DayOfWeek)
			}
			if dp.HourOfDay != tt.expected.HourOfDay {
				t.Errorf("HourOfDay = %d, want %d", dp.HourOfDay, tt.expected.HourOfDay)
			}
			if dp.Minutes != tt.expected.Minutes {
				t.Errorf("Minutes = %d, want %d", dp.Minutes, tt.expected.Minutes)
			}
			if dp.Seconds != tt.expected.Seconds {
				t.Errorf("Seconds = %d, want %d", dp.Seconds, tt.expected.Seconds)
			}
			if dp.Fault != tt.expected.Fault {
				t.Errorf("Fault = %v, want %v", dp.Fault, tt.expected.Fault)
			}
			if dp.WorkingDay != tt.expected.WorkingDay {
				t.Errorf("WorkingDay = %v, want %v", dp.WorkingDay, tt.expected.WorkingDay)
			}
			if dp.NoWD != tt.expected.NoWD {
				t.Errorf("NoWD = %v, want %v", dp.NoWD, tt.expected.NoWD)
			}
			if dp.NoYear != tt.expected.NoYear {
				t.Errorf("NoYear = %v, want %v", dp.NoYear, tt.expected.NoYear)
			}
			if dp.NoDate != tt.expected.NoDate {
				t.Errorf("NoDate = %v, want %v", dp.NoDate, tt.expected.NoDate)
			}
			if dp.NoDayOfWeek != tt.expected.NoDayOfWeek {
				t.Errorf("NoDayOfWeek = %v, want %v", dp.NoDayOfWeek, tt.expected.NoDayOfWeek)
			}
			if dp.NoTime != tt.expected.NoTime {
				t.Errorf("NoTime = %v, want %v", dp.NoTime, tt.expected.NoTime)
			}
			if dp.SummerTime != tt.expected.SummerTime {
				t.Errorf("SummerTime = %v, want %v", dp.SummerTime, tt.expected.SummerTime)
			}
			if dp.ExternalSync != tt.expected.ExternalSync {
				t.Errorf("ExternalSync = %v, want %v", dp.ExternalSync, tt.expected.ExternalSync)
			}
			if dp.ReliableSync != tt.expected.ReliableSync {
				t.Errorf("ReliableSync = %v, want %v", dp.ReliableSync, tt.expected.ReliableSync)
			}
		})
	}
}

func TestDPT_19001_String(t *testing.T) {
	// Test with all fields valid
	dp1 := DPT_19001{
		Year:         123, // 2023 (1900 + 123)
		Month:        5,
		DayOfMonth:   15,
		DayOfWeek:    1, // Monday
		HourOfDay:    14,
		Minutes:      30,
		Seconds:      45,
		SummerTime:   true,
		ExternalSync: true,
		ReliableSync: true,
	}

	if s := dp1.String(); !contains(s, "2023-05-15") || !contains(s, "Monday") || !contains(s, "14:30:45") {
		t.Errorf("Expected String() to contain date, day and time, got %s", s)
	}

	// Test with some fields invalid
	dp2 := DPT_19001{
		Year:        123, // 2023 (1900 + 123)
		Month:       5,
		DayOfMonth:  15,
		DayOfWeek:   1, // Monday
		HourOfDay:   14,
		Minutes:     30,
		Seconds:     45,
		NoYear:      true,
		NoDayOfWeek: true,
		SummerTime:  true,
	}

	if s := dp2.String(); contains(s, "2023") || contains(s, "Monday") {
		t.Errorf("Expected String() to not contain year and day, got %s", s)
	}

	// Test with all fields invalid
	dp3 := DPT_19001{
		NoYear:      true,
		NoDate:      true,
		NoDayOfWeek: true,
		NoTime:      true,
	}

	if s := dp3.String(); !contains(s, "Invalid DateTime") {
		t.Errorf("Expected String() to indicate invalid DateTime, got %s", s)
	}
}

func TestDPT_19001_IsValid(t *testing.T) {
	// Valid date and time
	dp1 := DPT_19001{
		Year:       123, // 2023 (1900 + 123)
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
	}
	if !dp1.IsValid() {
		t.Error("Expected IsValid() to return true for valid date and time")
	}

	// Invalid month
	dp2 := DPT_19001{
		Year:       123,
		Month:      13, // Invalid month
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
	}
	if dp2.IsValid() {
		t.Error("Expected IsValid() to return false for invalid month")
	}

	// Invalid day (February 30th)
	dp3 := DPT_19001{
		Year:       123,
		Month:      2,
		DayOfMonth: 30, // Invalid day for February
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
	}
	if dp3.IsValid() {
		t.Error("Expected IsValid() to return false for invalid day")
	}

	// Invalid hour
	dp4 := DPT_19001{
		Year:       123,
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  25, // Invalid hour
		Minutes:    30,
		Seconds:    45,
	}
	if dp4.IsValid() {
		t.Error("Expected IsValid() to return false for invalid hour")
	}

	// Invalid minutes
	dp5 := DPT_19001{
		Year:       123,
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    60, // Invalid minutes
		Seconds:    45,
	}
	if dp5.IsValid() {
		t.Error("Expected IsValid() to return false for invalid minutes")
	}

	// Invalid seconds
	dp6 := DPT_19001{
		Year:       123,
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    60, // Invalid seconds
	}
	if dp6.IsValid() {
		t.Error("Expected IsValid() to return false for invalid seconds")
	}

	// Valid with NoYear flag
	dp7 := DPT_19001{
		Year:       255, // Invalid year, but NoYear flag is set
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
		NoYear:     true,
	}
	if !dp7.IsValid() {
		t.Error("Expected IsValid() to return true when NoYear flag is set")
	}

	// Valid with NoDate flag
	dp8 := DPT_19001{
		Year:       123,
		Month:      13, // Invalid month, but NoDate flag is set
		DayOfMonth: 32, // Invalid day, but NoDate flag is set
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
		NoDate:     true,
	}
	if !dp8.IsValid() {
		t.Error("Expected IsValid() to return true when NoDate flag is set")
	}

	// Valid with NoTime flag
	dp9 := DPT_19001{
		Year:       123,
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  25, // Invalid hour, but NoTime flag is set
		Minutes:    60, // Invalid minutes, but NoTime flag is set
		Seconds:    60, // Invalid seconds, but NoTime flag is set
		NoTime:     true,
	}
	if !dp9.IsValid() {
		t.Error("Expected IsValid() to return true when NoTime flag is set")
	}
}

func TestDPT_19001_ToTime(t *testing.T) {
	// Valid date and time
	dp1 := DPT_19001{
		Year:       123, // 2023 (1900 + 123)
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
	}

	timeObj := dp1.ToTime()
	if timeObj == nil {
		t.Fatal("Expected ToTime() to return a non-nil time.Time object")
	}

	if timeObj.Year() != 2023 {
		t.Errorf("Expected year to be 2023, got %d", timeObj.Year())
	}
	if timeObj.Month() != time.May {
		t.Errorf("Expected month to be May, got %s", timeObj.Month())
	}
	if timeObj.Day() != 15 {
		t.Errorf("Expected day to be 15, got %d", timeObj.Day())
	}
	if timeObj.Hour() != 14 {
		t.Errorf("Expected hour to be 14, got %d", timeObj.Hour())
	}
	if timeObj.Minute() != 30 {
		t.Errorf("Expected minute to be 30, got %d", timeObj.Minute())
	}
	if timeObj.Second() != 45 {
		t.Errorf("Expected second to be 45, got %d", timeObj.Second())
	}

	// Invalid date (NoYear flag set)
	dp2 := DPT_19001{
		Year:       123,
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
		NoYear:     true,
	}

	if dp2.ToTime() != nil {
		t.Error("Expected ToTime() to return nil when NoYear flag is set")
	}

	// Invalid date (NoDate flag set)
	dp3 := DPT_19001{
		Year:       123,
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
		NoDate:     true,
	}

	if dp3.ToTime() != nil {
		t.Error("Expected ToTime() to return nil when NoDate flag is set")
	}

	// Invalid time (NoTime flag set)
	dp4 := DPT_19001{
		Year:       123,
		Month:      5,
		DayOfMonth: 15,
		HourOfDay:  14,
		Minutes:    30,
		Seconds:    45,
		NoTime:     true,
	}

	if dp4.ToTime() != nil {
		t.Error("Expected ToTime() to return nil when NoTime flag is set")
	}
}

func TestFromTime(t *testing.T) {
	// Create a time.Time object
	timeObj := time.Date(2023, time.May, 15, 14, 30, 45, 0, time.Local)

	// Convert to DPT_19001
	dp := FromTime(timeObj)

	// Check the values
	if dp.Year != 123 { // 2023 - 1900 = 123
		t.Errorf("Expected Year to be 123, got %d", dp.Year)
	}
	if dp.Month != 5 {
		t.Errorf("Expected Month to be 5, got %d", dp.Month)
	}
	if dp.DayOfMonth != 15 {
		t.Errorf("Expected DayOfMonth to be 15, got %d", dp.DayOfMonth)
	}
	if dp.HourOfDay != 14 {
		t.Errorf("Expected HourOfDay to be 14, got %d", dp.HourOfDay)
	}
	if dp.Minutes != 30 {
		t.Errorf("Expected Minutes to be 30, got %d", dp.Minutes)
	}
	if dp.Seconds != 45 {
		t.Errorf("Expected Seconds to be 45, got %d", dp.Seconds)
	}

	// Check weekday conversion (Monday = 1 in KNX)
	expectedWeekday := uint8((int(timeObj.Weekday())+6)%7 + 1)
	if dp.DayOfWeek != expectedWeekday {
		t.Errorf("Expected DayOfWeek to be %d, got %d", expectedWeekday, dp.DayOfWeek)
	}

	// Test year limits
	// Year before 1900
	timeObj1 := time.Date(1800, time.January, 1, 0, 0, 0, 0, time.Local)
	dp1 := FromTime(timeObj1)
	if dp1.Year != 0 { // Should be clamped to 1900
		t.Errorf("Expected Year to be 0 (1900), got %d", dp1.Year)
	}

	// Year after 2155
	timeObj2 := time.Date(2200, time.January, 1, 0, 0, 0, 0, time.Local)
	dp2 := FromTime(timeObj2)
	if dp2.Year != 255 { // Should be clamped to 2155
		t.Errorf("Expected Year to be 255 (2155), got %d", dp2.Year)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
