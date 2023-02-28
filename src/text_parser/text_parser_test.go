package textparser

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	type ParseTestCase struct {
		text    string
		err     error
		invoice *TextInvoice
	}

	cases := []ParseTestCase{
		{"10.53 something", nil, &TextInvoice{10.53, "something"}},
		{"something 10.53", nil, &TextInvoice{10.53, "something"}},
		{"something 10.", nil, &TextInvoice{10, "something"}},
		{"10. something", nil, &TextInvoice{10, "something"}},
		{"10.1 for food", nil, &TextInvoice{10.1, "for food"}},
		{"for food 10.1", nil, &TextInvoice{10.1, "for food"}},
		{"for 2 bags of food 10.1", nil, &TextInvoice{10.1, "for 2 bags of food"}},
		{"10.1 for 2 bags of food", nil, &TextInvoice{10.1, "for 2 bags of food"}},
		{"10.1", nil, &TextInvoice{10.1, ""}},
		{"10", nil, &TextInvoice{10, ""}},
	}

	for _, v := range cases {
		ans, err := Parse(v.text)
		if err != nil {
			t.Errorf("Unexpected error: %v, text: %s", err, v.text)
			continue
		}
		if *ans != *v.invoice {
			t.Errorf("Text: %s, Expected: %v, got: %v", v.text, v.invoice, ans)
		}
	}

	errorCases := []ParseTestCase{
		{"for food", ErrWrongFormat, nil},
		{"for 2 food", ErrWrongFormat, nil},
		{"for 2.3 food", ErrWrongFormat, nil},
		{"", ErrWrongFormat, nil},
	}

	for _, v := range errorCases {
		ans, err := Parse(v.text)
		if ans != nil {
			t.Errorf("Unexpected result: %v, text: %s", *ans, v.text)
			continue
		}
		if !errors.Is(err, v.err) {
			t.Errorf("Text: %s, Expected: %v, got: %v", v.text, v.err, err)
		}
	}
}
