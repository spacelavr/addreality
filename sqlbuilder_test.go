package sqlbuilder_test

import (
	"testing"

	"sqlbuilder"

	"github.com/stretchr/testify/assert"
)

type Row struct {
	Name       string
	GroupID    uint
	PlatformID uint
}
type Rows []Row

var (
	cases = []struct {
		name      string
		driver    int
		maxLine   int
		maxParams int
	}{
		{"pgsql", sqlbuilder.PgSQLDriver, sqlbuilder.PgSQLMaxLine, sqlbuilder.PgSQLMaxParams},
		{"mssql", sqlbuilder.MSSQLDriver, sqlbuilder.MSSQLMaxLine, sqlbuilder.MSSQLMaxParams},
	}
)

func TestNewInsertBuilder(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			builder, err := sqlbuilder.NewInsertBuilder(c.driver)
			assert.NoError(t, err)
			assert.NotNil(t, builder)
			assert.Equal(t, c.maxLine, builder.GetMaxLine())
			assert.Equal(t, c.maxParams, builder.GetMaxParams())
		})
	}
}

func TestBuilder_Append(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rows := Rows{}

			switch c.driver {
			case sqlbuilder.PgSQLDriver:
				// 3 - params count in a row, 1 - for new batch query
				for i := 0; i < c.maxParams/3+1; i++ {
					rows = append(rows, Row{Name: "device", GroupID: 1, PlatformID: 1})
				}
			case sqlbuilder.MSSQLDriver:
				// 1 - for new batch query
				for i := 0; i < c.maxLine+1; i++ {
					rows = append(rows, Row{Name: "device", GroupID: 1, PlatformID: 1})
				}
			default:
				t.Errorf("driver must be set, actual: %v", c.driver)
			}

			var b, err = sqlbuilder.NewInsertBuilder(c.driver)
			assert.NoError(t, err)

			for _, r := range rows {
				err = b.Append(r.Name, r.GroupID, r.PlatformID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuilder_ToSQL(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rows := Rows{}

			for i := 0; i < 3; i++ {
				rows = append(rows, Row{Name: "device", GroupID: 1, PlatformID: 1})
			}

			var b, err = sqlbuilder.NewInsertBuilder(c.driver)
			assert.NoError(t, err)

			for _, r := range rows {
				err = b.Append(r.Name, r.GroupID, r.PlatformID)
				assert.NoError(t, err)
			}

			bq, err := b.ToSQL()
			assert.NoError(t, err)
			if assert.NotNil(t, bq) {
				assert.Equal(t, "($1,$2,$3 ),($4,$5,$6 ),($7,$8,$9 ) ", bq[0].Query)
			}
		})
	}
}
