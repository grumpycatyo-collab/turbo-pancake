package source

import (
	"errors"
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source/db"
	"github.com/grumpycatyo-collab/turbo-pancake/business/sys/database"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

type Core struct {
	store db.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zerolog.Logger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

func (c Core) QueryCampaignsBySourceID(SourceID int, Domain string, Filter string) ([]Campaign, error) {
	isBlacklist := false
	if Filter == "white" {
		isBlacklist = false
	} else {
		isBlacklist = true
	}
	dbCampaigns, err := c.store.QueryCampaignsBySourceID(SourceID, Domain, isBlacklist)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return toCampaignSlice(dbCampaigns), nil
}
