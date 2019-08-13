package webmoney

import (
	"encoding/xml"
)

type MerchantRequest struct {
		XMLName						xml.Name	`xml:"merchant.request"`
		Wmid							string		`xml:"wmid"`
		LmiPayeePurse			string		`xml:"lmi_payee_purse"`
		LmiPaymentNo			string		`xml:"lmi_payment_no"`
		LmiPaymentNoType  string		`xml:"li_payment_no_type"`
		Sign							string		`xml:"sign"`
		Sha256						string		`xml:"sha256"`
		Md5								string		`xml:"md5"`
		SecretKey					string		`xml:"secret_key"`
}

type MerchantResponse struct {
	XmlName   xml.name					`xml:"merchant.request"`
	Operation MerchantOperation `xml:"operattion"`
	ErrorLog  ErrorLog					`xml:"errorlog"`
	Retval		string 						`xml:"retval"`
	Retdest   string						`xml:"retdesc"`
}

type MerchantOperation struct {
	XmlName					xml.name					`xml:"operation"`
	Wmtransid				string						`xml:"wmtransid,attr"`
	Wminvoiceid			string						`xml"wminvaiceid,attr"`
	Amount					string						`xml:"amount"`
	Operdate				string						`xml:"operdate"`
	Purpose					string						`xml:"purpose"`
	Purposefrom			string						`xml:"purposefrom"`
	Wmidfrom				string						`xml:"wmidfrom"`
	HoldPeriod			string						`xml:"hold_period"`
	HoldState				string						`xml:"hold_state"`
	Capitallerflag	string						`xml:"capitallerflag"`
	Enumflag				string						`xml:"enumflag"`
	IPAddress				string 						`xml:"IPAddress"`
	TelepayPhone		string 						`xml:"telepay_phone"`
	TelepayPaytype	string						`xml:"telepay_paytype"`
	PaymerNumber		string						`xml:"paymer_number"`
	PaymerEmail			string						`xml:"paymer_email"`
	PaymerType			string						`xml:"paymer_type"`
	CashierNumber		string						`xml:"cashier_number"`
	CashierDate			string						`xml:"cashier_date"`
	CashierAmount		string						`xml:"cashier_amount"`
	SdpType					string						`xml:"sdp_type"`
}

type ErrorLog struct {
	XmlName				xml.name						`xml:"errorlog"`
	LmiPayeePurse string							`xml:"lmi_payee_purse,attr"`
	LmiPaymentNo  string							`xml:"lmi_payment_no,attr"`
	Datecrt				string							`xml:"datecrt"`
	Dateupd				string							`xml:"dateupd"`
	DateS					string							`xml:"date_s"`
	DatePc				string							`xml:"date_pc"`
	DatePd				string							`xml:"date_pd"`
	PType					string							`xml:"p_type"`
	ErrCode				string							`xml:"err_code"`
	Siteid				string							`xml:"siteid"`
	Att						string							`xml:"att"`
	DateNotify		string							`xml:"date_notify"`
	ShopId				string							`xml:"shop_id"`
}

func (w *Wmclient) XMLTransGet(wmid, purse, no string)(MerchantResponse, error ) {
}



