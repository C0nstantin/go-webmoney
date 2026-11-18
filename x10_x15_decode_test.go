package webmoney

import (
	"encoding/xml"
	"testing"
)

// Регрессия: дочерние теги ответов X10/X15 не совпадали со спецификацией WM
// (<ininvoice>, <trust>), из-за чего списки декодировались пустыми без ошибки,
// а результат SetTrust падал с "expected element type <trast> but have <trust>".

func TestX10ResponseDecodeRealFormat(t *testing.T) {
	raw := `<ininvoices cnt="2">
  <ininvoice id="101" ts="20261117 12:00:00">
    <orderid>ORD-1</orderid>
    <customerwmid>123456789012</customerwmid>
    <storepurse>Z123456789012</storepurse>
    <amount>10.00</amount>
    <desc>Test invoice</desc>
    <state>0</state>
  </ininvoice>
  <ininvoice id="102" ts="20261117 12:05:00">
    <orderid>ORD-2</orderid>
    <amount>25.50</amount>
    <state>1</state>
  </ininvoice>
</ininvoices>`
	var resp InInvoicesResponse
	if err := xml.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.InvoiceList) != 2 {
		t.Fatalf("InvoiceList len = %d, want 2", len(resp.InvoiceList))
	}
	inv := resp.InvoiceList[0]
	if inv.Id != "101" || inv.OrderId != "ORD-1" || inv.Amount != "10.00" || inv.CustomerWmid != "123456789012" {
		t.Errorf("first invoice decoded wrong: %+v", inv)
	}
	if resp.InvoiceList[1].Id != "102" || resp.InvoiceList[1].Amount != "25.50" {
		t.Errorf("second invoice decoded wrong: %+v", resp.InvoiceList[1])
	}
}

func TestX15TrustListDecodeRealFormat(t *testing.T) {
	raw := `<trustlist cnt="1">
  <trust id="555" inv="1" trans="0" purse="1" transhist="1">
    <master>123456789012</master>
    <purse>Z123456789012</purse>
    <daylimit>100.00</daylimit>
    <dsum>10.00</dsum>
  </trust>
</trustlist>`
	var resp TrustListResponse
	if err := xml.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Cnt != "1" {
		t.Errorf("Cnt = %q, want 1", resp.Cnt)
	}
	if len(resp.TrustList) != 1 {
		t.Fatalf("TrustList len = %d, want 1", len(resp.TrustList))
	}
	tr := resp.TrustList[0]
	if tr.Id != "555" || tr.Inv != "1" || tr.Trast != "0" || tr.PurseAttr != "1" {
		t.Errorf("trust attributes decoded wrong: %+v", tr)
	}
	if tr.TransHist != "1" {
		t.Errorf("TransHist = %q, want 1 (attribute per spec)", tr.TransHist)
	}
	if tr.Master != "123456789012" || tr.Purse != "Z123456789012" || tr.DayLimit != "100.00" {
		t.Errorf("trust elements decoded wrong: %+v", tr)
	}
}

func TestX15SetTrustResultDecodeRealFormat(t *testing.T) {
	raw := `<trust id="555" inv="1" trans="1" purse="1" transhist="0">
  <purse>Z123456789012</purse>
  <master>123456789012</master>
</trust>`
	var tr Trust
	if err := xml.Unmarshal([]byte(raw), &tr); err != nil {
		t.Fatal(err)
	}
	if tr.Id != "555" || tr.Purse != "Z123456789012" || tr.Master != "123456789012" {
		t.Errorf("SetTrust result decoded wrong: %+v", tr)
	}
}