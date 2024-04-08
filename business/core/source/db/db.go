package db

import (
	"fmt"
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

func (s Store) QueryCampaignsBySourceID(sourceID int) ([]Campaign, error) {
	data := struct {
		SourceID int `db:"id"`
	}{
		SourceID: sourceID,
	}
	// TODO: Add loggers
	const q = `
    SELECT
        c.*
    FROM
        campaigns c
    INNER JOIN
        source_campaign sc ON c.id = sc.campaign_id
    WHERE 
        sc.source_id = :id`

	var campaigns []Campaign
	if err := database.NamedQuerySlice(s.log, s.db, q, data, &campaigns); err != nil {
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	return campaigns, nil
}
