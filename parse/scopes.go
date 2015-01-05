package parse

var genTemplate = `/*
  This code was generated by the Doc ORM Generator and isn't meant to be edited.
	If at all possible, please regenerate this file from your gp files instead of
	attempting to edit it to add changes.
*/

package {{ .Name }}

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Conn struct {
	*sql.DB
	AppConfig
	reformat bool
	returning bool
	Log *log.Logger
	{{ range .Tables }}
		{{ .Name }} {{ .Name }}Scope
	{{ end }}
}

func Open(driverName, dataSourceName string) (*Conn, error) {
	c := &Conn{}
	if driverName == "postgres" {
		c.reformat = true
		c.returning = true
	}
	var err error
	c.DB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	c.AppConfig = NewAppConfig(driverName)
	{{ range .Tables }}
	c.{{ .Name }} = New{{ .Name }}Scope(c)
	{{ end }}
	return c, nil
}

func (c *Conn) Clone() *Conn {
	c2 := &Conn{
		DB: c.DB,
		AppConfig: c.AppConfig,
		reformat: c.reformat,
		returning: c.returning,
		Log: c.Log,
	}
	{{ range .Tables }}
	c2.{{ .Name }} = New{{ .Name }}Scope(c2)
	{{ end }}
	return c2
}



{{ range $table := .Tables }}
type {{ .Name }}Scope struct {
	internalScope
}

func New{{ .Name }}Scope(c *Conn) {{ .Name }}Scope {
	return {{ .Name }}Scope {
		internalScope{
			conn:          c,
			table:         c.SQLTable("{{ .Name }}"),
			currentColumn: c.SQLTable("{{ .Name }}") + "." + c.SQLColumn("{{ .Name }}", "{{ .PrimaryKeyColumn.Name }}"),
		},
	}
}

func (scope {{ .Name }}Scope) SetConn(conn *Conn) Scope {
	scope.conn = conn
	return scope
}

func ({{ .Name }}Scope) scopeName() string {
	return "{{ .Name }}"
}

// basic conditions
func (scope {{ .Name }}Scope) Eq(val interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Eq(val)}
}

func (scope {{ .Name }}Scope) Neq(val interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Neq(val)}
}

func (scope {{ .Name }}Scope) Gt(val interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Gt(val)}
}

func (scope {{ .Name }}Scope) Gte(val interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Gte(val)}
}

func (scope {{ .Name }}Scope) Lt(val interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Lt(val)}
}

func (scope {{ .Name }}Scope) Lte(val interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Lte(val)}
}


// multi value conditions
func (scope {{ .Name }}Scope) Between(lower, upper interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Between(lower, upper)}
}

func (scope {{ .Name }}Scope) In(vals ...interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.In(vals...)}
}

func (scope {{ .Name }}Scope) NotIn(vals ...interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.NotIn(vals...)}
}

func (scope {{ .Name }}Scope) Like(str string) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Like(str)}
}

func (scope {{ .Name }}Scope) Where(sql string, vals ...interface{}) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.Where(sql, vals...)}
}

// ordering conditions
func (scope {{ .Name }}Scope) Order(ordering string) {{ .Name }}Scope {
	scope.order = append(scope.order, ordering)
	return scope
}

func (scope {{ .Name }}Scope) Desc() {{ .Name }}Scope {
	scope.order = append(scope.order, scope.currentColumn+" DESC")
	return scope
}

func (scope {{ .Name }}Scope) Asc() {{ .Name }}Scope {
	scope.order = append(scope.order, scope.currentColumn+" ASC")
	return scope
}

// Join funcs
func (scope {{ .Name }}Scope)	OuterJoin(things ...Scope) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.outerJoin("{{ .Name }}", things...)}
}

func (scope {{ .Name }}Scope)	InnerJoin(things ...Scope) {{ .Name }}Scope {
	return {{.Name}}Scope{scope.internalScope.innerJoin("{{ .Name }}", things...)}
}

// JoinBy allows you to specify the exact join SQL statment for one or more
// tables. You can also pass the Scope objects that you are manually joining, 
// which are recorded for future Joining to work off of or to be Include'd.
func (scope {{ .Name }}Scope)	JoinBy(joins string, joinedScopes ...Scope) {{ .Name }}Scope {
	scope.joins = append(scope.joins, joins)
	scope.joinedScopes = append(scope.joinedScopes, joinedScopes...)
	return scope
}

func (scope {{ .Name }}Scope) joinable() string {
	if scope.currentAlias != "" {
		return fmt.Sprintf(
			"%s AS %s",
			scope.conn.SQLTable("{{ .Name }}"),
			scope.currentAlias,
		)
	}
	return scope.conn.SQLTable("{{ .Name }}")
}

func (scope {{ .Name }}Scope) joinTable() string {
	if scope.currentAlias != "" {
		return scope.currentAlias
	}
	return scope.conn.SQLTable("{{ .Name }}")
}

// aggregation filtering
func (scope {{ .Name }}Scope) Having(sql string, vals ...interface{}) {{ .Name }}Scope {
	scope.having = append(scope.having, sql)
	scope.havevals = append(scope.havevals, vals...)
	return scope
}

func (scope {{ .Name }}Scope) GroupBySQL(cols ...string) {{ .Name }}Scope {
	scope.groupBy = append(scope.groupBy, cols...)
	return scope
}

// Result count filtering
func (scope {{ .Name }}Scope) Limit(limit int64) {{ .Name }}Scope {
	scope.limit = &limit
	return scope
}

func (scope {{ .Name }}Scope) Offset(offset int64) {{ .Name }}Scope {
	scope.offset = &offset
	return scope
}

// misc scope operations
func (scope {{ .Name }}Scope) Clear() {{ .Name }}Scope {
	goods := []condition{}
	for _, cond := range scope.conditions {
		if !strings.HasSuffix(cond.column, "."+scope.currentColumn) {
			goods = append(goods, cond)
		}
	}
	scope.conditions = goods
	return scope
}

func (scope {{ .Name }}Scope) ClearAll() {{ .Name }}Scope {
	scope.conditions = []condition{}
	return scope
}

func (scope {{ .Name }}Scope) Base() {{ .Name }}Scope {
	return New{{ .Name }}Scope(scope.conn)
}

// struct saving and loading
func (scope {{ .Name }}Scope) Find(id interface{}) ({{ .Name }}, error) {
	return scope.And(scope.Base().Eq(id)).Retrieve()
}

func (scope {{ .Name }}Scope) Retrieve() ({{ .Name }}, error) {
	val := &{{ .Name }}{}
	m := mapperFor{{ .Name }}(scope.conn, scope.includes)
	m.Current = &val
	scope.columns = m.Columns

	ss, vv := scope.QuerySQL()
	row := scope.conn.QueryRow(ss, vv...)
	err := row.Scan(m.Scanners...)
	if err != nil {
		err = fmt.Errorf("SQL: %s\n%s", ss, err.Error())
	}
	val.cached_conn = scope.conn
	return *val, err

}

func (scope {{ .Name }}Scope) RetrieveAll() ([]{{ .Name }}, error) {
	m := mapperFor{{ .Name }}(scope.conn, scope.includes)
	scope.columns = m.Columns

	ss, vv := scope.QuerySQL()
	rows, err := scope.conn.Query(ss, vv...)
	if err != nil {
		err = fmt.Errorf("SQL: %s\n%s", ss, err.Error())
		return []{{ .Name }}{}, err
	}
	defer rows.Close()

	vals := []{{ .Name }}{}

	for rows.Next() {
		temp := &{{ .Name }}{}
		m.Current = &temp
		err = rows.Scan(m.Scanners...)
		if err != nil {
			return []{{ .Name }}{}, err
		}
		temp.cached_conn = scope.conn
		vals = append(vals, *temp)
	}

	return vals, nil
}

func (scope {{ .Name }}Scope) SaveAll(vals []{{ .Name }}) error {
	for _, val := range vals {
		err := val.Save(scope.conn)
		if err != nil {
			return err
		}
	}
	return nil
}


// Scope attribute updating
func (scope {{ .Name }}Scope) Set(val interface{}) {{ .Name }}Scope {
	if scope.updates == nil {
		scope.updates = make(map[string]interface{})
	}
	colName := strings.TrimPrefix(scope.currentColumn, scope.conn.SQLTable("{{ $table.Name }}")+".")
	scope.updates[colName] = val
	return scope
}

func (scope {{ .Name }}Scope) Update() error {
	sql, vals := scope.UpdateSQL()
	_, err := scope.conn.Exec(sql, vals...)
	return err
}

// subset plucking
func (scope {{ .Name }}Scope) Pick(sql string) {{ .Name }}Scope {
	scope.isDistinct = false
	scope.currentColumn = sql

	return scope
}

func (scope {{ .Name }}Scope) PluckStruct(result interface{}) error {
	return scope.internalScope.pluckStruct("{{ .Name }}", result)
}

// direct sql
func (scope {{ .Name }}Scope) Count() int64 {
	return scope.{{ .PrimaryKeyColumn.Name }}().Distinct().CountOf()
}

func (scope {{ .Name }}Scope) CountBy(sql string) int64 {
	scope.columns = []string{sql}
	ss, sv := scope.QuerySQL()
	var value int64
	row := scope.conn.QueryRow(ss, sv...)
	err := row.Scan(&value)
	if err != nil {
		panic(err)
	}

	return value
}

func (scope {{ .Name }}Scope) CountOf() int64 {
	if scope.isDistinct {
		return scope.CountBy(fmt.Sprintf("COUNT(DISTINCT %s)", scope.currentColumn))
	}
	return scope.CountBy(fmt.Sprintf("COUNT(%s)", scope.currentColumn))
}

func (scope {{ .Name }}Scope) UpdateBySQL(sql string, vals ...interface{}) error {
	scope.columns = []string{""}
	ss, sv := scope.query()
	ss = strings.TrimPrefix(ss, "SELECT FROM "+scope.table)
	ss = fmt.Sprintf("UPDATE %s SET %s %s", scope.table, sql, ss)
	_, err := scope.conn.Exec(ss, append(vals, sv...))
	return err
}

func (scope {{ .Name }}Scope) Delete() error {
	sql, cv := scope.DeleteSQL()
	if sql == "" {
		if err, ok := cv[0].(error); ok {
			return err
		} else {
			return fmt.Errorf("Unspecified Error in DeleteSQL()")
		}
	}
	_, err := scope.conn.Exec(sql, cv...)
	if err != nil {
		return fmt.Errorf("Encountered error: %v\nSQL: %s %v", err, sql, cv)
	}
	return nil
}
func (scope {{ .Name }}Scope) condSQL() (string, []interface{}) {
	conds := []string{}
	vals := []interface{}{}
	for _, condition := range scope.conditions {
		conds = append(conds, condition.ToSQL())
		vals = append(vals, condition.vals...)
	}
	return strings.Join(conds, " AND "), vals
}

// special
func (scope {{ .Name }}Scope) Clone() {{ .Name }}Scope {
	return scope
}

func (scope {{ .Name }}Scope) QuerySQL() (string, []interface{}) {
	return scope.query()
}

func (scope {{ .Name }}Scope) UpdateSQL() (string, []interface{}) {
	sql := fmt.Sprintf(
		"UPDATE %s SET ",
		scope.conn.SQLTable("{{ $table.Name }}"),
	)

	updates := []string{}
	vals := []interface{}{}
	for col, val := range scope.updates {
		updates = append(updates, col + " = ?")
		vals = append(vals, val)
	}
	sql += strings.Join(updates, ", ")

	if len(scope.conditions) > 0 {
		cs, cv := scope.conditionSQL()
		sql += " WHERE " + cs
		vals = append(vals, cv...)
	}
	return sql, vals
}

func (scope {{ .Name }}Scope) DeleteSQL() (string, []interface{}) {
	delScope := scope.Clone()
	if len(scope.joins) > 0 || len(scope.having) > 0 {
		ids, err := scope.{{ .PrimaryKeyColumn.Name }}().Distinct().PluckInt()
		if err != nil {
			return "", []interface{}{err}
		}
		delScope = delScope.ClearAll().{{ .PrimaryKeyColumn.Name }}().In(ids)
	}
	cs, cv := scope.condSQL()

	if cs == "" {
		sql := fmt.Sprintf("DELETE FROM %s",scope.table, cs)
		return sql, []interface{}{}
	} else {
		sql := fmt.Sprintf("DELETE FROM %s WHERE %s",scope.table, cs)
		return sql, cv
	}
}

func (scope {{ .Name }}Scope) As(alias string) {{ .Name }}Scope {
	scope.currentAlias = alias
	return scope
}

func (scope {{ .Name }}Scope) Distinct() {{ .Name }}Scope {
	scope.isDistinct = true
	return scope
}

func (scope {{ .Name }}Scope) And(scopes ...Scope) {{ .Name }}Scope {
	for _, is := range scopes {
		scope.conditions = append(scope.conditions, is.conds()...)
	}
	return scope
}

func (scope {{ .Name }}Scope) Or(scopes ...Scope) {{ .Name }}Scope {
	c := condition{}
	ors := []string{}
	for _, oscope := range scopes {
		cond := []string{}
		conds := oscope.conds()
		if len(conds) == 1 {
			c.vals = append(c.vals, conds[0].vals...)
			ors = append(ors, conds[0].ToSQL())
		} else {
			for _, ocond := range conds {
				c.vals = append(c.vals, ocond.vals...)
				cond = append(cond, ocond.ToSQL())
			}
			ors = append(ors, "("+strings.Join(cond, " AND ")+")")
		}
	}
	c.cond = "(" + strings.Join(ors, " OR ") + ")"

	scope.conditions = append(scope.conditions, c)

	return scope
}

{{ range $column := .Columns }}
	{{ if $column.SimpleType }}
		func (scope {{ $table.Name }}Scope) {{ $column.Name }}() {{ $table.Name }}Scope {
			scope.currentColumn =
				scope.conn.SQLTable("{{ $table.Name }}") +
					"." +
					scope.conn.SQLColumn("{{ $table.Name }}", "{{ $column.Name }}")
			scope.currentAlias = ""
			scope.isDistinct = false
			return scope
		}

		type mapper{{ $table.Name }}To{{ $column.Name }} struct {
			Mapper *mapper{{ $table.Name }}
		}

		{{ if eq $column.GoType "int" }}
			func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
				{{ template "int_mapper" $column }}
			}
		{{ else if eq $column.GoType "string" }}
			func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
				{{ template "string_mapper" $column }}
			}
		{{ else if eq $column.GoType "&{time Time}" }}
			func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
				{{ template "time_mapper" $column }}
			}
		{{ else if eq $column.GoType "bool" }}
			func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
				{{ template "bool_mapper" $column }}
			}
		{{ else if eq $column.GoType "float64" }}
			func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
				{{ template "float64_mapper" $column }}
			}
		{{ else if eq $column.GoType "float32" }}
			func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
				{{ template "float32_mapper" $column }}
			}
		{{ else if eq $column.GoType "[]byte" }}
			func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
				{{ template "byte_mapper" $column }}
			}
		{{ end }}
	{{ end }}
	{{ if $column.Subrecord }}
		/*
			type scope{{ $table.Name }}{{ $column.Subrecord.Name }} struct {
				scope {{ $table.Name }}Scope
			}
			type {{ $table.Name }}{{ $column.Subrecord.Name }}Scope interface {
				Include() {{ $table.Name }}Scope
				{{ range $subcolumn := $column.Subcolumns }}
					{{ $subcolumn.Name }}() {{ $table.Name }}Scope
				{{ end }}
			}

			func (scope scope{{ $table.Name }}) {{ $column.Subrecord.Name }}() {{ $table.Name }}{{ $column.Subrecord.Name }}Scope {
				return scope{{ $table.Name }}{{ $column.Subrecord.Name }}{scope}
			}

			func (scope scope{{ $table.Name }}{{ $column.Subrecord.Name }}) Include() {{ $table.Name }}Scope {
				scope.scope.includes = append(scope.scope.includes, "{{ $column.Subrecord.Name }}")
				return scope.scope
			}
		
			{{ range $subcolumn := $column.Subcolumns }}

				func (scope scope{{ $table.Name }}{{ $column.Subrecord.Name }}) {{ $subcolumn.Name }}() {{ $table.Name }}Scope {
					scope.scope.currentColumn = 	
						scope.scope.conn.SQLTable("{{ $table.Name }}") +
							"." +
							scope.scope.conn.SQLColumn("{{ $table.Name }}", "{{ $subcolumn.Name }}")
					scope.scope.currentAlias = ""
					scope.scope.isDistinct = false
					return scope.scope
				}

				type mapper{{ $table.Name }}To{{ $subcolumn.Name }} struct {
					Mapper *mapper{{ $table.Name }}
				}

				{{ if eq $subcolumn.GoType "int" }}
					func (m mapper{{ $table.Name }}To{{ $subcolumn.Name }}) Scan(v interface{}) error {
						{{ template "int_mapper" $subcolumn }}
					}
				{{ else if eq $subcolumn.GoType "string" }}
					func (m mapper{{ $table.Name }}To{{ $subcolumn.Name }}) Scan(v interface{}) error {
						{{ template "string_mapper" $subcolumn }}
					}
				{{ else if eq $subcolumn.GoType "&{time Time}" }}
					func (m mapper{{ $table.Name }}To{{ $subcolumn.Name }}) Scan(v interface{}) error {
						{{ template "time_mapper" $subcolumn }}
					}
				{{ else if eq $column.GoType "[]byte" }}
					func (m mapper{{ $table.Name }}To{{ $column.Name }}) Scan(v interface{}) error {
						{{ template "byte_mapper" $column }}
					}
				{{ end }}
			{{ end }}
		*/
	{{ end }}
{{ end }}

type mapper{{ .Name }} struct {
	Current  **{{ .Name }}
	Columns  []string
	Scanners []interface{}
}

func mapperFor{{ .Name }}(c *Conn, includes []string) *mapper{{ .Name }} {
	m := &mapper{{ .Name }}{}
	m.Columns = []string{ {{ range $column := .Columns }} {{ if $column.SimpleType }} c.SQLTable("{{ $table.Name }}") + "." + c.SQLColumn("{{ $table.Name }}", "{{ $column.Name }}"), {{ end }} {{ end }} }
	m.Scanners = []interface{}{
		{{ range $column := .Columns }}
			{{ if $column.SimpleType }}
				mapper{{ $table.Name }}To{{ $column.Name }}{m},
			{{ end }}
		{{ end }}
	}

	{{ range $column := .Columns }}
		{{ if $column.Subrecord }}
			if drStringArray(includes).Includes("{{ $column.Subrecord.Name }}") {
				m.Columns = append(m.Columns,
					{{ range $subcolumn := $column.Subcolumns }}c.SQLTable("{{ $table.Name }}") + "." + c.SQLColumn("{{ $table.Name }}", "{{ $subcolumn.Name }}"),{{ end }}
				)

				m.Scanners = append(m.Scanners,
					{{ range $subcolumn := $column.Subcolumns }}mapper{{ $table.Name }}To{{ $subcolumn.Name }}{m},{{ end }}
				)
			}
		{{ end }}
	{{ end }}
	return m
}
{{ end }}

{{ range $table := .Tables }}
	{{ if $table.HasRelationship "ParentHasMany" }}
		{{ range $relate := .Relations }}
			{{ if $relate.IsHasMany }}
				func (t {{ $table.Name }}) {{ $relate.Name }}(c *Conn) ([]{{ $relate.Table }}, error) {
					return t.{{ $relate.Name }}Scope(c).RetrieveAll()
				}
				func (t {{ $table.Name }}) {{ $relate.Name }}Scope(c *Conn) {{ $relate.Table }}Scope {
					return c.{{ $relate.Table }}.{{ $relate.ColumnName }}().Eq(t.{{ $table.PrimaryKeyColumn.Name }})
				}
			{{ end }}
		{{ end }}
	{{ end }}
	{{ if $table.HasRelationship "ChildHasMany" }}
		{{ range $relate := .Relations }}
			{{ if $relate.IsChildHasMany }}
				func (t {{ $table.Name }}) {{ $relate.Name }}(c *Conn) ({{ $relate.Table }}, error) {
					return t.{{ $relate.Name }}Scope(c).Retrieve()
				}
				func (t {{ $table.Name }}) {{ $relate.Name }}Scope(c *Conn) {{ $relate.Table }}Scope {
					return c.{{ $relate.Table }}.Eq(t.{{ $relate.Name }}ID)
				}
			{{ end }}
		{{ end }}
	{{ end }}
	{{ if $table.HasRelationship "HasOne" }}
		{{ range $relate := .Relations }}
			{{ if $relate.IsHasOne }}
				func (t {{ $table.Name }}) {{ $relate.Name }}(c *Conn) ({{ $relate.Table }}, error) {
					return t.{{ $relate.Name }}Scope(c).Retrieve()
				}
				func (t {{ $table.Name }}) {{ $relate.Name }}Scope(c *Conn) {{ $relate.Table }}Scope {
					return c.{{ $relate.Table }}.{{ $relate.ColumnName }}().Eq(t.{{ $table.PrimaryKeyColumn.Name }})
				}
			{{ end }}
		{{ end }}
	{{ end }}
	{{ if $table.HasRelationship "BelongsTo" }}
		{{ range $relate := .Relations }}
			{{ if $relate.IsBelongsTo }}
				func (t {{ $table.Name }}) {{ $relate.Name }}(c *Conn) ({{ $relate.Table }}, error) {
					return t.{{ $relate.Name }}Scope(c).Retrieve()
				}
				func (t {{ $table.Name }}) {{ $relate.Name }}Scope(c *Conn) {{ $relate.Table }}Scope {
					return c.{{ $relate.Table }}.Eq(t.{{ $relate.Name }}ID)
				}
			{{ end }}
		{{ end }}
	{{ end }}
{{ end }}
`
