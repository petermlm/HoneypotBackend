package timelines

import (
	"context"
	"honeypot/settings"
	"log"
	"time"

	"github.com/go-pg/pg"
)

type Timelines interface {
	// General
	Close()

	// Migrations
	MigrateCmd([]string) error
	MigrationsTableExsits() bool

	// Queries
	InsertConnAttemp(*ConnAttemp) error
	GetTotalConsumptions(context.Context, string) (*CountResult, error)
	GetMapData(context.Context, string) ([]*MapDataEntry, error)
	GetConnAttemps(context.Context, string) ([]*ConnAttempSimple, error)
	GetTopConsumers(context.Context, string) ([]*MapDataEntry, error)
	GetTopFlavours(context.Context, string) ([]*PortCount, error)
	GetBytes(context.Context, string, string) ([]*BytesList, error)
	ExportData(context.Context) (string, error)
}

type timelines struct {
	db *pg.DB
}

func InitTimelines() Timelines {
	t := new(timelines)

	// Connect to database, waits one second between attempts
	log.Println("Connecting with Postgres...")
	for i := 0; i < settings.DatabaseConnRetries; i++ {
		t.db = pg.Connect(&pg.Options{
			Addr:     settings.DatabaseAddr,
			Database: settings.DatabaseDatabase,
			User:     settings.DatabaseUser,
			Password: settings.DatabasePassword,
		})

		dbConGood := t.isDbConGood()

		if dbConGood {
			break
		}

		log.Printf("...attemp %d\n", i)
		t.db = nil
		time.Sleep(time.Second)
	}

	if t.db == nil {
		panic("Can't connect with database")
	}

	log.Println("Database connection established")
	return t
}

func (t *timelines) Close() {
	t.db.Close()
	log.Println("Database connection closed")
}

func (t *timelines) isDbConGood() bool {
	var n int
	_, err := t.db.QueryOne(pg.Scan(&n), "SELECT 1")
	return err == nil
}
