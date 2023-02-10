package check_parser

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	MEGovTaxApiURL = "https://mapr.tax.gov.me/ic/api/verifyInvoice"
)

type TaxGovResponse struct {
	CreatedBy             any     `json:"createdBy"`
	CreationDate          any     `json:"creationDate"`
	LastUpdatedBy         any     `json:"lastUpdatedBy"`
	LastUpdateDate        any     `json:"lastUpdateDate"`
	Active                any     `json:"active"`
	ID                    float64 `json:"id"`
	Iic                   string  `json:"iic"`
	TotalPrice            float64 `json:"totalPrice"`
	InvoiceOrderNumber    float64 `json:"invoiceOrderNumber"`
	BusinessUnit          string  `json:"businessUnit"`
	CashRegister          string  `json:"cashRegister"`
	IssuerTaxNumber       string  `json:"issuerTaxNumber"`
	DateTimeCreated       string  `json:"dateTimeCreated"`
	InvoiceRequest        any     `json:"invoiceRequest"`
	InvoiceVersion        float64 `json:"invoiceVersion"`
	Fic                   string  `json:"fic"`
	IicReference          any     `json:"iicReference"`
	IicRefIssuingDate     any     `json:"iicRefIssuingDate"`
	SupplyDateOrPeriod    any     `json:"supplyDateOrPeriod"`
	CorrectiveInvoiceType any     `json:"correctiveInvoiceType"`
	BaddeptInv            any     `json:"baddeptInv"`
	PaymentMethod         []struct {
		ID       float64 `json:"id"`
		Vouchers any     `json:"vouchers"`
		Type     string  `json:"type"`
		Amount   float64 `json:"amount"`
		CompCard any     `json:"compCard"`
		AdvIIC   any     `json:"advIIC"`
		BankAcc  any     `json:"bankAcc"`
		TypeCode string  `json:"typeCode"`
	} `json:"paymentMethod"`
	Currency any `json:"currency"`
	Seller   struct {
		IDType  string `json:"idType"`
		IDNum   string `json:"idNum"`
		Name    string `json:"name"`
		Address string `json:"address"`
		Town    string `json:"town"`
		Country any    `json:"country"`
	} `json:"seller"`
	Buyer any `json:"buyer"`
	Items []struct {
		ID                 float64 `json:"id"`
		Name               string  `json:"name"`
		Code               string  `json:"code"`
		Unit               string  `json:"unit"`
		Quantity           float64 `json:"quantity"`
		UnitPriceBeforeVat float64 `json:"unitPriceBeforeVat"`
		UnitPriceAfterVat  float64 `json:"unitPriceAfterVat"`
		Rebate             float64 `json:"rebate"`
		RebateReducing     bool    `json:"rebateReducing"`
		PriceBeforeVat     float64 `json:"priceBeforeVat"`
		PriceAfterVat      float64 `json:"priceAfterVat"`
		VatRate            float64 `json:"vatRate"`
		VatAmount          float64 `json:"vatAmount"`
		ExemptFromVat      any     `json:"exemptFromVat"`
		VoucherSold        any     `json:"voucherSold"`
		Vd                 any     `json:"vd"`
		Vsn                any     `json:"vsn"`
		Investment         bool    `json:"investment"`
	} `json:"items"`
	SameTaxes []struct {
		ID             float64 `json:"id"`
		NumberOfItems  float64 `json:"numberOfItems"`
		PriceBeforeVat float64 `json:"priceBeforeVat"`
		VatRate        float64 `json:"vatRate"`
		ExemptFromVat  any     `json:"exemptFromVat"`
		VatAmount      float64 `json:"vatAmount"`
	} `json:"sameTaxes"`
	Fees                      any     `json:"fees"`
	Approvals                 []any   `json:"approvals"`
	IicRefs                   any     `json:"iicRefs"`
	InvoiceType               string  `json:"invoiceType"`
	TypeOfInvoice             string  `json:"typeOfInvoice"`
	IsSimplifiedInvoice       bool    `json:"isSimplifiedInvoice"`
	TypeOfSelfIss             any     `json:"typeOfSelfIss"`
	InvoiceNumber             string  `json:"invoiceNumber"`
	TcrCode                   string  `json:"tcrCode"`
	TaxFreeAmt                any     `json:"taxFreeAmt"`
	MarkUpAmt                 any     `json:"markUpAmt"`
	GoodsExAmt                any     `json:"goodsExAmt"`
	TotalPriceWithoutVAT      float64 `json:"totalPriceWithoutVAT"`
	TotalVATAmount            float64 `json:"totalVATAmount"`
	TotalPriceToPay           any     `json:"totalPriceToPay"`
	OperatorCode              string  `json:"operatorCode"`
	SoftwareCode              string  `json:"softwareCode"`
	IicSignature              string  `json:"iicSignature"`
	IsReverseCharge           bool    `json:"isReverseCharge"`
	PayDeadline               any     `json:"payDeadline"`
	ParagonBlockNum           any     `json:"paragonBlockNum"`
	TaxPeriod                 any     `json:"taxPeriod"`
	BankAccNum                any     `json:"bankAccNum"`
	Note                      any     `json:"note"`
	ListOfCorrectedInvoiceIIC []any   `json:"listOfCorrectedInvoiceIIC"`
	OriginalInvoice           any     `json:"originalInvoice"`
	BadDebtInvoice            any     `json:"badDebtInvoice"`
	IssuerInVat               bool    `json:"issuerInVat"`
	BadDebt                   bool    `json:"badDebt"`
}

func (self TaxGovResponse) Total() float64 {
	return self.TotalPrice
}

func (self TaxGovResponse) Date() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000+0000", self.DateTimeCreated)
}

func MERequestPaymentInfo(requestURL string) (*TaxGovResponse, error) {
	u, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	data := url.Values{
		"dateTimeCreated": query["crtd"],
		"iic":             query["iic"],
		"tin":             query["tin"],
	}
	response, err := http.PostForm(MEGovTaxApiURL, data)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := &TaxGovResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}
