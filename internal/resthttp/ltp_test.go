package resthttp

import (
	"encoding/json"
	"kraken/internal/resthttp/dto"
	"kraken/internal/services/domain"
	"kraken/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLtpHandler_Get(t *testing.T) {

	type args struct {
		uri string
	}

	type success struct {
		response *dto.LTPResponse
	}

	type want struct {
		statusCode int
		success    success
	}

	testTable := []struct {
		name    string
		prepare func(f *mocks.MockLTPService)
		args    args
		want    want
	}{
		{
			name: "services.LTPService return error",
			args: args{
				uri: "/",
			},
			prepare: func(f *mocks.MockLTPService) {
				f.EXPECT().LTP(gomock.Any()).Return(nil, errors.New("any error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "services.LTPService return valid data",
			args: args{
				uri: "/",
			},
			prepare: func(f *mocks.MockLTPService) {
				f.EXPECT().LTP(gomock.Any()).Return([]*domain.LTPPair{
					{
						Pair:   "pair_1",
						Amount: "1000",
					},
					{
						Pair:   "pair_2",
						Amount: "2000",
					},
				}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				success: success{response: &dto.LTPResponse{List: []*dto.LTP{
					{
						Pair:   "pair_1",
						Amount: "1000",
					},
					{
						Pair:   "pair_2",
						Amount: "2000",
					},
				}}},
			},
		},
	}

	for _, tc := range testTable {
		t.Run(
			tc.name,
			func(t *testing.T) {

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ltpService := mocks.NewMockLTPService(ctrl)

				if tc.prepare != nil {
					tc.prepare(ltpService)
				}

				handler := ltpHandler{
					ltpService: ltpService,
				}

				request := httptest.NewRequest(http.MethodGet, tc.args.uri, nil)
				response := httptest.NewRecorder()

				c := echo.New().NewContext(request, response)
				c.SetPath("")

				err := handler.Get(c)
				assert.NoError(t, err)
				assert.Equal(t, tc.want.statusCode, response.Code)

				var valueResponse *dto.LTPResponse
				if response.Code == http.StatusOK {
					err = json.NewDecoder(response.Body).Decode(&valueResponse)
					assert.NoError(t, err)
					assert.Equal(t, tc.want.success.response, valueResponse)
				}

			},
		)
	}

}
