package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kraken/internal/config"
	"kraken/internal/gateway/dto"
	"net/http"
	"path"
)

type (
	krakenGateway struct {
		cfg config.KrakenConfig
	}

	KrakenGateway interface {
		LTP(ctx context.Context, pair string) (*dto.KrakenLTPResponse, error)
	}
)

const (
	krakenProtocol = "https://"
)

//go:generate mockgen -source=kraken.go -destination=../../tests/mocks/kraken_gateway.go -package=mocks
func NewKrakenGateway(cfg config.KrakenConfig) KrakenGateway {
	return &krakenGateway{cfg: cfg}
}

// LTP send request to Kraken resource to fetch the details for pair
//
// Resource: [https://api.kraken.com/0/public/Ticker?pair=XBTUSD]
// Documentation: [https://docs.kraken.com/rest/#tag/Spot-Market-Data/operation/getTickerInformation]
func (g *krakenGateway) LTP(ctx context.Context, pair string) (*dto.KrakenLTPResponse, error) {

	requestContext, cnc := context.WithTimeout(ctx, g.cfg.RequestTimeOut)
	defer cnc()

	urlPath := fmt.Sprintf("0/public/Ticker?pair=%s", pair)
	requestURL := path.Join(g.cfg.HOST, urlPath)
	fullUrl := fmt.Sprintf("%s%s", krakenProtocol, requestURL)
	request, err := http.NewRequestWithContext(requestContext, http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("krakenGateway.make ltp request to kraken: %w", err)
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("krakenGateway.send ltp request to kraken: %w", err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("krakenGateway.read response body: %w", err)
	}

	var response dto.KrakenLTPResponse
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, fmt.Errorf("krakenGateway.unmarshal response body: %w", err)
	}

	return &response, nil
}
