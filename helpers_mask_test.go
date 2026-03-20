package webmoney

import (
	"testing"
)

func TestMaskXMLSecrets(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Mask secret_key",
			input:    "<request><wmid>123</wmid><secret_key>SUPER_SECRET</secret_key></request>",
			expected: "<request><wmid>123</wmid><secret_key>****</secret_key></request>",
		},
		{
			name:     "Mask multiple tags",
			input:    "<pcode>12345</pcode><sign>ABCEFG</sign>",
			expected: "<pcode>****</pcode><sign>****</sign>",
		},
		{
			name:     "Case insensitive matching",
			input:    "<PCODE>54321</PCODE><SeCreT_KeY>xyz</SeCreT_KeY>",
			expected: "<PCODE>****</PCODE><SeCreT_KeY>****</SeCreT_KeY>",
		},
		{
			name:     "No sensitive data",
			input:    "<invoice><amount>10.0</amount></invoice>",
			expected: "<invoice><amount>10.0</amount></invoice>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaskXMLSecrets(tc.input)
			if got != tc.expected {
				t.Errorf("MaskXMLSecrets() = %q, want %q", got, tc.expected)
			}
		})
	}
}