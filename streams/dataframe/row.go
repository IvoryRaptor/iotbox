package dataframe

import "strings"

type Row interface {
	GetStructField(name string) (int, *StructField)
	GetValue(name string) (*StructField, interface{})
	GetValues() []interface{}
}

type GenericRow struct {
	Device    string
	Timestamp int64
	Schema    *Schema
	Values    []interface{}
}

func (row *GenericRow) GetStructField(name string) (int, *StructField) {
	for index, field := range row.Schema.Fields {
		if strings.EqualFold(name, field.Name) {
			return index, field
		}
	}
	return -1, nil
}
func (row *GenericRow) GetValues() []interface{} {
	return row.Values
}
func (row *GenericRow) GetValue(name string) (*StructField, interface{}) {
	for index, field := range row.Schema.Fields {
		if strings.EqualFold(name, field.Name) {
			return field, row.Values[index]
		}
	}
	return nil, nil
}

type JoinedRow struct {
}
