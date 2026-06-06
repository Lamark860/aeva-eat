package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// Интеграционный тест GetManyByIDs/List против реального Postgres. Пропускается,
// если не задан AEVA_TEST_DSN — поэтому обычный `go test ./...` (и CI backend job
// без БД) его не запускает. Локально: поднять чистый PG, прогнать миграции,
// засеять пару строк и запустить с AEVA_TEST_DSN=postgres://...
func TestGetManyByIDs_Integration(t *testing.T) {
	dsn := os.Getenv("AEVA_TEST_DSN")
	if dsn == "" {
		t.Skip("AEVA_TEST_DSN not set — integration test skipped")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
	repo := NewPlaceRepo(db)

	// Сеем: 2 места, отзыв с автором и gem на place A.
	mustExec(t, db, `INSERT INTO users (id, username, password_hash, role) VALUES
		(901,'ituser','x','user') ON CONFLICT (id) DO NOTHING`)
	mustExec(t, db, `INSERT INTO places (id, name, city, created_by) VALUES
		(901,'IT Place A','Тест-Сити',901),(902,'IT Place B','Тест-Сити',901)
		ON CONFLICT (id) DO NOTHING`)
	mustExec(t, db, `INSERT INTO reviews (id, place_id, food_rating, service_rating, vibe_rating, is_gem)
		VALUES (9001,901,8,9,7,true) ON CONFLICT (id) DO NOTHING`)
	mustExec(t, db, `INSERT INTO review_authors (review_id, user_id) VALUES (9001,901)
		ON CONFLICT DO NOTHING`)
	t.Cleanup(func() {
		db.Exec(`DELETE FROM reviews WHERE id=9001`)
		db.Exec(`DELETE FROM places WHERE id IN (901,902)`)
		db.Exec(`DELETE FROM users WHERE id=901`)
	})

	// Порядок входа должен сохраняться: B, затем A.
	got, err := repo.GetManyByIDs([]int{902, 901})
	if err != nil {
		t.Fatalf("GetManyByIDs: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 places, got %d", len(got))
	}
	if got[0].ID != 902 || got[1].ID != 901 {
		t.Fatalf("input order not preserved: got %d,%d", got[0].ID, got[1].ID)
	}

	a := got[1] // place A
	if !a.IsGemPlace {
		t.Error("place A should be gem place")
	}
	if a.GemStatus == nil || len(a.GemStatus.MarkedBy) != 1 || a.GemStatus.MarkedBy[0].ID != 901 {
		t.Errorf("gem_status not batched correctly: %+v", a.GemStatus)
	}
	if len(a.Reviewers) != 1 || a.Reviewers[0].ID != 901 {
		t.Errorf("reviewers not batched correctly: %+v", a.Reviewers)
	}
	if a.AvgFood == nil || *a.AvgFood != 8 {
		t.Errorf("avg_food want 8, got %v", a.AvgFood)
	}

	b := got[0] // place B — без отзывов
	if b.IsGemPlace || b.GemStatus != nil || b.ReviewCount != 0 {
		t.Errorf("place B should have no reviews/gem: %+v", b)
	}

	// Soft-deleted место не возвращается.
	mustExec(t, db, `UPDATE places SET deleted_at = now() WHERE id = 902`)
	got2, err := repo.GetManyByIDs([]int{902, 901})
	if err != nil {
		t.Fatalf("GetManyByIDs after soft-delete: %v", err)
	}
	if len(got2) != 1 || got2[0].ID != 901 {
		t.Fatalf("soft-deleted place B should be excluded, got %d places", len(got2))
	}
}

func mustExec(t *testing.T, db *sql.DB, q string) {
	t.Helper()
	if _, err := db.Exec(q); err != nil {
		t.Fatalf("exec %q: %v", q, err)
	}
}
