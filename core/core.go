package core

import "quickcommerce/repo"

type Core struct {
	r *repo.Repo
}

func NewCore(r *repo.Repo) *Core {
	return &Core{r: r}
}

func (c *Core) Assign() {
	for len(c.r.OrderQueue) > 0 && len(c.r.PartnerQueue) > 0 {
		o := c.r.Orders[c.r.OrderQueue[0]]
		if o == nil || o.Status != "CREATED" {
			c.r.OrderQueue = c.r.OrderQueue[1:]
			continue
		}
		dp := c.r.DeliveryPartners[c.r.PartnerQueue[0]]
		if dp == nil || dp.Busy {
			c.r.PartnerQueue = c.r.PartnerQueue[1:]
			continue
		}

		c.r.OrderQueue = c.r.OrderQueue[1:]
		c.r.PartnerQueue = c.r.PartnerQueue[1:]

		o.Status = "ASSIGNED"
		o.DeliveryPartnerID = dp.ID
		dp.Busy = true
		dp.OrderID = o.ID
	}
}
