package input

import (
	"io"
	"reflect"
	"strings"
)

//go:generate go run ../bin/gen-types/gen-types.go

var Schemas = map[string]Schema{
	"AD15AirfoilInfo": AD15AirfoilInfoSchema,
	"AeroDyn14":       AeroDyn14Schema,
	"AeroDyn15":       AeroDyn15Schema,
	"AeroDynBlade":    AeroDynBladeSchema,
	"BeamDyn":         BeamDynSchema,
	"BeamDynBlade":    BeamDynBladeSchema,
	"ElastoDyn":       ElastoDynSchema,
	"ElastoDynBlade":  ElastoDynBladeSchema,
	"ElastoDynTower":  ElastoDynTowerSchema,
	"FAST":            FASTSchema,
	"InflowWind":      InflowWindSchema,
	"ServoDyn":        ServoDynSchema,
	"TurbSim":         TurbSimSchema,
}

type Schema []SchemaEntry

func NewSchema(name string, entries []SchemaEntry) Schema {
	for i := range entries {

		// Get entry
		e := &entries[i]

		// Skip heading entries
		if len(e.Heading) > 0 {
			continue
		}

		// Convert keyword into field
		e.Field = keywordToField(e.Keyword)

		// If entry is a table
		if e.Table != nil {

			// Construct type from struct name and field
			e.Type = name + e.Field

			// Get field from table column keyword
			for j := range e.Table.Columns {
				e.Table.Columns[j].Field = keywordToField(e.Table.Columns[j].Keyword)
			}
		}
	}
	return entries
}

func keywordToField(kw string) string {
	tmp := strings.ReplaceAll(strings.ReplaceAll(kw, "(", ""), ")", "")
	return strings.ToUpper(tmp[:1]) + tmp[1:]
}

type ParseFunc func(s any, v reflect.Value, lines []string) ([]string, error)
type FormatFunc func(w io.Writer, s any, field any, entry SchemaEntry) error

type SchemaEntry struct {
	Keyword      string
	Field        string
	Type         string
	Dims         int
	Desc         string
	Heading      string `json:",omitempty"`
	Default      any    `json:",omitempty"`
	CanBeDefault bool
	Unit         string
	Options      []Option    `json:",omitempty"`
	SkipIf       []Condition `json:",omitempty"`
	Show         []Condition `json:",omitempty"`
	Table        *Table      `json:",omitempty"`
	Parse        ParseFunc   `json:"-"`
	Format       FormatFunc  `json:"-"`
}

type Table struct {
	Columns []TableColumn
}

type TableColumn struct {
	Keyword  string
	Field    string
	Type     string
	Dims     int
	Unit     string
	Desc     string
	Optional bool
}

type Condition struct {
	Field    string `json:"field"`
	Relation string `json:"relation"`
	Value    any    `json:"value"`
}

type Option struct {
	Value any    `json:"value"`
	Text  string `json:"text"`
}

const (
	String = "string"
	Bool   = "bool"
	Int    = "int"
	Float  = "float64"
)
