package source

import (
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source/db"
	"unsafe"
)

type Campaign struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type Source struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func toCampaign(dbCampaign db.Campaign) Campaign {
	pu := (*Campaign)(unsafe.Pointer(&dbCampaign))
	return *pu
}

func toCampaignSlice(dbCampaigns []db.Campaign) []Campaign {
	users := make([]Campaign, len(dbCampaigns))
	for i, dbUsr := range dbCampaigns {
		users[i] = toCampaign(dbUsr)
	}
	return users
}
