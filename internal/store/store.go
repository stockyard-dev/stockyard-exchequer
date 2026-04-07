package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Budget is a tracked spending allocation. Allocated and Spent are stored
// as integer cents to avoid floating-point money bugs (matches the steward
// pattern). Period is one of: monthly, quarterly, yearly, project.
type Budget struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Allocated int    `json:"allocated"` // cents
	Spent     int    `json:"spent"`     // cents
	Period    string `json:"period"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "exchequer.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS budgets(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT DEFAULT '',
		allocated INTEGER DEFAULT 0,
		spent INTEGER DEFAULT 0,
		period TEXT DEFAULT 'monthly',
		start_date TEXT DEFAULT '',
		end_date TEXT DEFAULT '',
		notes TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_budgets_category ON budgets(category)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Budget) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Period == "" {
		e.Period = "monthly"
	}
	_, err := d.db.Exec(
		`INSERT INTO budgets(id, name, category, allocated, spent, period, start_date, end_date, notes, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.Category, e.Allocated, e.Spent, e.Period, e.StartDate, e.EndDate, e.Notes, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Budget {
	var e Budget
	err := d.db.QueryRow(
		`SELECT id, name, category, allocated, spent, period, start_date, end_date, notes, created_at
		 FROM budgets WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Name, &e.Category, &e.Allocated, &e.Spent, &e.Period, &e.StartDate, &e.EndDate, &e.Notes, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Budget {
	rows, _ := d.db.Query(
		`SELECT id, name, category, allocated, spent, period, start_date, end_date, notes, created_at
		 FROM budgets ORDER BY name ASC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Budget
	for rows.Next() {
		var e Budget
		rows.Scan(&e.ID, &e.Name, &e.Category, &e.Allocated, &e.Spent, &e.Period, &e.StartDate, &e.EndDate, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Budget) error {
	_, err := d.db.Exec(
		`UPDATE budgets SET name=?, category=?, allocated=?, spent=?, period=?, start_date=?, end_date=?, notes=?
		 WHERE id=?`,
		e.Name, e.Category, e.Allocated, e.Spent, e.Period, e.StartDate, e.EndDate, e.Notes, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM budgets WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM budgets`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []Budget {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ? OR category LIKE ? OR notes LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s)
	}
	if v, ok := filters["category"]; ok && v != "" {
		where += " AND category=?"
		args = append(args, v)
	}
	if v, ok := filters["period"]; ok && v != "" {
		where += " AND period=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, name, category, allocated, spent, period, start_date, end_date, notes, created_at
		 FROM budgets WHERE `+where+`
		 ORDER BY name ASC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Budget
	for rows.Next() {
		var e Budget
		rows.Scan(&e.ID, &e.Name, &e.Category, &e.Allocated, &e.Spent, &e.Period, &e.StartDate, &e.EndDate, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns aggregate budget metrics for the dashboard. Includes total
// budget count, total allocated and spent (cents), remaining (allocated -
// spent, can be negative), over-budget count, and breakdowns by category.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":           d.Count(),
		"total_allocated": 0,
		"total_spent":     0,
		"remaining":       0,
		"over_budget":     0,
		"by_category":     map[string]int{},
	}

	var allocated, spent int
	d.db.QueryRow(`SELECT COALESCE(SUM(allocated), 0) FROM budgets`).Scan(&allocated)
	d.db.QueryRow(`SELECT COALESCE(SUM(spent), 0) FROM budgets`).Scan(&spent)
	m["total_allocated"] = allocated
	m["total_spent"] = spent
	m["remaining"] = allocated - spent

	var overBudget int
	d.db.QueryRow(`SELECT COUNT(*) FROM budgets WHERE spent > allocated AND allocated > 0`).Scan(&overBudget)
	m["over_budget"] = overBudget

	if rows, _ := d.db.Query(`SELECT category, COUNT(*) FROM budgets WHERE category != '' GROUP BY category`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_category"] = by
	}

	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
