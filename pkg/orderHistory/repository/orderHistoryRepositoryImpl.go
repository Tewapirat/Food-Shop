package repository

import (
	"sync"

	"github.com/TewApirat/food-shop/pkg/orderHistory/model"
)

type orderHistoryRepositoryImpl struct {
	mu      sync.Mutex
	entries []model.OrderHistoryEntry
}

func NewOrderHistoryRepositoryImpl() OrderHistoryRepository {
	return &orderHistoryRepositoryImpl{
		entries: make([]model.OrderHistoryEntry, 0),
	}
}

func (r *orderHistoryRepositoryImpl) Add(entry model.OrderHistoryEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry)
	return nil
}

func (r *orderHistoryRepositoryImpl) List() ([]model.OrderHistoryEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]model.OrderHistoryEntry, len(r.entries))
	copy(out, r.entries)
	return out, nil
}

func (r *orderHistoryRepositoryImpl) Count() (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries), nil
}
