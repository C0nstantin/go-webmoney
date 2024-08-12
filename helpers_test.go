package webmoney

import (
	"os"
	"regexp"
	"testing"
	"time"
)

func TestReqn(t *testing.T) {
	reqn := Reqn()
	// Check if the length of the output is as expected
	expectedLength := 17 // YYYYMMDDHHMMSSmmm
	if len(reqn) != expectedLength {
		t.Errorf("Expected length of reqn to be %d, got %d", expectedLength, len(reqn))
	}

	// Check if the format of the output matches the expected pattern
	// The pattern is YYYYMMDDHHMMSSmmm where Y is year, M is month, D is day,
	// H is hour, M is minute, S is second, and m is millisecond
	pattern := `^\d{17}$`
	matched, err := regexp.MatchString(pattern, reqn)
	if err != nil {
		t.Fatal("Regex match failed:", err)
	}
	if !matched {
		t.Errorf("Reqn %s does not match the expected pattern %s", reqn, pattern)
	}

	// Optionally, you can check if the time is within a reasonable range (e.g., +/- 1 minute from now)
	// This is a bit more complex and might not be strictly necessary for all use cases
	loc, err := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)
	reqnTime, err := time.ParseInLocation("20060102150405", reqn[0:14], loc)
	if err != nil {
		t.Fatal("Failed to parse reqn time:", err)
	}
	if reqnTime.Before(now.Add(-time.Minute)) || reqnTime.After(now.Add(time.Minute)) {
		t.Errorf("Reqn time %s is not within a reasonable range of the current time", reqnTime)
	}
}
func TestOldReqn(t *testing.T) {
	os.Setenv("USE_OLD_REQN", "true")
	reqn := Reqn()
	// Check if the length of the output is as expected
	if len(reqn) != 14 {
		t.Errorf("Expected length of reqn to be 14, got %d", len(reqn))
	}
	// Check if the format of the output matches the expected pattern
	pattern := `^\d{14}$`
	matched, err := regexp.MatchString(pattern, reqn)
	if err != nil {
		t.Fatal("Regex match failed:", err)
	}
	if !matched {
		t.Errorf("Reqn %s does not match the expected pattern %s", reqn, pattern)
	}
	loc, _ := time.LoadLocation("UTC")
	parse, err := time.ParseInLocation("060102150405", reqn[0:12], loc)
	if err != nil {
		t.Fatal("Failed to parse reqn time:", err)
	}
	if parse.Before(time.Now().In(loc).Add(-time.Minute)) || parse.After(time.Now().In(loc).Add(time.Minute)) {
		t.Errorf("Reqn time %s is not within a reasonable range of the current time", parse)
	}
}

// TestUtf8ToWin tests the Utf8ToWin function to ensure it correctly encodes UTF-8 strings to Windows-1251.
func TestUtf8ToWin(t *testing.T) {
	// Example UTF-8 string and its expected Windows-1251 encoding.
	// Note: Ensure the example is correct and can be encoded in Windows-1251.
	utf8String := "Привет, мир!"                                       // "Hello, world!" in Russian.
	expectedWin1251String := "\xcf\xf0\xe8\xe2\xe5\xf2, \xec\xe8\xf0!" // Expected Windows-1251 encoding.

	// Encode the UTF-8 string to Windows-1251.
	win1251String, err := Utf8ToWin(utf8String)
	if err != nil {
		t.Fatalf("Utf8ToWin returned an error: %v", err)
	}

	// Compare the encoded string with the expected result.
	if win1251String != expectedWin1251String {
		t.Errorf("Utf8ToWin failed to encode UTF-8 to Windows-1251 correctly.\nGot: %s\nWant: %s", win1251String, expectedWin1251String)
	}
}
