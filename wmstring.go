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
	s := strings.ToLower(string(w))

	if !strings.HasPrefix(s, w.Prefix()) {
		return false
	}

	for i := len(w.Prefix()); i < len(s); i++ {
		ch := rune(w[i])
		switch ch {
		case '0':
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		case '6':
		case '7':
		case '8':
		case '9':
		case 'a':
		case 'b':
		case 'c':
		case 'd':
		case 'e':
		case 'f':
		case ';':
			break
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
	s := strings.ToLower(string(w))
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
		return strings.ToLower(string(w))
	}

	elements := strings.SplitN(string(w), ":", 2)
	if len(elements) < 2 {
		return strings.ToLower(string(w))
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
