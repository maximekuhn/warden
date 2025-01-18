package sqlite

import "testing"

func TestMigrate(t *testing.T) {
	db := createTmpDb()
	defer db.Close()

	if err := Migrate(db, 0); err != nil {
		t.Fatalf("Migrate(): expected ok got err %v", err)
	}
}
