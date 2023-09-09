/* Autogenerated file. Do not edit manually. */

package {{.Package}}

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/codegen/datamod/codec"
	"github.com/ethereum/go-ethereum/concrete/crypto"
	"github.com/ethereum/go-ethereum/concrete/lib"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = common.Big1
	_ = codec.EncodeAddress
)

var (
	{{.TableStructName}}DefaultKey = crypto.Keccak256([]byte("datamod.v1.{{.TableStructName}}"))
)

type {{.RowStructName}} struct {
	lib.DatastoreStruct
}

func New{{.RowStructName}}(store lib.DatastoreSlot) *{{.RowStructName}} {
	sizes := {{.SizesStr}}
	return &{{.RowStructName}}{*lib.NewDatastoreStruct(store, sizes)}
}

func (v *{{$.RowStructName}}) Get() (
{{- range .Schema.Values }}
	{{if eq .Type.Type 2}}*{{end}}{{.Type.GoType}},
{{- end }}
) {
	return {{ range .Schema.Values }}
		{{- if lt .Type.Type 2 -}}
		codec.{{.Type.DecodeFunc}}({{.Type.Size}}, {{if eq .Type.Type 0}}v.GetField{{else}}v.GetField_bytes{{end}}({{.Index}}))
		{{- else -}}
		&{{.Type.GoType}}{v.GetField_slot({{.Index}})}
		{{- end }}
		{{- if ne .Index (sub (len $.Schema.Values) 1) }},
		{{end}}
	{{- end }}
}

func (v *{{$.RowStructName}}) Set(
{{- range .Schema.Values }}
{{- if lt .Type.Type 2 }}
	{{.Name}} {{.Type.GoType}},
{{- end }}
{{- end }}
) {
{{- range .Schema.Values }}
{{- if lt .Type.Type 2 }}
	{{if eq .Type.Type 0}}v.SetField{{else if eq .Type.Type 1}}v.SetField_bytes{{end -}}
	({{ .Index }}, codec.{{.Type.EncodeFunc}}({{.Type.Size}}, {{.Name}}))
{{- end }}
{{- end }}
}
{{range .Schema.Values}}
{{- if lt .Type.Type 2 }}
func (v *{{$.RowStructName}}) Get{{.Title}}() {{.Type.GoType}} {
	data := {{if eq .Type.Type 0}}v.GetField{{else}}v.GetField_bytes{{end}}({{.Index}})
	return codec.{{.Type.DecodeFunc}}({{.Type.Size}}, data)
}

func (v *{{$.RowStructName}}) Set{{.Title}}(value {{.Type.GoType}}) {
	data := codec.{{.Type.EncodeFunc}}({{.Type.Size}}, value)
	{{if eq .Type.Type 0}}v.SetField{{else}}v.SetField_bytes{{end}}({{.Index}}, data)
}
{{ else }}
func (v *{{$.RowStructName}}) Get{{.Title}}() *{{.Type.GoType}} {
	dsSlot := v.GetField_slot({{.Index}})
	return &{{.Type.GoType}}{dsSlot}
}
{{ end}}
{{- end}}
{{- if .Schema.Keys }}
type {{.TableStructName}} struct {
	store lib.DatastoreSlot
}

func New{{.TableStructName}}(ds lib.Datastore) *{{.TableStructName}} {
	return &{{.TableStructName}}{ds.Get({{.TableStructName}}DefaultKey)}
}

func New{{.TableStructName}}WithKey(ds lib.Datastore, key []byte) *{{.TableStructName}} {
	return &{{.TableStructName}}{ds.Get(key)}
}

func (m *{{.TableStructName}}) Get(
{{- range .Schema.Keys }}
	{{.Name}} {{.Type.GoType}},
{{- end }}
) *{{.RowStructName}} {
	store := m.store.Mapping().GetNested(
		{{- range .Schema.Keys }}
		codec.{{.Type.EncodeFunc}}({{.Type.Size}}, {{.Name}}),
		{{- end }}
	)
	return New{{.RowStructName}}(store)
}
{{- else }}
type {{.TableStructName}} = {{.RowStructName}}

func New{{.TableStructName}}(ds lib.Datastore) *{{.TableStructName}} {
	store := ds.Get({{.TableStructName}}DefaultKey)
	return New{{.RowStructName}}(store)
}

func New{{.TableStructName}}WithKey(ds lib.Datastore, key []byte) *{{.TableStructName}} {
	store := ds.Get(key)
	return New{{.RowStructName}}(store)
}
{{- end }}
