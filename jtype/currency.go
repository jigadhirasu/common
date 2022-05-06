package jtype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func New(data []byte) (Currency, error) {
	data = bytes.ReplaceAll(data, []byte(`\"`), []byte(``))
	data = bytes.ReplaceAll(data, []byte(`"`), []byte(``))
	in := string(data)
	if !cashMatcher.MatchString(in) {
		return 0, fmt.Errorf(`str must match : "%s"`, regexpCurrency)
	}

	ln := 9
	if idx := strings.Index(in, "."); idx > -1 {
		ln -= len(in[idx+1:])
	}

	for i := 0; i < ln; i++ {
		in = in + "0"
	}

	r := strings.ReplaceAll(in, ".", "")
	out, err := strconv.ParseInt(r, 10, 64)
	return Currency(out), err
}

// 所有跟金額有關的整數需使用該倍率
const CurrencyUnit Currency = 1000000000

const regexpCurrency = "^-?[0-9]+.?[0-9]{0,9}$"

var cashMatcher = regexp.MustCompile(regexpCurrency)

type Currency int64

func (coin Currency) MarshalJSON() ([]byte, error) {
	cash := coin.String()
	if cash == "" {
		cash = "0"
	}
	return json.Marshal(cash)
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	cur, err := New(data)
	*c = cur
	return err
}

func (c Currency) String() string {
	dash := ""
	if c < 0 {
		dash = "-"
		c *= -1
	}
	s := fmt.Sprintf("%010d", c)
	ln := len(s)
	dot := ln - 9

	if s[dot:] == "000000000" {
		return dash + s[:dot]
	}

	rs := strings.Split(s[dot:], "")

	for i := 9; i > 0; i-- {
		if rs[i-1] == "0" {
			continue
		}

		return dash + s[:dot] + "." + s[dot:dot+i]
	}

	return dash + s[:dot] + "." + s[dot:]
}

func (c Currency) Int64() int64 {
	return int64(c)
}
