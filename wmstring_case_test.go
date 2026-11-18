package webmoney

import "testing"

// Регрессия: ToWmString/FromWmString применяли strings.ToLower ко всей
// строке и теряли регистр пользовательских данных («Invoice ABC-123» →
// «invoice abc-123»). Регистр должен сохраняться.

func TestToWmStringPreservesCase(t *testing.T) {
	cases := []string{
		"Invoice ABC-123",
		"Оплата Заказа №42",
		"Hello World",
	}
	for _, in := range cases {
		if got := WmString(in).ToWmString(); got != in {
			t.Errorf("ToWmString(%q) = %q, want unchanged", in, got)
		}
	}
}

func TestToWmStringStillReplacesDashes(t *testing.T) {
	in := "A\u2212B—C"
	want := "A-B-C"
	if got := WmString(in).ToWmString(); got != want {
		t.Errorf("ToWmString(%q) = %q, want %q", in, got, want)
	}
}

func TestWmStringRoundTripPreservesCase(t *testing.T) {
	// Строка вне WM-совместимого диапазона кодируется в u:<hex>;
	// декодирование должно вернуть исходную строку с регистром.
	in := "Привет 🌍 MiXeD"
	encoded := WmString(in).ToWmString()
	if !WmString(encoded).IsWmEncoded() {
		t.Fatalf("encoded form %q not recognized as WM-encoded", encoded)
	}
	if got := WmString(encoded).FromWmString(); got != in {
		t.Errorf("round trip = %q, want %q", got, in)
	}
}

func TestFromWmStringPlainPreservesCase(t *testing.T) {
	in := "Hello World"
	if got := WmString(in).FromWmString(); got != in {
		t.Errorf("FromWmString(%q) = %q, want unchanged", in, got)
	}
}

func TestIsWmEncodedUppercaseHex(t *testing.T) {
	if !WmString("u:4A;").IsWmEncoded() {
		t.Error("uppercase hex digits must be recognized as WM-encoded")
	}
	if got := WmString("u:4A;").FromWmString(); got != "J" {
		t.Errorf("FromWmString(u:4A;) = %q, want J", got)
	}
	if !WmString("U:41;").IsWmEncoded() {
		t.Error("prefix check must stay case-insensitive")
	}
}

func TestIsWmEncodedNoPanicOnShrinkingToLower(t *testing.T) {
	// U+212A (KELVIN SIGN) сворачивается в 1-байтовый 'k' — старый
	// байтовый цикл по результату ToLower мог выйти за границы.
	if WmString("u:\u212A;").IsWmEncoded() {
		t.Error("kelvin sign is not a hex digit")
	}
}