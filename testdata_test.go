package sanitiser

import "fmt"

type person struct {
	ID           string `sanitiser:"display"`
	Name         string
	Age          int `sanitiser:"display"`
	Email        string
	StringArray  [3]string  `sanitiser:"display"`
	String       []string   `sanitiser:"display"`
	Integers     []int      `sanitiser:"display"`
	Floats       []float64  `sanitiser:"display"`
	Bools        []bool     `sanitiser:"display"`
	Addresses    []address  `sanitiser:"display"`
	Address      address    `sanitiser:"display"`
	TaxAddress   *address   `sanitiser:"display"`
	TaxAddresses []*address `sanitiser:"display"`
	Container    container  `sanitiser:"display"`
}

type address struct {
	City   string `json:"city" sanitiser:"display"`
	State  string `json:"state" sanitiser:"display"`
	Street string `json:"street"`
	Zip    string `json:"zip"`
}

type structWithPrimitives struct {
	Int64   int64   `sanitiser:"display"`
	Int32   int32   `sanitiser:"display"`
	Int16   int16   `sanitiser:"display"`
	Int8    int8    `sanitiser:"display"`
	Int     int     `sanitiser:"display"`
	Uint64  uint64  `sanitiser:"display"`
	Uint32  uint32  `sanitiser:"display"`
	Uint16  uint16  `sanitiser:"display"`
	Uint8   uint8   `sanitiser:"display"`
	Uint    uint    `sanitiser:"display"`
	Bool    bool    `sanitiser:"display"`
	Rune    rune    `sanitiser:"display"`
	Byte    byte    `sanitiser:"display"`
	Float64 float64 `sanitiser:"display"`
	Float32 float32 `sanitiser:"display"`
	String  string  `sanitiser:"display"`
}

type structWithContainersFields struct {
	StringSlice  []string  `sanitiser:"display"`
	IntSlice     []int     `sanitiser:"display"`
	FloatSlice   []float64 `sanitiser:"display"`
	BoolSlice    []bool    `sanitiser:"display"`
	StructSlice  []address `sanitiser:"display"`
	PointerSlice []*int    `sanitiser:"display"`
	ArraySlice   [2]string `sanitiser:"display"`
}

type structWithComplexFields struct {
	Slice       []address `sanitiser:"display"`
	MaskedSlice []address
	Map         map[string]address `sanitiser:"display"`
	Array       [2]address         `sanitiser:"display"`
	Ptr         *address           `sanitiser:"display"`
	Struct      address            `sanitiser:"display"`
	Interface   interface{}        `sanitiser:"display"`
}

type structWithInterface struct {
	Interface interface{} `sanitiser:"display"`
}

type container struct {
	Persons []person `sanitiser:"display"`
}

type Printer interface {
	Print()
}

type printer struct {
	Name string `sanitiser:"display"`
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
