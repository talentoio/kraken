package dto

type (
	KrakenLTPResponse struct {
		Error  []string       `json:"error"`
		Result map[string]LTP `json:"result"`
	}

	LTP struct {
		Closed []string `json:"c"`
	}
)
