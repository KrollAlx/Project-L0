package handler

//func TestOrdersHandler_GetOrder(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//
//	expected := []models.Order{
//		{
//			Id: 1,
//			OrderUID: "b563feb7b2b84b6test",
//			TrackNumber: "WBILMTESTTRACK",
//			Entry: "WBIL",
//			Delivery: &models.Delivery{
//				Name: "Test Testov",
//				Phone: "+9720000000",
//				Zip: "2639809",
//				City: "Kiryat Mozkin",
//				Address: "Ploshad Mira 15",
//				Region: "Kraiot",
//				Email: "test@gmail.com",
//			},
//			Payment: &models.Payment{
//				Transaction: "b563feb7b2b84b6test",
//				RequestId: "",
//				Currency: "USD",
//				Provider: "wbpay",
//				Amount: 1817,
//				PaymentDt: 1637907727,
//				Bank: "alpha",
//				DeliveryCost: 1500,
//				GoodsTotal: 317,
//				CustomFee: 0,
//			},
//			Items: []models.Item{
//				{
//					ChrtId: 9934930,
//					TrackNumber: "WBILMTESTTRACK",
//					Price: 453,
//					Rid: "ab4219087a764ae0btest",
//					Name: "Mascaras",
//					Sale: 30,
//					Size: "0",
//					TotalPrice: 317,
//					NmId: 2389212,
//					Brand: "Vivienne Sabo",
//					Status: 202,
//				},
//			},
//			Locale: "en",
//			InternalSignature: "",
//			CustomerId: "test",
//			DeliveryService: "meest",
//			Shardkey: "9",
//			SmId: 99,
//			DateCreated: "2021-11-26T06:22:19Z",
//			OofShard: "1",
//		},
//	}
//	repo := mock_repository.NewMockOrders(ctl)
//	repo.EXPECT().GetAll().Return(expected, nil)
//	s := service.New(repo)
//	err := s.RestoreCache()
//
//	h := New(s)
//
//	rec := httptest.NewRecorder()
//	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
//
//	h.GetOrder(rec, req)
//
//	res := rec.Result()
//
//	require.NoError(t, err)
//	require.Equal(t, "200 OK", res.Status)
//}
