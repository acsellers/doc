package parse

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"code.google.com/p/go.tools/imports"
)

var schemaTemplate = `/*
  This code was generated by the Doctor ORM Generator and isn't meant to be edited.
	If at all possible, please regenerate this file from your gp files instead of
	attempting to edit it to add changes.
*/

package {{ .Name }}

import "github.com/acsellers/dr/schema"

func DefaultInt(col string) schema.Column {
	return schema.Column{Name: col, Type: "integer", Length: 10}
}

func DefaultString(col string) schema.Column {
	return schema.Column{Name: col, Type: "varchar", Length: 255}
}

func DefaultBool(col string) schema.Column {
	return schema.Column{Name: col, Type: "bool"}
}

func DefaultTime(col string) schema.Column {
	return schema.Column{Name: col, Type: "timestamp"}
}

func createRecord(c *Conn, cols []string, vals []interface{}, name, pkname string) (int, error) {
	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		c.SQLTable(name),
		strings.Join(cols, ", "),
		questions(len(cols)),
	)
	if c.returning {
		sql += " RETURNING " + c.SQLColumn(name, pkname)
		var pk int
		row := c.QueryRow(sql, vals...)
		err := row.Scan(&pk)
		return pk, err
	} else {
		result, err := c.Exec(sql, vals...)
		if err != nil {
			return 0, err
		}
		
		id, err := result.LastInsertId()
		return int(id), err
	}
}

func updateRecord(c *Conn, cols []string, vals []interface{}, name, pkname string) error {
	sql := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s=?",
		c.SQLTable(name),
		strings.Join(cols, ", "),
		c.SQLColumn(name, pkname),
	)
	_, err := c.Exec(sql, vals...)
	return err
}

func deleteRecord(c *Conn, val interface{}, name, pkname string) error {
	sql := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = ?",
		c.SQLTable(name),
		c.SQLColumn(name, pkname),
	)
	_, err := c.Exec(sql, val)
	return err

}

var Schema = schema.Schema{
	Tables: map[string]*schema.Table{
		{{ range $table := .Tables }}
			"{{ .Name }}": &schema.Table{
				Name: "{{ .Name }}",
				Columns: []schema.Column{
					{{ range $column := .Columns }}
						{{ if $column.Preset }}
							{{ if eq $column.GoType "int" }}
								DefaultInt("{{ $column.Name }}"),
							{{ end }}
							{{ if eq $column.GoType "string" }}
								DefaultString("{{ $column.Name }}"),
							{{ end }}
							{{ if eq $column.GoType "bool" }}
								DefaultBool("{{ $column.Name }}"),
							{{ end }}
							{{ if eq $column.GoType "&{time.Time}" }}
								DefaultTime("{{ $column.Name }}"),
							{{ end }}
						{{ else }}
							{{ if $column.SimpleType }}
								schema.Column{
									Name: "{{ $column.Name }}",
									Type: "{{ $column.Type }}",
									Length: {{ $column.Length }},
								},
							{{ end }}
							{{ if $column.Subrecord }}
								{{ range $subcolumn := $column.Subcolumns }}
									schema.Column{
										Name: "{{ $subcolumn.Name }}",
										Type: "{{ $subcolumn.Type }}",
										Length: {{ $subcolumn.Length }},
										IncludeName: "{{ $subcolumn.IncludeName }}",
									},
								{{ end }}
							{{ end }}
						{{ end }}
					{{ end }}
				},
			},
		{{ end }}
	},
}

func init() {
	{{ range $table := .Tables }}
		{{ if $table.HasRelationship "ParentHasMany" }}
			Schema.Tables["{{ .Name }}"].HasMany = []schema.ManyRelationship{
				{{ range $relate := $table.Relations }}
					{{ if $relate.IsHasMany }}
						schema.ManyRelationship{
							Schema.Tables["{{ $relate.ParentName }}"],
							Schema.Tables["{{ $relate.ChildName }}"],
							Schema.Tables["{{ $relate.ChildName }}"].FindColumn("{{ $relate.OperativeColumn }}"),
						},
					{{ end }}
				{{ end }}
			}
		{{ end }}
	{{ end }}

	{{ range $table := .Tables }}
		{{ if $table.HasRelationship "ChildHasMany" }}
			Schema.Tables["{{ .Name }}"].ChildOf = []schema.ManyRelationship{
				{{ range $relate := $table.Relations }}
					{{ if $relate.IsChildHasMany }}
						schema.ManyRelationship{
							Schema.Tables["{{ $relate.ParentName }}"],
							Schema.Tables["{{ $relate.ChildName }}"],
							Schema.Tables["{{ $relate.ChildName }}"].FindColumn("{{ $relate.OperativeColumn }}"),
						},
					{{ end }}
				{{ end }}
			}
		{{ end }}
	{{ end }}

	{{ range $table := .Tables }}
		{{ if $table.HasRelationship "HasOne" }}
			Schema.Tables["{{ .Name }}"].HasOne = []schema.OneRelationship{
				{{ range $relate := $table.Relations }}
					{{ if $relate.IsHasOne }}
						schema.OneRelationship{
							Schema.Tables["{{ $relate.ParentName }}"],
							Schema.Tables["{{ $relate.ChildName }}"],
							Schema.Tables["{{ $relate.ChildName }}"].FindColumn("{{ $relate.OperativeColumn }}"),
						},
					{{ end }}
				{{ end }}
			}
		{{ end }}
	{{ end }}

	{{ range $table := .Tables }}
		{{ if $table.HasRelationship "BelongsTo" }}
			Schema.Tables["{{ .Name }}"].BelongsTo = []schema.OneRelationship{
				{{ range $relate := $table.Relations }}
					{{ if $relate.IsBelongsTo }}
						schema.OneRelationship{
							Schema.Tables["{{ $relate.ParentName }}"],
							Schema.Tables["{{ $relate.ChildName }}"],
							Schema.Tables["{{ $relate.ChildName }}"].FindColumn("{{ $relate.OperativeColumn }}"),
						},
					{{ end }}
				{{ end }}
			}
		{{ end }}
	{{ end }}
}

{{ range $table := .Tables }}
func (t *{{ $table.Name }}) Save(c *Conn) error {

	// check the primary key vs the zero value, if they match then
	// we will assume we have a new record
	var pkz {{ .PrimaryKeyColumn.GoType }}
	if t.{{ .PrimaryKeyColumn.Name }} == pkz {
		return t.create(c)
	} else {
		return t.update(c)
	}
}

func (t *{{ $table.Name }}) simpleCols(c *Conn) []string {
	return []string{ {{ range $column := $table.Columns }}{{ if and (ne $column.Name $table.PrimaryKeyColumn.Name) $column.SimpleType }} c.SQLColumn("{{ $table.Name }}", "{{ $column.Name }}"),{{ end }}{{ end }} }	
}

func (t *{{ $table.Name }}) simpleVals() []interface{} {
	return []interface{}{ {{ range $column := $table.Columns }}{{ if and (ne $column.Name $table.PrimaryKeyColumn.Name) $column.SimpleType }} t.{{ $column.Name }},{{ end }}{{ end }} }	
}

func (t *{{ $table.Name }}) create(c *Conn) error {
	cols := t.simpleCols(c)
	vals := t.simpleVals()
	{{ range $column := $table.Columns }}
		{{ if $column.Subrecord }}
			{{ range $subcolumn := $column.Subcolumns }}
				{{ if $subcolumn.SimpleType }}
					if t.{{ $column.Subrecord.Name }}.{{ $subcolumn.Name }}{{ $subcolumn.NonZeroCheck }} {
						vals = append(vals, t.{{ $column.Subrecord.Name }}.{{ $subcolumn.Name }})
						cols = append(cols, c.SQLColumn("{{ $table.Name }}", "{{ $subcolumn.Name }}"))
					}
				{{ end }}
			{{ end }}
		{{ end }}
	{{ end }}

	pk ,err := createRecord(c, cols, vals, "{{ $table.Name }}", "{{ $table.PrimaryKeyColumn.Name }}")
	if err == nil {
			t.{{ $table.PrimaryKeyColumn.Name }} = pk
	}
	return err
}

func (t *{{ $table.Name }}) update(c *Conn) error {
	return updateRecord(c, t.simpleCols(c), append(t.simpleVals(), t.{{ $table.PrimaryKeyColumn.Name }}), "{{ $table.Name }}", "{{ $table.PrimaryKeyColumn.Name }}")
}

func (t {{ $table.Name }}) Delete(c *Conn) error {
	return deleteRecord(c, t.{{ $table.PrimaryKeyColumn.Name }}, "{{ $table.Name }}", "{{ $table.PrimaryKeyColumn.Name }}")
}
{{ end }}
`

var sTmpl *template.Template

func init() {
	var err error
	sTmpl, err = template.New("schema").Parse(schemaTemplate)
	if err != nil {
		panic(err)
	}
}

func (pkg *Package) WriteSchemaFile(w io.Writer) {
	b := &bytes.Buffer{}
	err := sTmpl.Execute(b, pkg)
	if err != nil {
		panic(err)
	}
	ib, err := imports.Process(pkg.ActiveFiles[0].AST.Name.Name+"_schema.go", b.Bytes(), nil)
	if err != nil {
		fmt.Println("Error in Gen File:", err)
		w.Write(b.Bytes())
		return
	}
	w.Write(ib)
}
