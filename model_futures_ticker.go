/*
 * Gate API v4
 *
 * APIv4 provides spot, margin and futures trading operations. There are public APIs to retrieve the real-time market statistics, and private APIs which needs authentication to trade on user's behalf.
 *
 * Contact: support@mail.gate.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gateapi

type FuturesTicker struct {
	// Futures contract
	Contract string `json:"contract,omitempty"`
	// Last trading price
	Last string `json:"last,omitempty"`
	// Change percentage.
	ChangePercentage string `json:"change_percentage,omitempty"`
	// Contract total size
	TotalSize string `json:"total_size,omitempty"`
	// Lowest trading price in recent 24h
	Low24h string `json:"low_24h,omitempty"`
	// Highest trading price in recent 24h
	High24h string `json:"high_24h,omitempty"`
	// Trade size in recent 24h
	Volume24h string `json:"volume_24h,omitempty"`
	// Trade volumes in recent 24h in BTC(deprecated, use `volume_24h_base`, `volume_24h_quote`, `volume_24h_settle` instead)
	Volume24hBtc string `json:"volume_24h_btc,omitempty"`
	// Trade volumes in recent 24h in USD(deprecated, use `volume_24h_base`, `volume_24h_quote`, `volume_24h_settle` instead)
	Volume24hUsd string `json:"volume_24h_usd,omitempty"`
	// Trade volume in recent 24h, in base currency
	Volume24hBase string `json:"volume_24h_base,omitempty"`
	// Trade volume in recent 24h, in quote currency
	Volume24hQuote string `json:"volume_24h_quote,omitempty"`
	// Trade volume in recent 24h, in settle currency
	Volume24hSettle string `json:"volume_24h_settle,omitempty"`
	// Recent mark price
	MarkPrice string `json:"mark_price,omitempty"`
	// Funding rate
	FundingRate string `json:"funding_rate,omitempty"`
	// Indicative Funding rate in next period
	FundingRateIndicative string `json:"funding_rate_indicative,omitempty"`
	// Index price
	IndexPrice string `json:"index_price,omitempty"`
	// Exchange rate of base currency and settlement currency in Quanto contract. Not existed in contract of other types
	QuantoBaseRate string `json:"quanto_base_rate,omitempty"`
}
