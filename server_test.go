package main

import (
	"testing"
)

func TestExtractURL(t *testing.T) {

	table := []struct {
		in    string
		out   string
		outOk bool
	}{
		{"/icon/cdn2.iconfinder.com/data/icons/antivirus-internet-security-flat/33/vpn_security-512.png", "https://cdn2.iconfinder.com/data/icons/antivirus-internet-security-flat/33/vpn_security-512.png", true},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			s, ok := extractURL(tt.in)

			if s != tt.out {
				t.Errorf("\nwant %v\ngot %v", tt.out, s)
			}

			if ok != tt.outOk {
				t.Errorf("\nwant %v\ngot %v", tt.outOk, ok)
			}
		})
	}

}
