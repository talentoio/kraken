package resthttp

import (
	"kraken/internal/resthttp/dto"
	"kraken/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	ltpHandler struct {
		ltpService services.LTPService
	}
)

func NewLTPHandler(ltpService services.LTPService) *ltpHandler {
	return &ltpHandler{ltpService: ltpService}
}

// Get godoc
//
//	@Summary		Get LTP data
//	@Description	Get LTP data for all available pairs
//	@Tags			LTP
//	@Produce		json
//
//	@Success		200	{object}	dto.LTPResponse
//	@Failure		500
//	@Router			/api/v1/ltp [get]
func (h *ltpHandler) Get(c echo.Context) error {

	ltp, err := h.ltpService.LTP(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var response = &dto.LTPResponse{
		List: make([]*dto.LTP, 0, len(ltp)),
	}

	for _, pairData := range ltp {
		response.List = append(response.List, &dto.LTP{
			Pair:   pairData.Pair,
			Amount: pairData.Amount,
		})
	}

	return c.JSON(http.StatusOK, response)
}
