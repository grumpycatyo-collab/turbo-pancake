package db

import (
	"github.com/grumpycatyo-collab/turbo-pancake/business/sys/database"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Store struct {
	log          *zerolog.Logger
	db           *sqlx.DB
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *zerolog.Logger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

func (s Store) QueryCampaignsBySourceID(sourceID int, domain string, isBlacklist bool) ([]Campaign, error) {
	data := struct {
		SourceID int    `db:"id"`
		Domain   string `db:"domain"`
	}{
		SourceID: sourceID,
		Domain:   domain,
	}

	var q string
	if isBlacklist {
		q = `
            SELECT DISTINCT 
                c.*
            FROM
                campaigns c
            INNER JOIN
                source_campaign sc ON c.id = sc.campaign_id
            WHERE 
                sc.source_id = :id
                AND c.domain NOT IN (:domain)
        `
	} else {
		q = `
            SELECT DISTINCT 
                c.*
            FROM
                campaigns c
            INNER JOIN
                source_campaign sc ON c.id = sc.campaign_id
            WHERE 
                sc.source_id = :id
                AND c.domain IN (:domain)
        `
	}

	var campaigns []Campaign
	if err := database.NamedQuerySlice(s.log, s.db, q, data, &campaigns); err != nil {
		return nil, err
	}

	return campaigns, nil
}
