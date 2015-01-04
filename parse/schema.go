package parse

var schemaTemplate = `/*
  This code was generated by the Doctor ORM Generator and isn't meant to be edited.
	If at all possible, please regenerate this file from your gp files instead of
	attempting to edit it to add changes.
*/

package {{ .Name }}

import "github.com/acsellers/dr/schema"

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
				Index: []schema.Index{
					schema.Index{
						Columns: []string{"{{ $table.PrimaryKeyColumn.Name }}"},
					},
					{{ range $index := $table.Indexes }}
						schema.Index{
							Columns: []string{ {{ range .Columns }}
								"{{ . }}",{{ end }}
							},
						},
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

func (t {{ $table.Name }}) Scope() {{ $table.Name }}Scope {
	return t.cached_conn.{{ $table.Name }}.{{ .PrimaryKeyColumn.Name }}().Eq(t.{{ .PrimaryKeyColumn.Name }})
}

func (t {{ $table.Name }}) ToScope(c *Conn) {{ $table.Name }}Scope {
	return c.{{ $table.Name }}.{{ .PrimaryKeyColumn.Name }}().Eq(t.{{ .PrimaryKeyColumn.Name }})
}

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
			t.cached_conn = c
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
