package check_parser

import (
	"testing"
	"time"
)

type dateTestCase struct {
	asString string
	expected time.Time
}

func TestTaxGovResponseDate(t *testing.T) {
	cases := []dateTestCase{
		{
			asString: "2023-02-02T16:31:36.000+0000",
			expected: time.Unix(1675355496, 0),
		},
	}
	for _, test := range cases {
		actual := TaxGovResponse{}
		actual.DateTimeCreated = test.asString
		date, err := actual.Date()
		if err != nil {
			t.Fatalf("Can't parse required date %s, err: %s", test.asString, err.Error())
		}
		if !date.Equal(test.expected) {
			t.Fatalf("Wrong results, want: %s, got: %s", test.expected, date)
		}
	}
}
