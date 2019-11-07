/*
 * Gate API v4
 *
 * APIv4 provides spot, margin and futures trading operations. There are public APIs to retrieve the real-time market statistics, and private APIs which needs authentication to trade on user's behalf.
 *
 * Contact: support@mail.gate.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gateapi

type FuturesOrderBook struct {
	// Asks order depth
	Asks []FuturesOrderBookItem `json:"asks"`
	// Bids order depth
	Bids []FuturesOrderBookItem `json:"bids"`
}
