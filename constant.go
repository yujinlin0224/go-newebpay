package newebpay

const (
	SupportedAPIVersion = "2.0"
)

type TradeInfoRespondType string

const (
	TradeInfoRespondTypeJSON   TradeInfoRespondType = "JSON"
	TradeInfoRespondTypeString TradeInfoRespondType = "String"
)

type TradeInfoLangType string

const (
	TradeInfoLangTypeEN   TradeInfoLangType = "en"
	TradeInfoLangTypeZHTW TradeInfoLangType = "zh-tw"
	TradeInfoLangTypeJP   TradeInfoLangType = "jp"
)

type TradeInfoEmailModify int

const (
	TradeInfoEmailModifyCanBeModified    TradeInfoEmailModify = 1
	TradeInfoEmailModifyCannotBeModified TradeInfoEmailModify = 0
)

type TradeInfoLoginType int

const (
	TradeInfoLoginTypeLoginIsRequired    TradeInfoLoginType = 1
	TradeInfoLoginTypeLoginIsNotRequired TradeInfoLoginType = 0
)

type TradeInfoPayMethod int

const (
	TradeInfoPayMethodEnabled    = 1
	TradeInfoPayMethodNotEnabled = 0
)

type TradeInfoBankType string

const (
	TradeInfoBankTypeBOT       TradeInfoBankType = "BOT"
	TradeInfoBankTypeHNCB      TradeInfoBankType = "HNCB"
	TradeInfoBankTypeFirstBank TradeInfoBankType = "FirstBank"
)

type TradeInfoCVSCOM int

const (
	TradeInfoCVSComEnabledWithoutPay        TradeInfoCVSCOM = 1
	TradeInfoCVSComEnabledWithPay           TradeInfoCVSCOM = 2
	TradeInfoCVSComEnabledWithAndWithoutPay TradeInfoCVSCOM = 3
	TradeInfoCVSComNotEnabled               TradeInfoCVSCOM = 0
)

type TradeInfoPayLgsType string

const (
	TradeInfoPayLgsTypeB2C TradeInfoPayLgsType = "B2C"
	TradeInfoPayLgsTypeC2C TradeInfoPayLgsType = "C2C"
)

type TradeInfoTokenTermDemand int

const (
	TradeInfoTokenTermDemandRequireExpirationDateAndCSC TradeInfoTokenTermDemand = 1
	TradeInfoTokenTermDemandRequireExpirationDate       TradeInfoTokenTermDemand = 2
	TradeInfoTokenTermDemandRequireCSC                  TradeInfoTokenTermDemand = 3
)

type encryptType int

const (
	encryptTypeAESGCM             = 1
	encryptTypeAESCBCPKCS7Padding = 0
)
