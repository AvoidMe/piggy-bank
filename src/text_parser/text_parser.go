package textparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrWrongFormat = errors.New("Wrong invoice format")
)

type TextInvoice struct {
	Amount float64
	Reason string
}

func groupmap(values []string, r *regexp.Regexp) map[string]string {
	keys := r.SubexpNames()

	d := make(map[string]string)
	for i := 1; i < len(keys); i++ {
		d[keys[i]] = strings.TrimSpace(values[i])
	}

	return d
}

func Parse(text string) (*TextInvoice, error) {
	ans := &TextInvoice{}
	// 10.3 something
	re := regexp.MustCompile("^(?P<amount>[0-9]+[.]?[0-9]*)(?P<reason>.*)")
	match := re.FindStringSubmatch(text)
	if len(match) != 3 {
		// something 10.3
		re = regexp.MustCompile("(?P<reason>.*?)(?P<amount>[0-9]+[.]?[0-9]*)$")
		match = re.FindStringSubmatch(text)
	}
	if len(match) != 3 {
		return nil, fmt.Errorf("%w: %s", ErrWrongFormat, text)
	}
	result := groupmap(match, re)
	amount, err := strconv.ParseFloat(result["amount"], 64)
	if err != nil {
		return nil, err
	}
	ans.Amount = amount
	ans.Reason = result["reason"]
	return ans, nil
}
