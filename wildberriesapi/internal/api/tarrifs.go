package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// TariffItem — структура одной записи из WB API тарифов
type TariffItem struct {
	KgvpBooking         float64 `json:"kgvpBooking"`
	KgvpMarketplace     float64 `json:"kgvpMarketplace"`
	KgvpPickup          float64 `json:"kgvpPickup"`
	KgvpSupplier        float64 `json:"kgvpSupplier"`
	KgvpSupplierExpress float64 `json:"kgvpSupplierExpress"`
	PaidStorageKgvp     float64 `json:"paidStorageKgvp"`
	ParentID            int     `json:"parentID"`
	ParentName          string  `json:"parentName"`
	SubjectID           int     `json:"subjectID"`
	SubjectName         string  `json:"subjectName"`
}

type TariffsBox struct {
	Data struct {
		DtNextBox     string `json:"dtNextBox"`
		DtTillMax     string `json:"dtTillMax"`
		WarehouseList []struct {
			BoxDeliveryBase                string `json:"boxDeliveryBase"`
			BoxDeliveryCoefExpr            string `json:"boxDeliveryCoefExpr"`
			BoxDeliveryLiter               string `json:"boxDeliveryLiter"`
			BoxDeliveryMarketplaceBase     string `json:"boxDeliveryMarketplaceBase"`
			BoxDeliveryMarketplaceCoefExpr string `json:"boxDeliveryMarketplaceCoefExpr"`
			BoxDeliveryMarketplaceLiter    string `json:"boxDeliveryMarketplaceLiter"`
			BoxStorageBase                 string `json:"boxStorageBase"`
			BoxStorageCoefExpr             string `json:"boxStorageCoefExpr"`
			BoxStorageLiter                string `json:"boxStorageLiter"`
			GeoName                        string `json:"geoName"`
			WarehouseName                  string `json:"warehouseName"`
		} `json:"warehouseList"`
	} `json:"data"`
}

type TariffsPallet struct {
	Data struct {
		DtNextPallet  string `json:"dtNextPallet"`
		DtTillMax     string `json:"dtTillMax"`
		WarehouseList []struct {
			PalletDeliveryExpr       string `json:"palletDeliveryExpr"`
			PalletDeliveryValueBase  string `json:"palletDeliveryValueBase"`
			PalletDeliveryValueLiter string `json:"palletDeliveryValueLiter"`
			PalletStorageExpr        string `json:"palletStorageExpr"`
			PalletStorageValueExpr   string `json:"palletStorageValueExpr"`
			WarehouseName            string `json:"warehouseName"`
		} `json:"warehouseList"`
	} `json:"data"`
}

// GetTariffs получает тарифы и комиссии по категориям товаров
func (c *WBClient) GetTariffs(ctx context.Context) ([]TariffItem, error) {
	allTariffs := make([]TariffItem, 0)

	url := "https://common-api.wildberries.ru/api/v1/tariffs/commission"
	body, err := c.doRequest(ctx, "GET", url, c.Tokens[0], nil)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("❌ failed to fetch tariffs")
		return nil, err
	}

	var resp struct {
		Report []TariffItem `json:"report"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		c.Logger.Error().Err(err).Msg("unmarshal error in GetTariffs response")
		return nil, err
	}

	allTariffs = append(allTariffs, resp.Report...)
	c.Logger.Info().Msgf("✅  tariffs loaded (%d records)", len(resp.Report))

	// Rate limit WB API
	time.Sleep(700 * time.Millisecond)

	return allTariffs, nil
}

// GetTariffsBox получает тарифы и комиссии по категориям товаров
func (c *WBClient) GetTariffsBox(ctx context.Context, date string) ([]TariffsBox, error) {
	allTariffs := make([]TariffsBox, 0)

	url := fmt.Sprintf("https://common-api.wildberries.ru/api/v1/tariffs/box?date=%s", date)
	body, err := c.doRequest(ctx, "GET", url, c.Tokens[0], nil)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("❌ failed to fetch tariffs")
		return nil, err
	}

	var resp struct {
		Response TariffsBox `json:"response"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		c.Logger.Error().Err(err).Msg("unmarshal error in GetTariffs response")
		return nil, err
	}

	allTariffs = append(allTariffs, resp.Response)
	c.Logger.Info().Msgf("✅  tariffs loaded (%d records)", len(resp.Response.Data.WarehouseList))

	// Rate limit WB API
	time.Sleep(700 * time.Millisecond)

	return allTariffs, nil
}

// GetTariffsPallet получает тарифы и комиссии по категориям товаров
func (c *WBClient) GetTariffsPallet(ctx context.Context, date string) ([]TariffsPallet, error) {
	allTariffs := make([]TariffsPallet, 0)

	url := fmt.Sprintf("https://common-api.wildberries.ru/api/v1/tariffs/pallet?date=%s", date)
	body, err := c.doRequest(ctx, "GET", url, c.Tokens[0], nil)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("❌ failed to fetch tariffs")
		return nil, err
	}

	var resp struct {
		Data TariffsPallet `json:"response"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		c.Logger.Error().Err(err).Msg("unmarshal error in GetTariffs response")
		return nil, err
	}

	allTariffs = append(allTariffs, resp.Data)
	c.Logger.Info().Msgf("✅  tariffs loaded (%d records)", len(resp.Data.Data.WarehouseList))

	// Rate limit WB API
	time.Sleep(700 * time.Millisecond)

	return allTariffs, nil
}
