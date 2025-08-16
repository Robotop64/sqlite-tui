package persistent

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	utils "SQLite-GUI/internal/utils"
)

type SourceType int

const (
	SRC_Unknown SourceType = iota
	SRC_File_BIN
	SRC_File_JSON
	SRC_Database_SQLite
)

type Source interface {
	Path() string
	SourceType() SourceType
	Load(path string) error
	Save() error
}

var Sources []*Source

func NewSource(path string) (Source, error) {
	if !utils.CheckPath(path) {
		return nil, fmt.Errorf("invalid path: %s", path)
	}
	var source Source
	file_ext := filepath.Ext(path)
	switch file_ext {
	case ".db", ".sqlite", ".sqlite3":
		source = &SQLDatabase{}
	}

	if err := source.Load(path); err != nil {
		return nil, fmt.Errorf("failed to load source: %v", err)
	}

	Sources = append(Sources, &source)
	return source, nil
}

type SQLDatabase struct {
	path    string
	Tables  []string
	Columns map[string][]SQLColumn
}
type SQLColumn struct {
	Name string
	Type string
}

func (db *SQLDatabase) Path() string           { return db.path }
func (db *SQLDatabase) SourceType() SourceType { return SRC_Database_SQLite }
func (db *SQLDatabase) Load(path string) error {
	if !utils.CheckPath(path) {
		return fmt.Errorf("invalid path: %s", path)
	}
	db.path = path
	if conn, err := sql.Open("sqlite3", db.path); err != nil {
		defer conn.Close()
		return fmt.Errorf("failed to open database: %v", err)
	} else {
		defer conn.Close()
		if rows, err := conn.Query("SELECT name FROM sqlite_master WHERE type='table';"); err == nil {
			defer rows.Close()

			for rows.Next() {
				var tableName string
				if err := rows.Scan(&tableName); err != nil {
					return err
				}
				db.Tables = append(db.Tables, tableName)
			}
		} else {
			return fmt.Errorf("failed to query database for tables: %v", err)
		}
		db.Columns = make(map[string][]SQLColumn, len(db.Tables))
		for _, table := range db.Tables {
			if row, err := conn.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 0;", table)); err != nil {
				return fmt.Errorf("failed to query table %s: %v", table, err)
			} else {
				defer row.Close()
				if columns, err := row.ColumnTypes(); err != nil {
					return fmt.Errorf("failed to get columns for table %s: %v", table, err)
				} else {
					for _, col := range columns {
						db.Columns[table] = append(db.Columns[table], SQLColumn{Name: col.Name(), Type: col.DatabaseTypeName()})
					}
				}
			}
		}
	}
	return nil
}
func (db *SQLDatabase) Save() error {
	return nil
}
