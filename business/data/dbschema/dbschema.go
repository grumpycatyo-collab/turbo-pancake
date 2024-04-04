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
	Name string `db:"name"`
}
type Campaign struct {
	Name string `db:"name"`
}

func Create(db *sqlx.DB) error {
	if err := database.StatusCheck(db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// NOTE: could've done better but the database/sql driver doesn't allow multiple queries at a time
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

	return nil
}

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

	_, err = db.NamedExec(`INSERT INTO sources (name) VALUES (:name)`, sourceStructs)
	if err != nil {
		return err
	}

	var campaignStructs []Campaign
	for i := 0; i < 100; i++ {
		campaignStructs = append(campaignStructs, Campaign{Name: fmt.Sprintf("Campaign_%d", rand.Intn(1000))})
	}

	_, err = db.NamedExec(`INSERT INTO campaigns (name) VALUES (:name)`, campaignStructs)
	if err != nil {
		return err
	}

	_, err = db.Exec(insertMidDoc)
	if err != nil {
		return err
	}

	return nil
}

func Show(db *sqlx.DB) error {
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

	result, err := db.Exec(top5Doc)
	if err != nil {
		return err
	}
	fmt.Println(result)

	return nil
}
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
