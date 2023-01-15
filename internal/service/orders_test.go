package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"project-L0/internal/models"
	mock_repository "project-L0/internal/repository/mocks"
	"testing"
)

func TestOrdersService_RestoreCache(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	expected := []models.Order{
		{
			OrderUID:    "b563feb7b2b84b6test",
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: &models.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: &models.Payment{
				Transaction:  "b563feb7b2b84b6test",
				RequestId:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []models.Item{
				{
					ChrtId:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmId:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerId:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmId:              99,
			DateCreated:       "2021-11-26T06:22:19Z",
			OofShard:          "1",
		},
	}
	repo := mock_repository.NewMockOrders(ctl)
	repo.EXPECT().GetAll().Return(expected, nil)
	service := New(repo)

	err := service.RestoreCache()
	require.NoError(t, err)
	require.Equal(t, expected, service.ordersCache)
}

func TestOrdersService_RestoreCacheError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_repository.NewMockOrders(ctl)
	repo.EXPECT().GetAll().Return(nil, errors.New("repository error"))
	service := New(repo)

	err := service.RestoreCache()
	require.Error(t, err)
}

func TestOrdersService_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	newOrder := models.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: &models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: &models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtId:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmId:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}
	repo := mock_repository.NewMockOrders(ctl)
	repo.EXPECT().Create(&newOrder).Return(nil)
	service := New(repo)

	err := service.Create(&newOrder)
	require.NoError(t, err)
	require.Len(t, service.ordersCache, 1)
}

func TestOrdersService_CreateError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	newOrder := models.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: &models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: &models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtId:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmId:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}
	repo := mock_repository.NewMockOrders(ctl)
	repo.EXPECT().Create(&newOrder).Return(errors.New("DB error"))
	service := New(repo)

	err := service.Create(&newOrder)
	require.Error(t, err)
	require.Len(t, service.ordersCache, 0)
}

func TestOrdersService_Get(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cache := []models.Order{
		{
			Id:          1,
			OrderUID:    "b563feb7b2b84b6test",
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: &models.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: &models.Payment{
				Transaction:  "b563feb7b2b84b6test",
				RequestId:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []models.Item{
				{
					ChrtId:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmId:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerId:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmId:              99,
			DateCreated:       "2021-11-26T06:22:19Z",
			OofShard:          "1",
		},
	}
	repo := mock_repository.NewMockOrders(ctl)
	service := New(repo)
	service.ordersCache = cache

	order, err := service.Get(1)
	require.NoError(t, err)
	require.Equal(t, order, service.ordersCache[0])
}

func TestOrdersService_GetError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cache := []models.Order{
		{
			Id:          1,
			OrderUID:    "b563feb7b2b84b6test",
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: &models.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: &models.Payment{
				Transaction:  "b563feb7b2b84b6test",
				RequestId:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []models.Item{
				{
					ChrtId:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmId:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerId:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmId:              99,
			DateCreated:       "2021-11-26T06:22:19Z",
			OofShard:          "1",
		},
	}
	repo := mock_repository.NewMockOrders(ctl)
	service := New(repo)
	service.ordersCache = cache

	order, err := service.Get(-1)
	require.Error(t, err)
	require.Equal(t, order, models.Order{})
}
