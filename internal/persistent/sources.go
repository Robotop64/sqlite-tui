package persistent

import (
	"path/filepath"
	"strings"
)

type SourceType int

const (
	SRC_Unknown SourceType = iota
	SRC_File_BIN
	SRC_File_JSON
	SRC_Database_SQLite
)

type Filter_Table struct {
	Columns []string

	NumRows uint
	SortCol int   // which column to sort by, ref to Columns field
	SortDir uint8 // 0 = asc, 1 = desc
}

type Filter_DB struct {
	TableName string
}

var Sources []*Source

type Source struct {
	Path       string
	SourceType SourceType
	Data       any
	Filter     any
}

func LoadSource(path string) Source {
	type_src := getSourceType(path)
	return Source{
		Path:       path,
		SourceType: type_src,
	}
}

func getSourceType(path string) SourceType {
	suffix := strings.ToLower(filepath.Ext(path))
	switch suffix {
	case ".crsty":
		return SRC_File_BIN
	case ".json":
		return SRC_File_JSON
	case ".sqlite":
		return SRC_Database_SQLite
	default:
		return SRC_Unknown
	}
}
