package dto

type (
	LTPResponse struct {
		List []*LTP `json:"ltp"`
	}

	LTP struct {
		Pair   string `json:"pair"`
		Amount string `json:"amount"`
	}
)
