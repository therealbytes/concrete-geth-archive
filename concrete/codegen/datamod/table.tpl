/* Autogenerated file. Do not edit manually. */

package {{.Package}}

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/codegen/datamod/codec"
	"github.com/ethereum/go-ethereum/concrete/crypto"
	"github.com/ethereum/go-ethereum/concrete/lib"
)

// Reference imports to suppress errors if they are not used.
var (
	_ = crypto.Keccak256
	_ = big.NewInt
	_ = common.Big1
	_ = codec.EncodeAddress
)

// var (
//	{{.TableStructName}}DefaultKey = crypto.Keccak256([]byte("datamod.v1.{{.TableStructName}}"))
// )

func {{.TableStructName}}DefaultKey() []byte {
	return crypto.Keccak256([]byte("datamod.v1.{{.TableStructName}}"))
}

type {{.RowStructName}} struct {
	lib.IDatastoreStruct
}

func New{{.RowStructName}}(dsSlot lib.DatastoreSlot) *{{.RowStructName}} {
	sizes := {{.SizesStr}}
	return &{{.RowStructName}}{lib.NewDatastoreStruct(dsSlot, sizes)}
}

func (v *{{$.RowStructName}}) Get() (
{{- range .Schema.Values }}
	{{ .Name }} {{if eq .Type.Type 2}}*{{end}}{{.Type.GoType}},
{{- end }}
) {
	return {{ range .Schema.Values }}
		{{- if lt .Type.Type 2 -}}
		codec.{{.Type.DecodeFunc}}({{.Type.Size}}, {{if eq .Type.Type 0}}v.GetField{{else}}v.GetField_bytes{{end}}({{.Index}}))
		{{- else -}}
		New{{.Type.GoType}}FromSlot(v.GetField_slot({{.Index}}))
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
func (v *{{$.RowStructName}}) Get{{.PascalCase}}() {{.Type.GoType}} {
	data := {{if eq .Type.Type 0}}v.GetField{{else}}v.GetField_bytes{{end}}({{.Index}})
	return codec.{{.Type.DecodeFunc}}({{.Type.Size}}, data)
}

func (v *{{$.RowStructName}}) Set{{.PascalCase}}(value {{.Type.GoType}}) {
	data := codec.{{.Type.EncodeFunc}}({{.Type.Size}}, value)
	{{if eq .Type.Type 0}}v.SetField{{else}}v.SetField_bytes{{end}}({{.Index}}, data)
}
{{ else }}
func (v *{{$.RowStructName}}) Get{{.PascalCase}}() *{{.Type.GoType}} {
	dsSlot := v.GetField_slot({{.Index}})
	return New{{.Type.GoType}}FromSlot(dsSlot)
}
{{ end}}
{{- end}}
type {{.TableStructName}} struct {
	dsSlot lib.DatastoreSlot
}

func New{{.TableStructName}}(ds lib.Datastore) *{{.TableStructName}} {
	dsSlot := ds.Get({{.TableStructName}}DefaultKey())
	return &{{.TableStructName}}{dsSlot}
}

func New{{.TableStructName}}FromSlot(dsSlot lib.DatastoreSlot) *{{.TableStructName}} {
	return &{{.TableStructName}}{dsSlot}
}

{{- if .Schema.Keys }}
func (m *{{.TableStructName}}) Get(
{{- range .Schema.Keys }}
	{{.Name}} {{.Type.GoType}},
{{- end }}
) *{{.RowStructName}} {
	dsSlot := m.dsSlot.Mapping().GetNested(
		{{- range .Schema.Keys }}
		codec.{{.Type.EncodeFunc}}({{.Type.Size}}, {{.Name}}),
		{{- end }}
	)
	return New{{.RowStructName}}(dsSlot)
}
{{- else }}
func (m *{{.TableStructName}}) Get() *{{.RowStructName}} {
	return New{{.RowStructName}}(m.dsSlot)
}
{{- end }}