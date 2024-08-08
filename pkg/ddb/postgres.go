package ddb

import (
	"fmt"

	"go.kyoto.codes/zen/v3/slice"
)

type Postgres struct {
	Connection
}

func (p *Postgres) systemSchemas() []string {
	return []string{"pg_catalog", "information_schema"}
}

func (p *Postgres) QueryTables() ([]Table, error) {
	// Query the database for the tables
	data, err := p.QueryData("SELECT table_name,table_schema FROM information_schema.tables")
	if err != nil {
		return nil, err
	}
	// Convert the data to a slice of Table objects
	tables := slice.Map(data.Rows, func(r []any) Table {
		return Table{
			Name:   r[0].(string),
			Schema: r[1].(string),
		}
	})
	// Mark system tables
	tables = slice.Map(tables, func(t Table) Table {
		if slice.Contains(p.systemSchemas(), t.Schema) {
			t.IsSystem = true
		}
		return t
	})
	// Return
	return tables, nil
}

func (p *Postgres) QueryColumns(table string) ([]Column, error) {
	// Query the database for the columns
	data, err := p.QueryData(fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s'", table))
	if err != nil {
		return nil, err
	}
	// Convert the data to a slice of Column objects
	columns := slice.Map(data.Rows, func(r []any) Column {
		return Column{
			Name: r[0].(string),
			Type: r[1].(string),
		}
	})
	// Return
	return columns, nil
}
