package db

import (
	"regexp"
	"testing"
)

type sarifSettingsFixture struct {
	id         string
	raceID     string
	insertions int
}

func (f *sarifSettingsFixture) GetSetting(key string) (string, bool, error) {
	return f.id, f.id != "", nil
}
func (f *sarifSettingsFixture) InsertSettingIfAbsent(key, value string) (bool, error) {
	f.insertions++
	if f.raceID != "" {
		f.id = f.raceID
		return false, nil
	}
	f.id = value
	return true, nil
}

func TestStoredSARIFSourceID(t *testing.T) {
	a, b := &sarifSettingsFixture{}, &sarifSettingsFixture{}
	first, err := storedSARIFSourceID(a)
	if err != nil || !regexp.MustCompile("^[a-f0-9]{32}$").MatchString(first) || len(first) != 32 {
		t.Fatalf("expected valid opaque ID, got %q (%v)", first, err)
	}
	second, err := storedSARIFSourceID(a)
	if err != nil || second != first || a.insertions != 1 {
		t.Fatalf("not stable: %q, %v, writes=%d", second, err, a.insertions)
	}
	other, err := storedSARIFSourceID(b)
	if err != nil || first == other {
		t.Fatalf("independent source shares identifier: %q, %v", other, err)
	}
	race := &sarifSettingsFixture{raceID: "1234567890abcdef1234567890abcdef"}
	if got, err := storedSARIFSourceID(race); err != nil || got != race.raceID {
		t.Fatalf("did not use persisted winner: %q, %v", got, err)
	}
	broken := &sarifSettingsFixture{id: "not-an-id"}
	if _, err := storedSARIFSourceID(broken); err == nil || broken.insertions != 0 {
		t.Fatal("must not replace an unrecognized persisted ID")
	}
}
