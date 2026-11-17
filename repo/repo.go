package repo

import (
	"fmt"
	"quickcommerce/models"
	"sync"
)

type Repo struct {
	mu               sync.Mutex
	Customers        map[string]models.Customer
	DeliveryPartners map[string]*models.DeliveryPartner
	Orders           map[string]*models.Order
	OrderQueue       []string
	PartnerQueue     []string
	seq              int
}

func NewRepo() *Repo {
	return &Repo{
		Customers:        make(map[string]models.Customer),
		DeliveryPartners: make(map[string]*models.DeliveryPartner),
		Orders:           make(map[string]*models.Order),
		seq:              100,
	}
}

func (r *Repo) Lock()   { r.mu.Lock() }
func (r *Repo) Unlock() { r.mu.Unlock() }

func (r *Repo) NextID() string {
	r.seq++
	return fmt.Sprintf("ORD-%d", r.seq)
}
