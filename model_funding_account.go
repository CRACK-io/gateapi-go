/*
 * Gate API v4
 *
 * APIv4 provides spot, margin and futures trading operations. There are public APIs to retrieve the real-time market statistics, and private APIs which needs authentication to trade on user's behalf.
 *
 * Contact: support@mail.gate.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gateapi

type FundingAccount struct {
	// Currency name
	Currency string `json:"currency,omitempty"`
	// Available assets to lend, which is identical to spot account `available`
	Available string `json:"available,omitempty"`
	// Locked amount. i.e. amount in `open` loans
	Locked string `json:"locked,omitempty"`
	// Amount that is loaned but not repaid
	Lent string `json:"lent,omitempty"`
	// Amount used for lending. total_lent = lent + locked
	TotalLent string `json:"total_lent,omitempty"`
}
