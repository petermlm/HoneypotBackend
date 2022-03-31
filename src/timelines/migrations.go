package timelines

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func (t *timelines) MigrateCmd(commands []string) error {
	oldVersion, newVersion, err := migrations.Run(t.db, commands...)
	if err != nil {
		return err
	}

	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}

	return nil
}

func (t *timelines) MigrationsTableExsits() bool {
	res, _ := t.db.Exec(`
           SELECT 1
           FROM   information_schema.tables
           WHERE  table_name = 'gopg_migrations';
    `)
	return res.RowsReturned() > 0
}
