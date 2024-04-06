package dbschema

import (
	_ "embed"
	"fmt"
	"github.com/grumpycatyo-collab/turbo-pancake/business/sys/database"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
)

var (
	//go:embed sql/create_tables/create_campaigns.sql
	createCampaignsDoc string

	//go:embed sql/create_tables/create_sources.sql
	createSourcesDoc string

	//go:embed sql/create_tables/create_source_campaign.sql
	createMidDoc string

	//go:embed sql/drop_tables.sql
	dropTablesDoc string

	//go:embed sql/selects/select_top_5.sql
	top5Doc string

	//go:embed sql/selects/union.sql
	unionDoc string

	//go:embed sql/selects/select_non_linked_campaigns.sql
	nullCountCampaigns string
)

type Source struct {
	Name string `db:"name"`
}
type Campaign struct {
	Name string `db:"name"`
}

type Top5 struct {
	SourceId      int64  `db:"id"`
	SourceName    string `db:"name"`
	CampaignCount int64  `db:"campaign_count"`
}

// Create initial tables
func Create(db *sqlx.DB) error {
	if err := database.StatusCheck(db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in defer:", r)
		}
	}()

	if _, err := tx.Exec(createSourcesDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if _, err := tx.Exec(createCampaignsDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if _, err := tx.Exec(createMidDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// Seed data into sources, campaigns and source_campaign tables
func Seed(db *sqlx.DB) error {
	if err := database.StatusCheck(db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered in defer:", r)
		}
	}()

	rand.Seed(time.Now().UnixNano())

	var sourceStructs []Source
	for i := 0; i < 100; i++ {
		sourceStructs = append(sourceStructs, Source{Name: fmt.Sprintf("Source_%d", rand.Intn(1000))})
	}

	for _, source := range sourceStructs {
		if _, err := tx.Exec(`INSERT INTO sources (name) VALUES (?)`, source.Name); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	var campaignStructs []Campaign
	for i := 0; i < 100; i++ {
		campaignStructs = append(campaignStructs, Campaign{Name: fmt.Sprintf("Campaign_%d", rand.Intn(1000))})
	}

	for _, campaign := range campaignStructs {
		if _, err := tx.Exec(`INSERT INTO campaigns (name) VALUES (?)`, campaign.Name); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	rows, err := tx.Query(`SELECT sources.id AS source_id, campaigns.id AS campaign_id FROM sources JOIN campaigns`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	defer rows.Close()

	sourceCampaignMap := make(map[int][]int)
	for rows.Next() {
		var sourceID, campaignID int
		err := rows.Scan(&sourceID, &campaignID)
		if err != nil {
			return err
		}
		sourceCampaignMap[sourceID] = append(sourceCampaignMap[sourceID], campaignID)
	}

	stmt, err := tx.Prepare(`INSERT INTO source_campaign (source_id, campaign_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	rand.Seed(time.Now().UnixNano())
	for sourceID, campaignIDs := range sourceCampaignMap {
		numCampaigns := rand.Intn(11)
		for i := 0; i < numCampaigns && i < len(campaignIDs); i++ {
			if _, err := stmt.Exec(sourceID, campaignIDs[i]); err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		}
	}

	return tx.Commit()
}

// Show an example how to work with some of SELECT
func Show(db *sqlx.DB) error {
	if err := database.StatusCheck(db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	// Example on selecting top 5 sources
	result := []Top5{}
	err := db.Select(&result, top5Doc)
	if err != nil {
		return err
	}
	id, name, count := result[0], result[1], result[2]
	fmt.Printf("%#v\n%#v\n%#v\n", id, name, count)

	return nil
}

// DropAll table by necessity
func DropAll(db *sqlx.DB) error {
	if err := database.StatusCheck(db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(dropTablesDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
