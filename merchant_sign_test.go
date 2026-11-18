package webmoney

import "testing"

// По спецификации WM sha256-подпись X20/X20(TransConfirm)/X22 включает
// secret_key в конец sign-source (при classic-аутентификации s == "").
// Регрессия: раньше аргумент s игнорировался и подпись считалась без секрета.

func TestX20GetSignSourceIncludesSecret(t *testing.T) {
	x := X20Request{
		LmiPayeePurse:       "Z123456789012",
		LmiPaymentNo:        "42",
		LmiClientnumber:     "79001234567",
		LmiClientnubmerType: "phone",
	}
	want := "Z123456789012" + "42" + "79001234567" + "phone" + "topsecret"
	got, err := x.GetSignSource("topsecret")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("X20 sign source = %q, want %q", got, want)
	}
	// classic-режим: секрет не участвует
	got, err = x.GetSignSource("")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Z1234567890124279001234567phone" {
		t.Errorf("X20 classic sign source = %q", got)
	}
}

func TestX202GetSignSourceIncludesSecret(t *testing.T) {
	x := X202Request{
		LmiPayeePurse:       "Z123456789012",
		LmiWminvoiceid:      "777",
		LmiClientnumberCode: "1234",
	}
	want := "Z123456789012" + "777" + "1234" + "topsecret"
	got, err := x.GetSignSource("topsecret")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("X202 sign source = %q, want %q", got, want)
	}
	got, err = x.GetSignSource("")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Z1234567890127771234" {
		t.Errorf("X202 classic sign source = %q", got)
	}
}

func TestX22GetSignSourceIncludesSecret(t *testing.T) {
	x := X22Request{
		Validityperiodinhours: "24",
		Paymenttags: Paymenttags{
			LmiPayeePurse: "Z123456789012",
			LmiPaymentNo:  "42",
		},
	}
	want := "Z123456789012" + "42" + "24" + "topsecret"
	got, err := x.GetSignSource("topsecret")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("X22 sign source = %q, want %q", got, want)
	}
	got, err = x.GetSignSource("")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Z1234567890124224" {
		t.Errorf("X22 classic sign source = %q", got)
	}
}

func TestX18GetSignSourceIncludesSecret(t *testing.T) {
	x := X18Request{
		LmiPayeePurse: "Z123456789012",
		LmiPaymentNo:  "42",
	}
	got, err := x.GetSignSource("topsecret")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Z12345678901242topsecret" {
		t.Errorf("X18 sign source = %q", got)
	}
}