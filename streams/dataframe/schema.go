package dataframe

import "fmt"

type DataType interface {
	ValueToString(value interface{}) string
	Equal(value1 interface{}, value2 interface{}) bool
}
type BooleanType struct {
}

func (t *BooleanType) ValueToString(value interface{}) string {
	return fmt.Sprint(value.(bool))
}

func (t *BooleanType) Equal(value1 interface{}, value2 interface{}) bool {
	return value1.(bool) == value2.(bool)
}

type ByteType struct {
}
type ShortType struct {
}
type IntegerType struct {
}

func (t *IntegerType) ValueToString(value interface{}) string {
	return fmt.Sprint(value.(int))
}

func (t *IntegerType) Equal(value1 interface{}, value2 interface{}) bool {
	return value1.(int) == value2.(int)
}

type LongType struct {
}
type FloatType struct {
}
type DoubleType struct {
}
type StringType struct {
}
type DateType struct {
}
type DecimalType struct {
}
type TimestampType struct {
}
type BinaryType struct {
}
type ArrayType struct {
}

type StructField struct {
	Name     string
	DataType DataType
	Nullable bool
}

type Schema struct {
	Name   string
	Fields []*StructField
}
