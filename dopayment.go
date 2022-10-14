package newebpay

import (
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	DoPaymentURLForProd = "https://core.newebpay.com/MPG/mpg_gateway"
	DoPaymentURLForTest = "https://ccore.newebpay.com/MPG/mpg_gateway"
)

type DoPaymentRequestTradeInfo struct {
	MerchantID      string                   `url:"MerchantID" validate:"min=1,max=15"`
	RespondType     TradeInfoRespondType     `url:"RespondType" validate:"oneof=JSON String"`
	TimeStamp       int64                    `url:"TimeStamp" validate:"number"`
	Version         string                   `url:"Version" validate:"oneof=2.0"`
	LangType        TradeInfoLangType        `url:"LangType,omitempty" validate:"omitempty,oneof=en zh-tw jp"`
	MerchantOrderNo string                   `url:"MerchantOrderNo" validate:"min=1,max=30"`
	Amt             int                      `url:"Amt" validate:"min=1,max=4294967295"`
	ItemDesc        string                   `url:"ItemDesc" validate:"min=1,max=50"`
	TradeLimit      int                      `url:"TradeLimit,omitempty" validate:"min=0,max=900"`
	ExpireDate      string                   `url:"ExpireDate,omitempty" validate:"omitempty,datetime=20060102"`
	ReturnURL       string                   `url:"ReturnURL,omitempty" validate:"omitempty,url,max=50"`
	NotifyURL       string                   `url:"NotifyURL,omitempty" validate:"omitempty,url,max=50"`
	CustomerURL     string                   `url:"CustomerURL,omitempty" validate:"omitempty,url,max=50"`
	ClientBackURL   string                   `url:"ClientBackURL,omitempty" validate:"omitempty,url,max=50"`
	Email           string                   `url:"Email,omitempty" validate:"omitempty,email"`
	EmailModify     TradeInfoEmailModify     `url:"EmailModify,omitempty" validate:"oneof=1 0"`
	LoginType       TradeInfoLoginType       `url:"LoginType,omitempty" validate:"oneof=1 0"`
	OrderComment    string                   `url:"OrderComment,omitempty" validate:"omitempty,max=300"`
	Credit          TradeInfoPayMethod       `url:"CREDIT,omitempty" validate:"oneof=1 0"`
	AndroidPay      TradeInfoPayMethod       `url:"ANDROIDPAY,omitempty" validate:"oneof=1 0"`
	SamsungPay      TradeInfoPayMethod       `url:"SAMSUNGPAY,omitempty" validate:"oneof=1 0"`
	LinePay         TradeInfoPayMethod       `url:"LINEPAY,omitempty" validate:"oneof=1 0"`
	ImageUrl        string                   `url:"ImageUrl,omitempty" validate:"omitempty,url"`
	InstFlag        string                   `url:"InstFlag,omitempty" validate:"omitempty,oneof=1|split=0 3 6 12 18 24 30"`
	CreditRed       TradeInfoPayMethod       `url:"CreditRed,omitempty" validate:"oneof=1 0"`
	UnionPay        TradeInfoPayMethod       `url:"UNIONPAY,omitempty" validate:"oneof=1 0"`
	WebATM          TradeInfoPayMethod       `url:"WEBATM,omitempty" validate:"oneof=1 0"`
	VACC            TradeInfoPayMethod       `url:"VACC,omitempty" validate:"oneof=1 0"`
	BankType        TradeInfoBankType        `url:"BankType,omitempty" validate:"omitempty,split=BOT HNCB FirstBank"`
	CVS             TradeInfoPayMethod       `url:"CVS,omitempty" validate:"oneof=1 0"`
	Barcode         TradeInfoPayMethod       `url:"BARCODE,omitempty" validate:"oneof=1 0"`
	ESunWallet      TradeInfoPayMethod       `url:"ESUNWALLET,omitempty" validate:"oneof=1 0"`
	TaiwanPay       TradeInfoPayMethod       `url:"TAIWANPAY,omitempty" validate:"oneof=1 0"`
	CVSCOM          TradeInfoCVSCOM          `url:"CVSCOM,omitempty" validate:"oneof=1 2 3 0"`
	EZPay           TradeInfoPayMethod       `url:"EZPAY,omitempty" validate:"oneof=1 0"`
	EZPWeChat       TradeInfoPayMethod       `url:"EZPWECHAT,omitempty" validate:"oneof=1 0"`
	EZPAlipay       TradeInfoPayMethod       `url:"EZPALIPAY,omitempty" validate:"oneof=1 0"`
	LgsType         TradeInfoPayLgsType      `url:"LgsType,omitempty" validate:"omitempty,oneof=B2C C2C"`
	NTCB            TradeInfoPayMethod       `url:"NTCB,omitempty" validate:"oneof=1 0"`
	NTCBLocate      string                   `url:"NTCBLocate,omitempty" validate:"omitempty,number,len=3"`
	NTCBStartDate   string                   `url:"NTCBStartDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	NTCBEndDate     string                   `url:"NTCBEndDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	TokenTerm       string                   `url:"TokenTerm,omitempty" validate:"omitempty,max=20"`
	TokenTermDemand TradeInfoTokenTermDemand `url:"TokenTermDemand,omitempty" validate:"omitempty,oneof=1 2 3"`
}

type DoPaymentRequestData struct {
	MerchantID  string      `url:"MerchantID" json:"MerchantID"`
	TradeInfo   string      `url:"TradeInfo" json:"TradeInfo"`
	TradeSha    string      `url:"TradeSha" json:"TradeSha"`
	Version     string      `url:"Version" json:"Version"`
	EncryptType encryptType `url:"EncryptType" json:"EncryptType"`
}

func MakeDoPaymentRequestData(cipher Cipher, tradeInfo DoPaymentRequestTradeInfo) (DoPaymentRequestData, error) {
	tradeInfo.RespondType = TradeInfoRespondTypeJSON // TODO: For now only support JSON as RespondType
	if tradeInfo.TimeStamp == 0 {
		tradeInfo.TimeStamp = time.Now().Unix()
	}
	tradeInfo.Version = SupportedAPIVersion // This package only support API version 2.0
	err := validate.Struct(tradeInfo)
	if err != nil {
		return DoPaymentRequestData{}, err
	}
	tradeInfoForm, err := query.Values(tradeInfo)
	if err != nil {
		return DoPaymentRequestData{}, err
	}
	encodedTradeInfo := tradeInfoForm.Encode()
	encryptedTradeInfo := cipher.Encrypt(encodedTradeInfo)
	hashedTradeInfo := sha256HashHex(fmt.Sprintf("HashKey=%s&%s&HashIV=%s",
		string(cipher.hashKey), encryptedTradeInfo, string(cipher.hashIV)))
	return DoPaymentRequestData{
		MerchantID:  tradeInfo.MerchantID,
		TradeInfo:   encryptedTradeInfo,
		TradeSha:    hashedTradeInfo,
		Version:     tradeInfo.Version,
		EncryptType: encryptTypeAESCBCPKCS7Padding, // TODO: For now only support AES/CBC/PKCS7Padding
	}, nil
}

type DoPaymentResponseTradeInfo struct {
	Status  string
	Message string
	Result  DoPaymentResponseTradeInfoResult
}

type DoPaymentResponseTradeInfoResult struct {
	MerchantID        string   `json:"MerchantID"`
	Amt               int      `json:"Amt"`
	TradeNo           string   `json:"TradeNo"`
	MerchantOrderNo   string   `json:"MerchantOrderNo"`
	PaymentType       string   `json:"PaymentType"`
	RespondType       string   `json:"RespondType"`
	PayTime           *string  `json:"PayTime,omitempty"`
	IP                string   `json:"IP"`
	EscrowBank        *string  `json:"EscrowBank,omitempty"`
	AuthBank          *string  `json:"AuthBank,omitempty"`
	RespondCode       *string  `json:"RespondCode,omitempty"`
	Auth              **string `json:"Auth,omitempty"`
	Card6No           *string  `json:"Card6No,omitempty"`
	Card4No           *string  `json:"Card4No,omitempty"`
	Inst              *int     `json:"Inst,omitempty"`
	InstFirst         *int     `json:"InstFirst,omitempty"`
	InstEach          *int     `json:"InstEach,omitempty"`
	ECI               *string  `json:"ECI,omitempty"`
	TokenUseStatus    *int     `json:"TokenUseStatus,omitempty"`
	RedAmt            *int     `json:"RedAmt,omitempty"`
	PaymentMethod     *string  `json:"PaymentMethod,omitempty"`
	DCCAmt            *float64 `json:"DCC_Amt,omitempty"`
	DCCRate           *float64 `json:"DCC_Rate,omitempty"`
	DCCMarkup         *float64 `json:"DCC_Markup,omitempty"`
	DCCCurrency       *string  `json:"DCC_Currency,omitempty"`
	DCCCurrencyCode   *int     `json:"DCC_Currency_Code,omitempty"`
	PayBankCode       *string  `json:"PayBankCode,omitempty"`
	PayerAccount5Code *string  `json:"PayerAccount5Code,omitempty"`
	CodeNo            *string  `json:"CodeNo,omitempty"`
	StoreType         *int     `json:"StoreType,omitempty"`
	StoreID           *string  `json:"StoreID,omitempty"`
	Barcode1          *string  `json:"Barcode_1,omitempty"`
	Barcode2          *string  `json:"Barcode_2,omitempty"`
	Barcode3          *string  `json:"Barcode_3,omitempty"`
	RepayTimes        *int     `json:"RepayTimes,omitempty"`
	PayStore          *string  `json:"PayStore,omitempty"`
	StoreCode         *string  `json:"StoreCode,omitempty"`
	StoreName         *string  `json:"StoreName,omitempty"`
	StoreAddr         *string  `json:"StoreAddr,omitempty"`
	TradeType         *int     `json:"TradeType,omitempty"`
	CVSCOMName        *string  `json:"CVSCOMName,omitempty"`
	CVSCOMPhone       *string  `json:"CVSCOMPhone,omitempty"`
	LgsNo             *string  `json:"LgsNo,omitempty"`
	LgsType           *string  `json:"LgsType,omitempty"`
	ChannelID         *string  `json:"ChannelID,omitempty"`
	ChannelNo         *string  `json:"ChannelNo,omitempty"`
	PayAmt            *int     `json:"PayAmt,omitempty"`
	RedDisAmt         *int     `json:"RedDisAmt,omitempty"`
	ExpireDate        *string  `json:"ExpireDate,omitempty"`
	BankCode          *string  `json:"BankCode,omitempty"`
}
