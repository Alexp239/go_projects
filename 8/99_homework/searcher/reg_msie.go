// Code generated by re2dfa (https://github.com/opennota/re2dfa).

package main

import "unicode/utf8"

func matchMSIE(s string) (end int) {
	end = -1
	var r rune
	var rlen int
	i := 0
	_, _, _ = r, rlen, i
	r, rlen = utf8.DecodeRuneInString(s[i:])
	if rlen == 0 {
		return
	}
	i += rlen
	switch {
	case r == 77:
		goto s2
	}
	return
s2:
	r, rlen = utf8.DecodeRuneInString(s[i:])
	if rlen == 0 {
		return
	}
	i += rlen
	switch {
	case r == 83:
		goto s3
	}
	return
s3:
	r, rlen = utf8.DecodeRuneInString(s[i:])
	if rlen == 0 {
		return
	}
	i += rlen
	switch {
	case r == 73:
		goto s4
	}
	return
s4:
	r, rlen = utf8.DecodeRuneInString(s[i:])
	if rlen == 0 {
		return
	}
	i += rlen
	switch {
	case r == 69:
		end = i
	}
	return
}
