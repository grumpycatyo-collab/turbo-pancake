package source

import (
	"errors"
	"fmt"
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

func (c Core) QueryCampaignsBySourceID(SourceID int) ([]Campaign, error) {
	dbCampaigns, err := c.store.QueryCampaignsBySourceID(SourceID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toCampaignSlice(dbCampaigns), nil
}
