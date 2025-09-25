package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/zenmaster911/L0/pkg/cache"
	"github.com/zenmaster911/L0/pkg/cache/cache_mocks"
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/service"
	"github.com/zenmaster911/L0/pkg/service/mocks"
)

func TestHandler_CreateOrder(t *testing.T) {
	type mockBehaviour func(s *mocks.OrderMock, ca *cache_mocks.RedisCacheInterfaceMock, reply *model.Reply)

	testTable := []struct {
		name                 string
		inputArgs            model.Reply
		inputErr             error
		mockBehavior         mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputArgs: model.Reply{
				OrderUid:    "b563feb7b2b84b6test",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
				Delivery: model.Delivery{
					Name:    "test test",
					Phone:   "+1111111111",
					Zip:     "111111",
					City:    "test",
					Address: "test test test",
					Region:  "test",
					Email:   "test@tes.ru",
				},
				Payment: model.Payment{
					Transaction:  "b563feb7b2b84b6test",
					RequestId:    "",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       12345,
					PaymentDt:    1637907727,
					Bank:         "test",
					DeliveryCost: 100,
					GoodsTotal:   111,
					CustomFee:    0,
				},
				Items: []model.DeliveryItem{
					model.DeliveryItem{
						ChrtId:      11111,
						TrackNumber: "WBILMTESTTRACK",
						Price:       12345,
						Rid:         "ab4219087a764ae0btest",
						Name:        "tester",
						Sale:        1,
						Size:        "1",
						TotalPrice:  12345,
						NmId:        11111,
						Brand:       "test",
						Status:      202,
					},
				},
				Locale:            "en",
				InternalSignature: "",
				CustomerId:        "test",
				DeliveryService:   "test",
				Shardkey:          "9",
				SmId:              99,
				DateCreated:       "2021-11-26T06:22:19Z",
				OofShard:          "1",
			},
			inputErr: nil,
			mockBehavior: func(s *mocks.OrderMock, ca *cache_mocks.RedisCacheInterfaceMock, reply *model.Reply) {
				s.CreateOrderMock.Expect(reply).Return("b563feb7b2b84b6test", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"order_uid":"b563feb7b2b84b6test"}`,
		}, {
			name: "validation fail",
			inputArgs: model.Reply{
				OrderUid:    "b563feb7b2b84b6test",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
				Delivery: model.Delivery{
					Name:    "test test",
					Phone:   "",
					Zip:     "111111",
					City:    "test",
					Address: "test test test",
					Region:  "test",
					Email:   "test@tes.ru",
				},
				Payment: model.Payment{
					Transaction:  "b563feb7b2b84b6test",
					RequestId:    "",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       12345,
					PaymentDt:    1637907727,
					Bank:         "test",
					DeliveryCost: 100,
					GoodsTotal:   111,
					CustomFee:    0,
				},
				Items: []model.DeliveryItem{
					model.DeliveryItem{
						ChrtId:      11111,
						TrackNumber: "WBILMTESTTRACK",
						Price:       12345,
						Rid:         "ab4219087a764ae0btest",
						Name:        "tester",
						Sale:        1,
						Size:        "1",
						TotalPrice:  12345,
						NmId:        11111,
						Brand:       "test",
						Status:      202,
					},
				},
				Locale:            "en",
				InternalSignature: "",
				CustomerId:        "test",
				DeliveryService:   "test",
				Shardkey:          "9",
				SmId:              99,
				DateCreated:       "2021-11-26T06:22:19Z",
				OofShard:          "1",
			},
			inputErr:             nil,
			mockBehavior:         func(s *mocks.OrderMock, ca *cache_mocks.RedisCacheInterfaceMock, reply *model.Reply) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Validation failed","fields":["Phone is required"]}`,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			mockService := mocks.NewOrderMock(c)
			mockCache := cache_mocks.NewRedisCacheInterfaceMock(c)
			tt.mockBehavior(mockService, mockCache, &tt.inputArgs)
			services := &service.Service{Order: mockService}
			cache := &cache.Cache{RedisCacheInterface: mockCache}
			h := NewHandler(services, cache)

			router := chi.NewRouter()
			router.Post("/order", h.CreateOrder)
			w := httptest.NewRecorder()

			body, _ := json.Marshal(tt.inputArgs)
			req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String(), w.Body.String())

		})
	}
}

func TestHandler_GetOrder(t *testing.T) {
	type mockBehaviour func(s *mocks.OrderMock, ca *cache_mocks.RedisCacheInterfaceMock, reply *model.Reply)

	testTable := []struct {
		name         string
		inputArgs    model.Reply
		inputErr     error
		mockBehavior mockBehaviour
	}
}
