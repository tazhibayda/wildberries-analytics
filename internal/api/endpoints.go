package api

var WBBaseURLs = map[string]string{
	"statistics": "https://statistics-api.wildberries.ru/api/v1/supplier",
	"catalog":    "https://suppliers-api.wildberries.ru/api/v3",
	"advert":     "https://advert-api.wb.ru/adv/v0",
	"analytics":  "https://seller-analytics-api.wildberries.ru/api/v1/supplier",
	"finance":    "https://suppliers-api.wildberries.ru/api/v2",
	"search":     "https://catalog-analytics.wildberries.ru/api/v1",
}

type WBEndpoint struct {
	Name string
	URL  string
}

var WBEndpoints = struct {
	// === Statistics ===
	Sales  WBEndpoint
	Orders WBEndpoint
	Stocks WBEndpoint

	// === Paid Storage ===
	PaidStorageStart    WBEndpoint
	PaidStorageStatus   WBEndpoint
	PaidStorageDownload WBEndpoint

	// === Tariffs / Prices ===
	Prices  WBEndpoint
	Tariffs WBEndpoint

	// === Advertising ===
	AdvertCampaigns    WBEndpoint
	AdvertFullStats    WBEndpoint
	AdvertAutoStat     WBEndpoint
	AdvertStatWords    WBEndpoint
	AdvertKeywordsStat WBEndpoint

	// === Finance ===
	FinanceOps WBEndpoint
	Returns    WBEndpoint
	Supplies   WBEndpoint
}{
	Sales:  WBEndpoint{"sales", WBBaseURLs["statistics"] + "/sales"},
	Orders: WBEndpoint{"orders", WBBaseURLs["statistics"] + "/orders"},
	Stocks: WBEndpoint{"stocks", WBBaseURLs["statistics"] + "/stocks"},

	PaidStorageStart:    WBEndpoint{"paid_storage", WBBaseURLs["analytics"] + "/paidStorage"},
	PaidStorageStatus:   WBEndpoint{"paid_storage_status", WBBaseURLs["analytics"] + "/paidStorage/status/%s"},
	PaidStorageDownload: WBEndpoint{"paid_storage_download", WBBaseURLs["analytics"] + "/paidStorage/download/%s"},

	Prices:  WBEndpoint{"prices", WBBaseURLs["catalog"] + "/prices"},
	Tariffs: WBEndpoint{"tariffs", WBBaseURLs["catalog"] + "/tariffs"},

	AdvertCampaigns:    WBEndpoint{"advert_campaigns", WBBaseURLs["advert"] + "/campaigns"},
	AdvertFullStats:    WBEndpoint{"advert_fullstats", WBBaseURLs["advert"] + "/fullstats"},
	AdvertAutoStat:     WBEndpoint{"advert_auto_stat_words", WBBaseURLs["advert"] + "/auto/stat/words"},
	AdvertStatWords:    WBEndpoint{"advert_stat_words", WBBaseURLs["advert"] + "/stat/words"},
	AdvertKeywordsStat: WBEndpoint{"advert_keywords_stat", WBBaseURLs["advert"] + "/keywords/stats"},

	FinanceOps: WBEndpoint{"finance_ops", WBBaseURLs["finance"] + "/finances/operations"},
	Returns:    WBEndpoint{"returns", WBBaseURLs["finance"] + "/returns"},
	Supplies:   WBEndpoint{"supplies", WBBaseURLs["finance"] + "/supplies"},
}
