package dbs

import (
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestDBConfig(t *testing.T) {
	if db := FetchDB("configuration"); db != nil {
		rows, err := db.Queryx("select * from configuration_user")
		if err != nil {
			t.Fatal(err)
		}

		defer func(rows *sqlx.Rows) {
			_ = rows.Close()
		}(rows)
		for rows.Next() {
			// users := make(map[string]interface{})
			user, err := rows.SliceScan()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%+v", user)
		}
	}
}
