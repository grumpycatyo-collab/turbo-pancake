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

	//go:embed sql/insert_source_campaign.sql
	insertMidDoc string

	//go:embed sql/selects/select_top_5.sql
	top5Doc string

	//go:embed sql/selects/union.sql
	unionDoc string

	//go:embed sql/selects/select_non_linked_campaigns.sql
	nullCountCampaigns string
)

type Source struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
type Campaign struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// An example how to view information from the selects
type Top5 struct {
	SourceId      int64  `db:"id"`
	SourceName    string `db:"name"`
	CampaignCount int64  `db:"campaign_count"`
}

// Creating initial tables
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

	// NOTE: could've done better but the database/sql driver doesn't allow multiple queries at a time
	_, err = tx.Exec(createSourcesDoc)
	if err != nil {
		return err
	}

	_, err = tx.Exec(createCampaignsDoc)
	if err != nil {
		return err
	}

	_, err = tx.Exec(createMidDoc)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Seeding data into sources, campaigns and source_campaign tables
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
		_, err := tx.Exec(`INSERT INTO sources (name) VALUES (?)`, source.Name)
		if err != nil {
			return err
		}
	}

	var campaignStructs []Campaign
	for i := 0; i < 100; i++ {
		campaignStructs = append(campaignStructs, Campaign{Name: fmt.Sprintf("Campaign_%d", rand.Intn(1000))})
	}

	for _, campaign := range campaignStructs {
		_, err := tx.Exec(`INSERT INTO campaigns (name) VALUES (?)`, campaign.Name)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Showing an example how to work with some of SELECTs
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

// Dropping all table by necessity
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
