/* Autogenerated file. Do not edit manually. */

package {{.Package}}

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/codegen/datamod"
	"github.com/ethereum/go-ethereum/concrete/crypto"
	"github.com/ethereum/go-ethereum/concrete/lib"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = common.Big1
	_ = datamod.EncodeAddress
)

var (
	{{.TableStructName}}DefaultKey = crypto.Keccak256([]byte("datamod.v1.{{.TableStructName}}"))
)

type {{.RowStructName}} struct {
	lib.StorageStruct
}

func New{{.RowStructName}}(slot lib.StorageSlot) *{{.RowStructName}} {
	sizes := {{.SizesStr}}
	return &{{.RowStructName}}{*lib.NewStorageStruct(slot, sizes)}
}

func (v *{{$.RowStructName}}) Get() (
{{- range .Schema.Values }}
	{{.Type.GoType}},
{{- end }}
) {
	return {{ range .Schema.Values -}}
		datamod.{{.Type.DecodeFunc}}(
			{{- .Type.Size }}, v.GetField({{.Index}}))
			{{- if eq .Index (sub (len $.Schema.Values) 1) }}{{else}}, {{end}}
	{{- end }}
}

func (v *{{$.RowStructName}}) Set(
{{- range .Schema.Values }}
	{{.Name}} {{.Type.GoType}},
{{- end }}
) {
{{- range .Schema.Values }}
	v.SetField({{.Index}}, datamod.{{.Type.EncodeFunc}}({{.Type.Size}}, {{.Name}}))
{{- end }}
}
{{range .Schema.Values}}
func (v *{{$.RowStructName}}) Get{{.Title}}() {{.Type.GoType}} {
	data := v.GetField({{.Index}})
	return datamod.{{.Type.DecodeFunc}}({{.Type.Size}}, data)
}

func (v *{{$.RowStructName}}) Set{{.Title}}(value {{.Type.GoType}}) {
	data := datamod.{{.Type.EncodeFunc}}({{.Type.Size}}, value)
	v.SetField({{.Index}}, data)
}
{{end}}
{{- if .Schema.Keys }}
type {{.TableStructName}} struct {
	mapping lib.Mapping
}

func New{{.TableStructName}}(ds lib.Datastore) *{{.TableStructName}} {
	return &{{.TableStructName}}{ds.Mapping({{.TableStructName}}DefaultKey)}
}

func New{{.TableStructName}}WithKey(ds lib.Datastore, key []byte) *{{.TableStructName}} {
	return &{{.TableStructName}}{ds.Mapping(key)}
}

func (m *{{.TableStructName}}) Get(
{{- range .Schema.Keys }}
	{{.Name}} {{.Type.GoType}},
{{- end }}
) *{{.RowStructName}} {
	return New{{.RowStructName}}(
		m.mapping.
		{{- range .Schema.Keys -}}
		{{- if eq .Index (sub (len $.Schema.Keys) 1) -}}
			Value(datamod.{{.Type.EncodeFunc}}({{.Type.Size}}, {{.Name}})),
		{{- else -}}
			Mapping(datamod.{{.Type.EncodeFunc}}({{.Type.Size}}, {{.Name}})).
		{{- end -}}
		{{end}}
	)
}
{{- else }}
type {{.TableStructName}} = {{.RowStructName}}

func New{{.TableStructName}}(ds lib.Datastore) *{{.TableStructName}} {
	return New{{.RowStructName}}(ds.Value({{.TableStructName}}DefaultKey))
}

func New{{.TableStructName}}WithKey(ds lib.Datastore, key []byte) *{{.TableStructName}} {
	return New{{.RowStructName}}(ds.Value(key))
}
{{- end }}