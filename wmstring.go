package webmoney

import "fmt"
import "strings"
import "strconv"

type WmString string

func (w WmString) IsCompatibleWithWmXml() bool {
	for _, ch := range w {
		if ch >= '\u0400' && ch <= '\u052F' || ch >= '\u0000' && ch <= '\u007F' || ch == '\u2116' || ch == '\u2212' || ch == '\u00AB' || ch == '\u00BB' || ch == '\u2010' || ch == '\u2013' || ch == '\u2014' {
			continue
		}
		return false
	}
	return true
}

func (w WmString) IsWmEncoded() bool {
	s := string(w)
	prefix := w.Prefix()
	if len(s) < len(prefix) || !strings.EqualFold(s[:len(prefix)], prefix) {
		return false
	}

	// Hex digits are case-insensitive; iterate runes of the lowercased
	// payload (byte indexing into a ToLower result can shift or panic).
	for _, ch := range strings.ToLower(s[len(prefix):]) {
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', ';':
		default:
			return false
		}
	}

	return true
}

func (w WmString) Prefix() string {
	return "u:"
}

func (w WmString) ToWmString() string {
	s := string(w)
	if w.IsCompatibleWithWmXml() {
		s = strings.ReplaceAll(s, string("\u2212"), "-")
		s = strings.ReplaceAll(s, "—", "-")
		return s
	}
	DebugLog("String in WmString.ToWmString", s)
	var builder strings.Builder
	builder.WriteString(w.Prefix())
	for _, c := range s {
		fmt.Fprintf(&builder, "%x;", c)
	}
	return builder.String()
}

func (w WmString) FromWmString() string {

	if !w.IsWmEncoded() {
		return string(w)
	}

	elements := strings.SplitN(string(w), ":", 2)
	if len(elements) < 2 {
		return string(w)
	}
	el := strings.Split(elements[1], ";")
	var builder strings.Builder

	for _, i := range el {
		if len(i) == 0 {
			continue
		}
		if s, err := strconv.ParseInt(i, 16, 64); err == nil {
			builder.WriteRune(rune(s))
		}
	}

	return builder.String()
}
