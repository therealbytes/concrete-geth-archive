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

func New{{.RowStructName}}WithHooks(row *{{.RowStructName}}) *{{.RowStructName}} {
	return &{{.RowStructName}}{lib.NewDatastoreStructWithHooks(row.IDatastoreStruct)}
}

type {{.TableStructName}}WithHooks struct {
	{{.TableStructName}}
	OnSetRow func(rowKey []interface{}, column int, value []byte)
}

func New{{.TableStructName}}WithHooks(table *{{.TableStructName}}) *{{.TableStructName}}WithHooks {
	return &{{.TableStructName}}WithHooks{
		{{.TableStructName}}: *table,
	}
}

func (m *{{.TableStructName}}WithHooks) Get(
{{- range .Schema.Keys }}
	{{.Name}} {{.Type.GoType}},
{{- end }}
) *{{.RowStructName}} {
	row := m.{{.TableStructName}}.Get(
	{{- range .Schema.Keys }}
		{{.Name}},
	{{- end }}
	)
	row = New{{.RowStructName}}WithHooks(row)
	row.IDatastoreStruct.(*lib.DatastoreStructWithHooks).OnSetField = func(column int, value []byte) {
		if m.OnSetRow != nil {
			m.OnSetRow([]interface{}{
				{{- range .Schema.Keys }}
					{{.Name}},
				{{- end }}
			}, column, value)
		}
	}
	return row
}
