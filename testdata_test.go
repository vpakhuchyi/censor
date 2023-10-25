package sanitiser

type person struct {
	ID           string `log:"display"`
	Name         string
	Age          int `log:"display"`
	Email        string
	StringArray  [3]string  `log:"display"`
	String       []string   `log:"display"`
	Integers     []int      `log:"display"`
	Floats       []float64  `log:"display"`
	Bools        []bool     `log:"display"`
	Addresses    []address  `log:"display"`
	Address      address    `log:"display"`
	TaxAddress   *address   `log:"display"`
	TaxAddresses []*address `log:"display"`
	Container    container  `log:"display"`
}

type address struct {
	City   string `json:"city" log:"display"`
	State  string `json:"state" log:"display"`
	Street string `json:"street"`
	Zip    string `json:"zip"`
}

type structWithPrimitives struct {
	Int64   int64   `log:"display"`
	Int32   int32   `log:"display"`
	Int16   int16   `log:"display"`
	Int8    int8    `log:"display"`
	Int     int     `log:"display"`
	Uint64  uint64  `log:"display"`
	Uint32  uint32  `log:"display"`
	Uint16  uint16  `log:"display"`
	Uint8   uint8   `log:"display"`
	Uint    uint    `log:"display"`
	Bool    bool    `log:"display"`
	Rune    rune    `log:"display"`
	Byte    byte    `log:"display"`
	Float64 float64 `log:"display"`
	Float32 float32 `log:"display"`
	String  string  `log:"display"`
}

type structWithContainersFields struct {
	StringSlice  []string  `log:"display"`
	IntSlice     []int     `log:"display"`
	FloatSlice   []float64 `log:"display"`
	BoolSlice    []bool    `log:"display"`
	StructSlice  []address `log:"display"`
	PointerSlice []*int    `log:"display"`
	ArraySlice   [2]string `log:"display"`
}

type structWithComplexFields struct {
	Slice       []address `log:"display"`
	MaskedSlice []address
	Map         map[string]address `log:"display"`
	Array       [2]address         `log:"display"`
	Ptr         *address           `log:"display"`
	Struct      address            `log:"display"`
}

type container struct {
	Persons []person `log:"display"`
}
