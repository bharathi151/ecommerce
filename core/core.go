package core

import (
	"quickcommerce/models"
	"quickcommerce/repo"
)

type Core struct {
	repository *repo.Repo
}

func NewCore(repository *repo.Repo) *Core {
	return &Core{repository: repository}
}

func (c *Core) Assign() {
	for len(c.repository.OrderQueue) > 0 && len(c.repository.PartnerQueue) > 0 {
		order := c.repository.Orders[c.repository.OrderQueue[0]]
		if order == nil || order.Status != models.OrderStatusCreated {
			c.popOrder()
			continue
		}
		partner := c.repository.DeliveryPartners[c.repository.PartnerQueue[0]]
		if partner == nil || partner.Status == models.PartnerStatusBusy {
			c.popPartner()
			continue
		}

		c.popOrder()
		c.popPartner()

		order.Status = models.OrderStatusAssigned
		order.DeliveryPartnerID = partner.ID
		partner.Status = models.PartnerStatusBusy
		partner.OrderID = order.ID
	}
}

func (c *Core) popOrder() {
	c.repository.OrderQueue[0] = ""
	c.repository.OrderQueue = c.repository.OrderQueue[1:]
	if len(c.repository.OrderQueue) == 0 {
		c.repository.OrderQueue = nil
	}
}

func (c *Core) popPartner() {
	c.repository.PartnerQueue[0] = ""
	c.repository.PartnerQueue = c.repository.PartnerQueue[1:]
	if len(c.repository.PartnerQueue) == 0 {
		c.repository.PartnerQueue = nil
	}
}
