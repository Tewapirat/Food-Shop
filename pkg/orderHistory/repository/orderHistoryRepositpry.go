package repository

import "github.com/TewApirat/food-shop/pkg/orderHistory/model"

type OrderHistoryRepository interface {
	Add(entry model.OrderHistoryEntry) error
	List() ([]model.OrderHistoryEntry, error)
	Count() (int, error)
}
