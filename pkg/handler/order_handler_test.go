package handler

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/service"
	"github.com/zenmaster911/L0/pkg/service/mocks"
)

func testHandler_CreateOrder(t *testing.T) {
	type mockBehaviour func(s *mocks.OrderMock, reply *model.Reply)

	testTable := []struct {
		name                string
		inputArgs           model.Reply
		inputErr            error
		mockBehavior        mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
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
			mockBehavior: func(s *mocks.OrderMock, reply *model.Reply) {
				s.CreateOrderMock.Expect(reply).Return("b563feb7b2b84b6test", nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "b563feb7b2b84b6test",
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			mockService := mocks.NewOrderMock(c)
			tt.mockBehavior(mockService, &tt.inputArgs)
			services := &service.Service{Order: mockService}
			h := NewHandler(services)

		})
	}
}
