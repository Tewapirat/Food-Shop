package repository

import (
	"github.com/stretchr/testify/mock"

	"github.com/TewApirat/food-shop/pkg/orderHistory/model"
)

type OrderHistoryRepositoryMock struct {
	mock.Mock
}

func (m *OrderHistoryRepositoryMock)Add(entry model.OrderHistoryEntry) error{
	args := m.Called(entry)
	return args.Error(0)
}

func (m *OrderHistoryRepositoryMock)List() ([]model.OrderHistoryEntry, error){
	args := m.Called()
	return args.Get(0).([]model.OrderHistoryEntry), args.Error(1)
}

func (m *OrderHistoryRepositoryMock)Count() (int, error){
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}
