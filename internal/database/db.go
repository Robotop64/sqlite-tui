package database

import (
	"database/sql"
	"errors"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Robotop64/sqlite-tui/internal/utils"
	persistent "github.com/Robotop64/sqlite-tui/internal/utils/persistent"
)

// "github.com/mattn/go-sqlite3"

type filter struct {
	Columns    []string
	MaxNumRows int
	SortCol    string //which column to sort by
	SortDir    uint8  // 0: None, 1: Ascending, 2: Descending
}

type filterMode int8

const (
	Composite filterMode = iota
	Single
)

type fieldMode int8

const (
	like fieldMode = iota
	exact
	contains
)

type Schema struct {
	TableNames  []string
	CreationSQL []string
}

var ActiveSchema Schema

var dbPath string

var rootPath string
var activeTarget persistent.Target

func SetTarget(profile_root string, target persistent.Target) error {
	rootPath = filepath.Dir(profile_root)
	activeTarget = target

	dbPath = target.DatabasePath
	if dbPath == "" {
		return errors.New("Database path is empty")
	}
	if !filepath.IsAbs(dbPath) {
		dbPath = utils.RelativeToAbsolutePath(rootPath, dbPath)
	}
	return nil
}

func SetSchema() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, sql FROM sqlite_master WHERE type='table'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var schema Schema
	for rows.Next() {
		var name, sql string
		if err := rows.Scan(&name, &sql); err != nil {
			return err
		}
		schema.TableNames = append(schema.TableNames, name)
		schema.CreationSQL = append(schema.CreationSQL, sql)
	}
	ActiveSchema = schema

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
