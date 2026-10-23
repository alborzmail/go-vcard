package vcard

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestParameterValues(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  []string
	}{
		{`"geo:12,34"`, []string{"geo:12,34"}},
		{`"one;two"`, []string{"one;two"}},
		{`home,"one,two",work`, []string{"home", "one,two", "work"}},
		{`"one,two","three:four"`, []string{"one,two", "three:four"}},
		{`"say ^'hello^'^^^nnext"`, []string{"say \"hello\"^\nnext"}},
		{`"literal\n\path"`, []string{`literal\n\path`}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			input := "BEGIN:VCARD\r\nVERSION:4.0\r\nADR;GEO=" + tt.input + ":;;12 Main St;Berlin;;;Germany\r\nEND:VCARD\r\n"
			card, err := NewDecoder(strings.NewReader(input)).Decode()
			if err != nil {
				t.Fatal(err)
			}
			check := func(c Card) {
				t.Helper()
				f := c.Get("ADR")
				if f == nil {
					t.Fatal("ADR lost")
				}
				if !reflect.DeepEqual(f.Params["GEO"], tt.want) {
					t.Errorf("GEO = %q, want %q", f.Params["GEO"], tt.want)
				}
				a := c.Address()
				if a.StreetAddress != "12 Main St" || a.Locality != "Berlin" || a.Country != "Germany" {
					t.Errorf("address corrupted: %+v", a)
				}
			}
			check(card)
			var buf bytes.Buffer
			if err := NewEncoder(&buf).Encode(card); err != nil {
				t.Fatal(err)
			}
			again, err := NewDecoder(&buf).Decode()
			if err != nil {
				t.Fatal(err)
			}
			check(again)
		})
	}
}

func TestFormatParameterDelimiters(t *testing.T) {
	for _, value := range []string{"geo:12,34", "one;two", "one,two"} {
		if got, want := formatParam("X", value), `X="`+value+`"`; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}
