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
	TitleЗатратыНаДоставкуТранспортом            = "Затраты на доставку транспортом"
	TitleЗатратыНаСертификациюПредприятия        = "Затраты на сертификацию предприятия"
	TitleЗатратыНаРекламуИкуЗаРубежом            = "Затраты на рекламу ИКУ за рубежом"
	TitleЗатратыНаПереводКаталогаИку             = "Затраты на перевод каталога ИКУ"
	TitleЗатратыНаАрендуПомещенияИку             = "Затраты на аренду помещения ИКУ"
	TitleЗатратыНаСертификациюИку                = "Затраты на сертификацию ИКУ"
	TitleЗатратыНаДемонстрациюИку                = "Затраты на демонстрацию ИКУ"
	TitleЗатратыНаФранчайзинг                    = "Затраты на франчайзинг"
	TitleЗатратыНаРегистрациюТоварныхЗнаков      = "Затраты на регистрацию товарных знаков"
	TitleЗатратыНаАренду                         = "Затраты на аренду"
	TitleЗатратыНаПеревод                        = "Затраты на перевод"
	TitleЗатратыНаРекламуТоваровЗаРубежом        = "Затраты на рекламу товаров за рубежом"
	TitleЗатратыНаУчастиеВВыставках              = "Затраты на участие в выставках"
	TitleЗатратыНаУчастиеВВыставкахИку           = "Затраты на участие в выставках ИКУ"
	TitleЗатратыНаСоответствиеТоваровТребованиям = "Затраты на соответствие товаров требованиям"

	SheetsЗатратыНаДоставкуТранспортом            = "expenses_dostavka_view_agg"
	SheetsЗатратыНаСертификациюПредприятия        = "expenses_certifikazia_predpriyatia_view_agg"
	SheetsЗатратыНаРекламуИкуЗаРубежом            = "expenses_reklama_iku_view_agg"
	SheetsЗатратыНаПереводКаталогаИку             = "expenses_perevod_iku_view_agg"
	SheetsЗатратыНаАрендуПомещенияИку             = "expenses_arenda_iku_view_agg"
	SheetsЗатратыНаСертификациюИку                = "expenses_certifikazia_iku_view_agg"
	SheetsЗатратыНаДемонстрациюИку                = "expenses_demonstrazia_iku_view_agg"
	SheetsЗатратыНаФранчайзинг                    = "expenses_franchaizing_view_agg"
	SheetsЗатратыНаРегистрациюТоварныхЗнаков      = "expenses_registrazia_tovar_znakov_view_agg"
	SheetsЗатратыНаАренду                         = "expenses_arenda_view_agg"
	SheetsЗатратыНаПеревод                        = "expenses_perevod_view_agg"
	SheetsЗатратыНаРекламуТоваровЗаРубежом        = "expenses_reklaman_view_agg"
	SheetsЗатратыНаУчастиеВВыставках              = "expenses_uchastie_vystavka_view_agg"
	SheetsЗатратыНаУчастиеВВыставкахИку           = "expenses_uchastie_vystavka_iku_view_agg"
	SheetsЗатратыНаСоответствиеТоваровТребованиям = "expenses_sootvetstvie_tovara_view_agg"
)
