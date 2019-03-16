package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	connectDB()
	os.Exit(m.Run())
}


var testsTripCountByMedallionsUncached = map[string]struct {
	medallions []string
	date       string
	expected   []MedallionsCount
}{
	"empty medallions, date": {
		[]string{},
		"",
		[]MedallionsCount{},
	},
	"empty medallions": {
		[]string{},
		"2013-12-01",
		[]MedallionsCount{},
	},
	"empty date": {
		[]string{"2B1A06E9228B7278227621EF1B879A1D", "CFC043F3E41A505744D0FF5E63D007DD"},
		"",
		[]MedallionsCount{},
	},
	"date out of range": {
		[]string{"2B1A06E9228B7278227621EF1B879A1D", "CFC043F3E41A505744D0FF5E63D007DD"},
		"2003-12-01",
		[]MedallionsCount{},
	},
	"good query": {
		[]string{"2B1A06E9228B7278227621EF1B879A1D", "CFC043F3E41A505744D0FF5E63D007DD"},
		"2013-12-01",
		[]MedallionsCount{{"2B1A06E9228B7278227621EF1B879A1D", 4}, {"CFC043F3E41A505744D0FF5E63D007DD", 2}},
	},
}

func TestTripCountByMedallionsBypass(t *testing.T) {
	trips := NewTrips()
	for testname, testrow := range testsTripCountByMedallionsUncached {
		result := trips.CountByMedallionsBypass(testrow.medallions, testrow.date)
		assert.ElementsMatch(t, testrow.expected, result, testname)
	}
}
