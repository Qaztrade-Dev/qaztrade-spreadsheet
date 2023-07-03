package domain

type Sheet struct {
	ApplicationID string
	SheetTitle    string
	SheetID       uint64
	TotalRows     uint64
	TotalSum      float64
}

type Manager struct {
	TotalRows uint64
	TotalSum  float64
	Sheets    []*Sheet
	index     int
}

const (
	TitleDostavka          = "Затраты на доставку транспортом"
	TitleCertification     = "Затраты на сертификацию предприятия"
	TitleReklamaIKU        = "Затраты на рекламу ИКУ за рубежом"
	TitlePerevodIKU        = "Затраты на перевод каталога ИКУ"
	TitleArendaIKU         = "Затраты на аренду помещения ИКУ"
	TitleCertificationIKU  = "Затраты на сертификацию ИКУ"
	TitleDemonstrationIKU  = "Затраты на демонстрацию ИКУ"
	TitleFranchising       = "Затраты на франчайзинг"
	TitleRegistrationTovar = "Затраты на регистрацию товарных знаков"
	TitleArenda            = "Затраты на аренду"
	TitlePerevod           = "Затраты на перевод"
	TitleReklama           = "Затраты на рекламу товаров за рубежом"
	TitleUchastie          = "Затраты на участие в выставках"
	TitleUchastieIKU       = "Затраты на участие в выставках ИКУ"
	TitleSootvetstvie      = "Затраты на соответствие товаров требованиям"

	SheetsDostavkaAgg          = "expenses_dostavka_view_agg"
	SheetsCertificationAgg     = "expenses_certifikazia_predpriyatia_view_agg"
	SheetsReklamaIKUAgg        = "expenses_reklama_iku_view_agg"
	SheetsPerevodIKUAgg        = "expenses_perevod_iku_view_agg"
	SheetsArendaIKUAgg         = "expenses_arenda_iku_view_agg"
	SheetsCertificationIKUAgg  = "expenses_certifikazia_iku_view_agg"
	SheetsDemonstrationIKUAgg  = "expenses_demonstrazia_iku_view_agg"
	SheetsFranchisingAgg       = "expenses_franchaizing_view_agg"
	SheetsRegistrationTovarAgg = "expenses_registrazia_tovar_znakov_view_agg"
	SheetsArendaAgg            = "expenses_arenda_view_agg"
	SheetsPerevodAgg           = "expenses_perevod_view_agg"
	SheetsReklamaAgg           = "expenses_reklaman_view_agg"
	SheetsUchastieAgg          = "expenses_uchastie_vystavka_view_agg"
	SheetsUchastieIKUAgg       = "expenses_uchastie_vystavka_iku_view_agg"
	SheetsSootvetstvieAgg      = "expenses_sootvetstvie_tovara_view_agg"
)
