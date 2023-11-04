package censor

import "fmt"

type person struct {
	ID           string `censor:"display"`
	Name         string
	Age          int `censor:"display"`
	Email        string
	StringArray  [3]string  `censor:"display"`
	String       []string   `censor:"display"`
	Integers     []int      `censor:"display"`
	Floats       []float64  `censor:"display"`
	Bools        []bool     `censor:"display"`
	Addresses    []address  `censor:"display"`
	Address      address    `censor:"display"`
	TaxAddress   *address   `censor:"display"`
	TaxAddresses []*address `censor:"display"`
	Container    container  `censor:"display"`
}

type address struct {
	City   string `json:"city" censor:"display"`
	State  string `json:"state" censor:"display"`
	Street string `json:"street"`
	Zip    string `json:"zip"`
}

type structWithPrimitives struct {
	Int64      int64      `censor:"display"`
	Int32      int32      `censor:"display"`
	Int16      int16      `censor:"display"`
	Int8       int8       `censor:"display"`
	Int        int        `censor:"display"`
	Uint64     uint64     `censor:"display"`
	Uint32     uint32     `censor:"display"`
	Uint16     uint16     `censor:"display"`
	Uint8      uint8      `censor:"display"`
	Uint       uint       `censor:"display"`
	Bool       bool       `censor:"display"`
	Rune       rune       `censor:"display"`
	Byte       byte       `censor:"display"`
	Float64    float64    `censor:"display"`
	Float32    float32    `censor:"display"`
	String     string     `censor:"display"`
	Complex64  complex64  `censor:"display"`
	Complex128 complex128 `censor:"display"`
}

type structWithContainersFields struct {
	StringSlice  []string  `censor:"display"`
	IntSlice     []int     `censor:"display"`
	FloatSlice   []float64 `censor:"display"`
	BoolSlice    []bool    `censor:"display"`
	StructSlice  []address `censor:"display"`
	PointerSlice []*int    `censor:"display"`
	ArraySlice   [2]string `censor:"display"`
}

type structWithComplexFields struct {
	Slice       []address `censor:"display"`
	MaskedSlice []address
	Map         map[string]address `censor:"display"`
	Array       [2]address         `censor:"display"`
	Ptr         *address           `censor:"display"`
	Struct      address            `censor:"display"`
	Interface   interface{}        `censor:"display"`
}

type structWithInterface struct {
	Interface interface{} `censor:"display"`
}

type container struct {
	Persons []person `censor:"display"`
}

type Printer interface {
	Print()
}

type printer struct {
	Name string `censor:"display"`
}

func (p *printer) Print() {
	fmt.Println(p.Name)
}

type sensPrinter struct {
	Name string
}

func (s *sensPrinter) Print() {
	fmt.Println(s.Name)
}
