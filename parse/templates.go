package parse

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"code.google.com/p/go.tools/imports"
	"github.com/acsellers/inflections"
)

var genTemplate = `{{ define "int_mapper" }}
	if v == nil {
		// do nothing, use zero value
	} else if s, ok := v.(int64); ok {
		{{ if .MustNull }}
			var temp int
			temp = int(s)
			(*m.Mapper.Current).{{ .Name }} = &temp
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = int(s)
		{{ end }}
	} else if b, ok := v.([]byte); ok {
		i, err := strconv.Atoi(string(b))
		{{ if .MustNull }}
			(*m.Mapper.Current).ID = &i
		{{ else }}
			(*m.Mapper.Current).ID = i
		{{ end }}
		return err
	}
	return nil
{{ end }}
{{ define "string_mapper" }}
	if v == nil {
		// do nothing, use zero value
	} else if s, ok := v.(string); ok {
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &s
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = s
		{{ end }}
	} else if s, ok := v.([]byte); ok {
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &string(s)
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = string(s)
		{{ end }}
	}

	return nil
{{ end }}
{{ define "time_mapper" }}
	if v == nil {
		// do nothing, use zero value
	} else if s, ok := v.(time.Time); ok {
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &s
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = s
		{{ end }}
	}

	return nil
{{ end }}
{{ define "bool_mapper" }}
	if v == nil {
		// it is false or null, the zero values
	} else if b, ok := v.(bool); ok {
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &b
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = b
		{{ end }}
	} else if i, ok := v.(int64); ok {
		{{ if .MustNull }}
			var tb bool
			if i == 0 {
				(*m.Mapper.Current).{{ .Name }} = &tb
			} else {
				tb = true
				(*m.Mapper.Current).{{ .Name }} = &tb
			}
		{{ else }}
			if i != 0 {
				(*m.Mapper.Current).{{ .Name }} = true
			}
		{{ end }}
	}

	return nil
{{ end }}
{{ define "float64_mapper" }}
	if v == nil {
		// it is false or null, the zero values
	} else if f, ok := v.(float64); ok {
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &f
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = f
		{{ end }}
	} else if i, ok := v.(int64); ok {
		tf := float64(i)
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &tf
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = tf
		{{ end }}
	}	else {
		return fmt.Errorf("Value not recognized as float64, received %v", v)
	}

	return nil
{{ end }}

{{ define "float32_mapper" }}
	if v == nil {
		// it is false or null, the zero values
	} else if f, ok := v.(float64); ok {
		mf := float32(f)
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &mf
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = mf
		{{ end }}
	} else if i, ok := v.(int64); ok {
		tf := float32(i)
		{{ if .MustNull }}
			(*m.Mapper.Current).{{ .Name }} = &tf
		{{ else }}
			(*m.Mapper.Current).{{ .Name }} = tf
		{{ end }}
	}	else {
		return fmt.Errorf("Value not recognized as float32, received %v", v)
	}

	return nil
{{ end }}
/*
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

type Scope interface {
	condSQL() (string, []interface{})
	ToSQL() (string, []interface{})
	scopeName() string
	Conn() *Conn
	SetConn(*Conn) Scope
	joinOn(Scope) (string, bool)
	joinable() string
	joinTable() string
}

var (
	{{ range $table := .Tables }}
		{{ if public $table.Name }}
			{{ if eq (plural $table.Name) $table.Name }}
				{{ $table.Name }}s {{ .Name }}Scope = &scope{{ .Name }}{}
			{{ else }}
				{{ plural $table.Name }} {{ .Name }}Scope = &scope{{ .Name }}{}
			{{ end }}
		{{ end }}
	{{ end }}
)

{{ range $table := .Tables }}
type {{ .Name }}Scope interface {
	// column scopes
	{{ range $column := .Columns }}{{ if $column.SimpleType }}{{ $column.Name }}() {{ $table.Name }}Scope{{ end }}
{{ end }}
	{{ range $column := .Columns }}{{ if $column.Subrecord }}{{ $column.Subrecord.Name }}() {{ $table.Name }}{{ $column.Subrecord.Name }}Scope{{ end }}
{{ end }}

	// Basic conditions
	Eq(val interface{}) {{ .Name }}Scope
	Neq(val interface{}) {{ .Name }}Scope
	Gt(val interface{}) {{ .Name }}Scope
	Gte(val interface{}) {{ .Name }}Scope
	Lt(val interface{}) {{ .Name }}Scope
	Lte(val interface{}) {{ .Name }}Scope
	Like(string) {{ .Name }}Scope

	// multi value conditions
	Between(lower, upper interface{}) {{ .Name }}Scope
	In(vals ...interface{}) {{ .Name }}Scope
	NotIn(vals ...interface{}) {{ .Name }}Scope
	Where(sql string, vals ...interface{}) {{ .Name }}Scope 
	// ordering conditions
	Order(ordering string) {{ .Name }}Scope
	Desc() {{ .Name }}Scope
	Asc() {{ .Name }}Scope

	// Join funcs
	OuterJoin(things ...Scope) {{ .Name }}Scope
	InnerJoin(things ...Scope) {{ .Name }}Scope
	JoinBy(joins string, joined ...Scope) {{ .Name }}Scope

	// Aggregation filtering
	Having(sql string, vals ...interface{}) {{ .Name }}Scope
	GroupBySQL(cols ...string) {{ .Name }}Scope

	// Result count filtering
	Limit(limit int64) {{ .Name }}Scope
	Offset(offset int64) {{ .Name }}Scope

	// Misc. Scope operations
	Clear() {{ .Name }}Scope 
	ClearAll() {{ .Name }}Scope
	Base() {{ .Name }}Scope

	// Struct instance saving and loading
	Find(id interface{}) ({{ .Name }}, error)
	Retrieve() ({{ .Name }}, error)
	RetrieveAll() ([]{{ .Name }}, error)
	SaveAll(vals []{{ .Name }}) error

	// Scope attribute updating
	Set(val interface{}) {{ .Name }}Scope
	Update() error

	// Subset plucking
	Pick(sql string) {{ .Name }}Scope
	PluckString() ([]string, error)
	PluckInt() ([]int64, error)
	PluckTime() ([]time.Time, error)
	PluckStruct(result interface{}) error

	// Direct SQL operations
	Count() int64
	CountBy(sql string) int64
	CountOf() int64
	UpdateSQL(sql string, vals ...interface{}) error
	Delete() error

	// Special operations
	Scope

	// special operations
	Clone() {{ .Name }}Scope
	As(alias string) {{ .Name }}Scope
	Distinct() {{ .Name }}Scope
	And(...{{ .Name }}Scope) {{ .Name }}Scope
	Or(...{{ .Name }}Scope) {{ .Name }}Scope
}
{{ end }}

type Conn struct {
	*sql.DB
	AppConfig
	reformat bool
	returning bool
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
	{{ range .Tables }}
	c.{{ .Name }} = new{{ .Name }}Scope(c)
	{{ end }}
	return c, nil
}

func (c *Conn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.DB.Exec(c.FormatQuery(query), args...)
}

func (c *Conn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.DB.Query(c.FormatQuery(query), args...)
}

func (c *Conn) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.DB.QueryRow(c.FormatQuery(query), args...)
}

func (c *Conn) FormatQuery(query string) string {
	if !c.reformat {
		return query
	}

	parts := strings.Split(query, "?")
	var newQuery []string
	for i, part := range parts[:len(parts)-1] {
		newQuery = append(newQuery, fmt.Sprintf("%s$%d", part, i+1))
	}
	newQuery = append(newQuery, parts[len(parts)-1])

	return strings.Join(newQuery, "")
}

func (c *Conn) Close() error {
	return c.DB.Close()
}

type AppConfig struct {
	defaultConfig
}

type defaultConfig struct {
	SpecialTables map[string]string
	SpecialColumns map[string]map[string]string
}

func (c *defaultConfig) SQLTable(table string) string {
	if c == nil || c.SpecialTables[table] == "" {
		return strings.ToLower(table)
	}
	return c.SpecialTables[table]
}

func (c *defaultConfig) SQLColumn(table, column string) string {
	if c == nil || c.SpecialColumns[table] == nil || c.SpecialColumns[table][column] == "" {
		return strings.ToLower(column)
	}
	return c.SpecialColumns[table][column]
}

{{ range $table := .Tables }}
type scope{{ .Name }} struct {
	conn                        *Conn
	table                       string
	columns                     []string
	order                       []string
	joins                       []string
	joinedScopes                []Scope
	includes                    []string
	conditions                  []condition
	having                      []string
	havevals                    []interface{}
	groupBy                     []string
	currentColumn, currentAlias string
	isDistinct                  bool
	limit, offset               *int64
	updates                     map[string]interface{}
}

func new{{ .Name }}Scope(c *Conn) *scope{{ .Name }} {
	return &scope{{ .Name }}{
		conn:          c,
		table:         c.SQLTable("{{ .Name }}"),
		currentColumn: c.SQLTable("{{ .Name }}") + "." + c.SQLColumn("{{ .Name }}", "{{ .PrimaryKeyColumn.Name }}"),
	}
}

func (scope scope{{ .Name }}) Conn() *Conn {
	return scope.conn
}
func (scope scope{{ .Name }}) SetConn(conn *Conn) Scope {
	scope.conn = conn
	return scope
}
func (scope{{ .Name }}) scopeName() string {
	return "{{ .Name }}"
}

func (s *scope{{ .Name }}) query() (string, []interface{}) {
	// SELECT (columns) FROM (table) (joins) WHERE (conditions)
	// GROUP BY (grouping) HAVING (havings)
	// ORDER BY (orderings) LIMIT (limit) OFFSET (offset)
	sql := []string{}
	vals := []interface{}{}
	if len(s.columns) == 0 {
		sql = append(sql, "SELECT", s.table+".*")
	} else {
		sql = append(sql, "SELECT", strings.Join(s.columns, ", "))
	}
	// if s.source == nil { // subquery
	//
	// } else {
	sql = append(sql, "FROM", s.table)
	// }
	sql = append(sql, s.joins...)

	if len(s.conditions) > 0 {
		cs, cv := s.conditionSQL()
		sql = append(sql, "WHERE", cs)
		vals = append(vals, cv...)
	}

	// if len(s.groupings) > 0 {
	//   sql = append(sql , "GROUP BY")
	//   for _, grouping := range s.groupings {
	//     sql = append(sql, grouping.ToSQL()
	//   }
	// }

	if len(s.having) > 0 {
		sql = append(sql, "HAVING")
		sql = append(sql, s.having...)
		vals = append(vals, s.havevals...)
	}

	if len(s.order) > 0 {
		sql = append(sql, "ORDER BY")
		sql = append(sql, s.order...)
	}

	if s.limit != nil {
		sql = append(sql, "LIMIT", fmt.Sprintf("%v", *s.limit))
	}

	if s.offset != nil {
		sql = append(sql, "OFFSET", fmt.Sprintf("%v", *s.offset))
	}

	return strings.Join(sql, " "), vals
}

func (scope scope{{ .Name }}) conditionSQL() (string, []interface{}) {
	var vals []interface{}
	conds := []string{}
	for _, condition := range scope.conditions {
		conds = append(conds, condition.ToSQL())
		vals = append(vals, condition.vals...)
	}
	return strings.Join(conds, " AND "), vals
}

// basic conditions
func (scope scope{{ .Name }}) Eq(val interface{}) {{ .Name }}Scope {
	c := condition{column: scope.currentColumn}
	if val == nil {
		c.cond = "IS NULL"
	} else {
		c.cond = "= ?"
		c.vals = []interface{}{val}
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) Neq(val interface{}) {{ .Name }}Scope {
	c := condition{column: scope.currentColumn}
	if val == nil {
		c.cond = "IS NOT NULL"
	} else {
		c.cond = "<> ?"
		c.vals = []interface{}{val}
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) Gt(val interface{}) {{ .Name }}Scope {
	c := condition{
		column: scope.currentColumn,
		cond:   "> ?",
		vals:   []interface{}{val},
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) Gte(val interface{}) {{ .Name }}Scope {
	c := condition{
		column: scope.currentColumn,
		cond:   ">= ?",
		vals:   []interface{}{val},
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) Lt(val interface{}) {{ .Name }}Scope {
	c := condition{
		column: scope.currentColumn,
		cond:   "< ?",
		vals:   []interface{}{val},
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) Lte(val interface{}) {{ .Name }}Scope {

	c := condition{
		column: scope.currentColumn,
		cond:   "<= ?",
		vals:   []interface{}{val},
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

// multi value conditions
func (scope scope{{ .Name }}) Between(lower, upper interface{}) {{ .Name }}Scope {
	c := condition{
		column: scope.currentColumn,
		cond:   "BETWEEN ? AND ?",
		vals:   []interface{}{lower, upper},
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) In(vals ...interface{}) {{ .Name }}Scope {
	if len(vals) == 0 {
		if reflect.TypeOf(vals[0]).Kind() == reflect.Slice {
			rv := reflect.ValueOf(vals[0])
			vals = make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				vals[i] = rv.Index(i).Interface()
			}
		}
	}

	vc := make([]string, len(vals))
	c := condition{
		column: scope.currentColumn,
		cond:   "IN (" + strings.Join(vc, "?, ") + "?)",
		vals:   vals,
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) NotIn(vals ...interface{}) {{ .Name }}Scope {
	vc := make([]string, len(vals))
	c := condition{
		column: scope.currentColumn,
		cond:   fmt.Sprintf("NOT IN (%!s(MISSING)?)", strings.Join(vc, "?, ")),
		vals:   vals,
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) Like(str string) {{ .Name }}Scope {
	c := condition{
		column: scope.currentColumn,
		cond:   "LIKE ?",
		vals:   []interface{}{str},
	}

	scope.conditions = append(scope.conditions, c)
	return scope
}

func (scope scope{{ .Name }}) Where(sql string, vals ...interface{}) {{ .Name }}Scope {
	c := condition{
		cond: sql,
		vals: vals,
	}
	scope.conditions = append(scope.conditions, c)
	return scope
}

// ordering conditions
func (scope scope{{ .Name }}) Order(ordering string) {{ .Name }}Scope {
	scope.order = append(scope.order, ordering)
	return scope
}

func (scope scope{{ .Name }}) Desc() {{ .Name }}Scope {
	scope.order = append(scope.order, scope.currentColumn+" DESC")
	return scope
}

func (scope scope{{ .Name }}) Asc() {{ .Name }}Scope {
	scope.order = append(scope.order, scope.currentColumn+" ASC")
	return scope
}

// Join funcs
func (scope scope{{ .Name }})	OuterJoin(things ...Scope) {{ .Name }}Scope {
	for _, thing := range things {
		thing = thing.SetConn(scope.conn)
		if joinString, ok := scope.joinOn(thing); ok {
			scope.joins = append(scope.joins, fmt.Sprintf(
				"LEFT JOIN %s ON %s",
				thing.joinable(),
				joinString, 
			))
			scope.joinedScopes = append(scope.joinedScopes, thing)
			continue
		} else {
			for _, joinscope := range scope.joinedScopes {
				if joinString, ok := joinscope.joinOn(thing); ok {
					scope.joins = append(scope.joins, fmt.Sprintf(
						"LEFT JOIN %s ON %s",
						thing.joinable(),
						joinString, 
					))
					scope.joinedScopes = append(scope.joinedScopes, thing)
					continue		
				}
			}
		}
		// error
	}
	return scope
}
func (scope scope{{ .Name }})	InnerJoin(things ...Scope) {{ .Name }}Scope {
	for _, thing := range things {
		thing = thing.SetConn(scope.conn)
		if joinString, ok := scope.joinOn(thing); ok {
			scope.joins = append(scope.joins, fmt.Sprintf(
				"INNER JOIN %s ON %s",
				thing.joinable(),
				joinString, 
			))
			scope.joinedScopes = append(scope.joinedScopes, thing)
			continue
		} else {
			for _, joinscope := range scope.joinedScopes {
				if joinString, ok := joinscope.joinOn(thing); ok {
					scope.joins = append(scope.joins, fmt.Sprintf(
						"INNER JOIN %s ON %s",
						thing.joinable(),
						joinString, 
					))
					scope.joinedScopes = append(scope.joinedScopes, thing)
					continue		
				}
			}
		}
		// error
	}

	return scope
}

// JoinBy allows you to specify the exact join SQL statment for one or more
// tables. You can also pass the Scope objects that you are manually joining, 
// which are recorded for future Joining to work off of or to be Include'd.
func (scope scope{{ .Name }})	JoinBy(joins string, joinedScopes ...Scope) {{ .Name }}Scope {
	scope.joins = append(scope.joins, joins)
	return scope
}

func (scope scope{{ .Name }}) joinOn(joinee Scope) (string, bool) {
	ts := Schema.Tables["{{ .Name }}"]
	for _, hm := range ts.HasMany {
		if (hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName()) || hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName() {
			pkc := hm.Parent.PrimaryKeyColumn()
			return fmt.Sprintf(
				"%s.%s = %s.%s",
				scope.conn.SQLTable(hm.Parent.Name),
				scope.conn.SQLColumn(hm.Parent.Name, pkc.Name),
				scope.conn.SQLTable(hm.Child.Name),
				scope.conn.SQLColumn(hm.Child.Name, hm.ChildColumn.Name),
			), true
		}
	}
	for _, hm := range ts.ChildOf {
		if (hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName()) || hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName() {
			pkc := hm.Parent.PrimaryKeyColumn()
			return fmt.Sprintf(
				"%s.%s = %s.%s",
				scope.conn.SQLTable(hm.Parent.Name),
				scope.conn.SQLColumn(hm.Parent.Name, pkc.Name),
				scope.conn.SQLTable(hm.Child.Name),
				scope.conn.SQLColumn(hm.Child.Name, hm.ChildColumn.Name),
			), true
		}
	}
	for _, hm := range ts.HasOne {
		if (hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName()) || hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName() {
			pkc := hm.Parent.PrimaryKeyColumn()
			return fmt.Sprintf(
				"%s.%s = %s.%s",
				scope.conn.SQLTable(hm.Parent.Name),
				scope.conn.SQLColumn(hm.Parent.Name, pkc.Name),
				scope.conn.SQLTable(hm.Child.Name),
				scope.conn.SQLColumn(hm.Child.Name, hm.ChildColumn.Name),
			), true
		}
	}
	for _, hm := range ts.BelongsTo {
		if (hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName()) || hm.Parent.Name == scope.scopeName() && hm.Child.Name == joinee.scopeName() {
			pkc := hm.Parent.PrimaryKeyColumn()
			return fmt.Sprintf(
				"%s.%s = %s.%s",
				scope.conn.SQLTable(hm.Parent.Name),
				scope.conn.SQLColumn(hm.Parent.Name, pkc.Name),
				scope.conn.SQLTable(hm.Child.Name),
				scope.conn.SQLColumn(hm.Child.Name, hm.ChildColumn.Name),
			), true
		}
	}
	return "", false
}

func (scope scope{{ .Name }}) joinable() string {
	if scope.currentAlias != "" {
		return fmt.Sprintf(
			"%s AS %s",
			scope.conn.SQLTable("{{ .Name }}"),
			scope.currentAlias,
		)
	}
	return scope.conn.SQLTable("{{ .Name }}")
}

func (scope scope{{ .Name }}) joinTable() string {
	if scope.currentAlias != "" {
		return scope.currentAlias
	}
	return scope.conn.SQLTable("{{ .Name }}")
}

// aggregation filtering
func (scope scope{{ .Name }}) Having(sql string, vals ...interface{}) {{ .Name }}Scope {
	scope.having = append(scope.having, sql)
	scope.havevals = append(scope.havevals, vals...)
	return scope
}

func (scope scope{{ .Name }}) GroupBySQL(cols ...string) {{ .Name }}Scope {
	scope.groupBy = append(scope.groupBy, cols...)
	return scope
}

// Result count filtering
func (scope scope{{ .Name }}) Limit(limit int64) {{ .Name }}Scope {
	scope.limit = &limit
	return scope
}

func (scope scope{{ .Name }}) Offset(offset int64) {{ .Name }}Scope {
	scope.offset = &offset
	return scope
}

// misc scope operations
func (scope scope{{ .Name }}) Clear() {{ .Name }}Scope {
	goods := []condition{}
	for _, cond := range scope.conditions {
		if !strings.HasSuffix(cond.column, "."+scope.currentColumn) {
			goods = append(goods, cond)
		}
	}
	scope.conditions = goods
	return scope
}

func (scope scope{{ .Name }}) ClearAll() {{ .Name }}Scope {
	scope.conditions = []condition{}
	return scope
}

func (scope scope{{ .Name }}) Base() {{ .Name }}Scope {
	return new{{ .Name }}Scope(scope.conn)
}

// struct saving and loading
func (scope scope{{ .Name }}) Find(id interface{}) ({{ .Name }}, error) {
	return scope.And(scope.Base().Eq(id)).Retrieve()
}

func (scope scope{{ .Name }}) Retrieve() ({{ .Name }}, error) {
	val := &{{ .Name }}{}
	m := mapperFor{{ .Name }}(scope.conn, scope.includes)
	m.Current = &val
	scope.columns = m.Columns

	ss, vv := scope.ToSQL()
	row := scope.conn.QueryRow(ss, vv...)
	err := row.Scan(m.Scanners...)
	if err != nil {
		err = fmt.Errorf("SQL: %s\n%s", ss, err.Error())
	}
	return *val, err

}

func (scope scope{{ .Name }}) RetrieveAll() ([]{{ .Name }}, error) {
	m := mapperFor{{ .Name }}(scope.conn, scope.includes)
	scope.columns = m.Columns

	ss, vv := scope.ToSQL()
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
		vals = append(vals, *temp)
	}

	return vals, nil
}

func (scope scope{{ .Name }}) SaveAll(vals []{{ .Name }}) error {
	for _, val := range vals {
		err := val.Save(scope.conn)
		if err != nil {
			return err
		}
	}
	return nil
}


// Scope attribute updating
func (scope scope{{ .Name }}) Set(val interface{}) {{ .Name }}Scope {
	if scope.updates == nil {
		scope.updates = make(map[string]interface{})
	}
	colName := strings.TrimPrefix(scope.currentColumn, scope.conn.SQLTable("{{ $table.Name }}")+".")
	scope.updates[colName] = val
	return scope
}

func (scope scope{{ .Name }}) Update() error {
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

	_, err := scope.conn.Exec(sql, vals...)
	return err
}

// subset plucking
func (scope scope{{ .Name }}) Pick(sql string) {{ .Name }}Scope {
	scope.isDistinct = false
	scope.currentColumn = sql

	return scope
}

func (scope scope{{ .Name }}) PluckString() ([]string, error) {
	if scope.isDistinct{
		scope.currentColumn = "DISTINCT " + scope.currentColumn
	}
	scope.columns = []string{scope.currentColumn}
	ss, vv := scope.ToSQL()
	rows, err := scope.conn.Query(ss, vv...)
	if err != nil {
		return []string{}, err
	}
	vals := []string{}
	defer rows.Close()
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			return []string{}, err
		}
		vals = append(vals, temp)
	}

	return vals, nil
}

func (scope scope{{ .Name }}) PluckInt() ([]int64, error) {
	if scope.isDistinct{
		scope.currentColumn = "DISTINCT " + scope.currentColumn
	}

	scope.columns = []string{scope.currentColumn}
	ss, vv := scope.ToSQL()
	rows, err := scope.conn.Query(ss, vv...)
	if err != nil {
		return []int64{}, err
	}
	vals := []int64{}
	defer rows.Close()
	for rows.Next() {
		var temp int64
		err = rows.Scan(&temp)
		if err != nil {
			return []int64{}, err
		}
		vals = append(vals, temp)
	}

	return vals, nil
}

func (scope scope{{ .Name }}) PluckTime() ([]time.Time, error) {
	if scope.isDistinct{
		scope.currentColumn = "DISTINCT " + scope.currentColumn
	}

	scope.columns = []string{scope.currentColumn}
	ss, vv := scope.ToSQL()
	rows, err := scope.conn.Query(ss, vv...)
	if err != nil {
		return []time.Time{}, err
	}
	vals := []time.Time{}
	defer rows.Close()
	for rows.Next() {
		var temp time.Time
		err = rows.Scan(&temp)
		if err != nil {
			return []time.Time{}, err
		}
		vals = append(vals, temp)
	}

	return vals, nil
}

func (scope scope{{ .Name }}) PluckStruct(result interface{}) error {
	panic("UNIMPLEMENTED")
	return nil
}

// direct sql
func (scope scope{{ .Name }}) Count() int64 {
	return scope.{{ .PrimaryKeyColumn.Name }}().Distinct().CountOf()
}

func (scope scope{{ .Name }}) CountBy(sql string) int64 {
	scope.columns = []string{sql}
	ss, sv := scope.ToSQL()
	var value int64
	row := scope.conn.QueryRow(ss, sv...)
	err := row.Scan(&value)
	if err != nil {
		panic(err)
	}

	return value
}

func (scope scope{{ .Name }}) CountOf() int64 {
	if scope.isDistinct {
		return scope.CountBy(fmt.Sprintf("COUNT(DISTINCT %s)", scope.currentColumn))
	}
	return scope.CountBy(fmt.Sprintf("COUNT(%s)", scope.currentColumn))
}

func (scope scope{{ .Name }}) UpdateSQL(sql string, vals ...interface{}) error {
	scope.columns = []string{""}
	ss, sv := scope.query()
	ss = strings.TrimPrefix(ss, "SELECT FROM "+scope.table)
	ss = fmt.Sprintf("UPDATE %s SET %s %s", scope.table, sql, ss)
	_, err := scope.conn.Exec(ss, append(vals, sv...))
	return err
}

func (scope scope{{ .Name }}) Delete() error {
	delScope := scope.Clone()
	if len(scope.joins) > 0 || len(scope.having) > 0 {
		ids, err := scope.{{ .PrimaryKeyColumn.Name }}().Distinct().PluckInt()
		if err != nil {
			return err
		}
		delScope = delScope.ClearAll().{{ .PrimaryKeyColumn.Name }}().In(ids)
	}
	cs, cv := scope.condSQL()
	if cs == "" {
		sql := fmt.Sprintf("DELETE FROM %s",scope.table, cs)
		_, err := scope.conn.Exec(sql, cv)
		return err
	} else {
		sql := fmt.Sprintf("DELETE FROM %s WHERE %s",scope.table, cs)
		_, err := scope.conn.Exec(sql, cv...)
		if err != nil {
			return fmt.Errorf("Encountered error: %v\nSQL: %s %v", err, sql, cv)
		}
		return nil
	}
}
func (scope scope{{ .Name }}) condSQL() (string, []interface{}) {
	conds := []string{}
	vals := []interface{}{}
	for _, condition := range scope.conditions {
		conds = append(conds, condition.ToSQL())
		vals = append(vals, condition.vals...)
	}
	return strings.Join(conds, " AND "), vals
}

// special
func (scope scope{{ .Name }}) Clone() {{ .Name }}Scope {
	return scope
}

func (scope scope{{ .Name }}) ToSQL() (string, []interface{}) {
	return scope.query()
}

func (scope scope{{ .Name }}) As(alias string) {{ .Name }}Scope {
	scope.currentAlias = alias
	return scope
}

func (scope scope{{ .Name }}) Distinct() {{ .Name }}Scope {
	scope.isDistinct = true
	return scope
}

func (scope scope{{ .Name }}) And(scopes ...{{ .Name }}Scope) {{ .Name }}Scope {
	for _, is := range scopes {
		scope.conditions = append(scope.conditions, is.(scope{{ .Name }}).conditions...)
	}
	return scope
}

func (scope scope{{ .Name }}) Or(scopes ...{{ .Name }}Scope) {{ .Name }}Scope {
	c := condition{}
	ors := []string{}
	for _, oscope := range scopes {
		cond := []string{}
		ascope := oscope.(scope{{ .Name }})
		if len(ascope.conditions) == 1 {
			c.vals = append(c.vals, ascope.conditions[0].vals...)
			ors = append(ors, ascope.conditions[0].ToSQL())
		} else {
			for _, ocond := range ascope.conditions {
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
		func (scope scope{{ $table.Name }}) {{ $column.Name }}() {{ $table.Name }}Scope {
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
		{{ end }}
	{{ end }}
	{{ if $column.Subrecord }}
		type scope{{ $table.Name }}{{ $column.Subrecord.Name }} struct {
			scope scope{{ $table.Name }}
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
			{{ end }}
		{{ end }}
	{{ end }}
{{ end }}

type mapper{{ .Name }} struct {
	Current  **{{ .Name }}
	Columns  []string
	Scanners []interface{}
}

func mapperFor{{ .Name }}(c *Conn, includes []string) *mapper{{ .Name }} {
	m := &mapper{{ .Name }}{}
	m.Columns = []string{
		{{ range $column := .Columns }}
			{{ if $column.SimpleType }}
				c.SQLTable("{{ $table.Name }}") + "." + c.SQLColumn("{{ $table.Name }}", "{{ $column.Name }}"),
			{{ end }}
		{{ end }}
	}
	m.Scanners = []interface{}{
		{{ range $column := .Columns }}
			{{ if $column.SimpleType }}
				mapper{{ $table.Name }}To{{ $column.Name }}{m},
			{{ end }}
		{{ end }}
	}

	{{ range $column := .Columns }}
		{{ if $column.Subrecord }}
			if StringArray(includes).Includes("{{ $column.Subrecord.Name }}") {
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

type StringArray []string

func (sa StringArray) Includes(s string) bool{
	for _, si := range sa {
		if si == s {
			return true
		}
	}
	return false
}

type condition struct {
	column string
	cond   string
	vals   []interface{}
}

func (c condition) ToSQL() string {
	if c.column == "" {
		return c.cond
	}
	return c.column + " " + c.cond
}

func questions(n int) string {
	chars := make([]byte, n*2-1)
	for i, _ := range chars {
		if i % 2 == 0 {
			chars[i] = '?'
		} else {
			chars[i] = ','
		}
	}
	return string(chars)
}
`

var tmpl *template.Template

func init() {
	rg := regexp.MustCompile(`^[A-Z].*`)
	var err error
	tmpl, err = template.New("gen").
		Funcs(template.FuncMap{
		"plural": inflections.Pluralize,
		"public": func(s string) bool {
			return rg.MatchString(s)
		},
	}).
		Parse(genTemplate)
	if err != nil {
		panic(err)
	}
}

func (pkg *Package) OutputTemplates() {
	b := &bytes.Buffer{}
	err := tmpl.Execute(b, pkg)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(pkg.ActiveFiles[0].AST.Name.Name + "_gen.go")
	if err != nil {
		fmt.Println("Could not write schema file")
	}
	defer f.Close()

	ib, err := imports.Process(pkg.ActiveFiles[0].AST.Name.Name+"_gen.go", b.Bytes(), nil)
	if err != nil {
		fmt.Println("Error in Gen File:", err)
		f.Write(b.Bytes())
		return
	}
	f.Write(ib)
}
