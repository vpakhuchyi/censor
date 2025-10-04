package cache

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestField is a test helper struct for testing TypeCache with field slices.
type TestField struct {
	Name     string
	IsMasked bool
}

func TestTypeCache_Set(t *testing.T) {
	t.Parallel()

	t.Run("set new value", func(t *testing.T) {
		t.Parallel()

		// GIVEN: An empty type cache.
		c := NewTypeCache[[]TestField](10)
		type TestStruct struct{ Field string }
		typ := reflect.TypeOf(TestStruct{})
		fields := []TestField{{Name: "Field", IsMasked: false}}

		// WHEN: Set is called with a new type-value pair.
		c.Set(typ, fields)

		// THEN: The value should be stored and retrievable.
		retrieved, found := c.Get(typ)
		assert.True(t, found)
		assert.Equal(t, fields, retrieved)
	})

	t.Run("set duplicate key keeps original value", func(t *testing.T) {
		t.Parallel()

		// GIVEN: A type cache with an existing value.
		c := NewTypeCache[[]TestField](10)
		type TestStruct struct{ Field string }
		typ := reflect.TypeOf(TestStruct{})
		originalFields := []TestField{{Name: "Field", IsMasked: false}}
		c.Set(typ, originalFields)

		// WHEN: Set is called again with the same type but different value.
		newFields := []TestField{{Name: "Field", IsMasked: true}}
		c.Set(typ, newFields)

		// THEN: The original value should be preserved (Set returns early for duplicates).
		retrieved, found := c.Get(typ)
		assert.True(t, found)
		assert.Equal(t, originalFields, retrieved)
	})

	t.Run("evict oldest entry when capacity exceeded", func(t *testing.T) {
		t.Parallel()

		// GIVEN: A type cache at full capacity.
		c := NewTypeCache[[]TestField](3)
		type TestStruct1 struct{ Field1 string }
		type TestStruct2 struct{ Field2 int }
		type TestStruct3 struct{ Field3 bool }
		type TestStruct4 struct{ Field4 float64 }

		type1 := reflect.TypeOf(TestStruct1{})
		type2 := reflect.TypeOf(TestStruct2{})
		type3 := reflect.TypeOf(TestStruct3{})
		type4 := reflect.TypeOf(TestStruct4{})

		fields1 := []TestField{{Name: "Field1", IsMasked: false}}
		fields2 := []TestField{{Name: "Field2", IsMasked: true}}
		fields3 := []TestField{{Name: "Field3", IsMasked: false}}
		fields4 := []TestField{{Name: "Field4", IsMasked: true}}

		c.Set(type1, fields1)
		c.Set(type2, fields2)
		c.Set(type3, fields3)

		// WHEN: A new entry is added when cache is at capacity.
		c.Set(type4, fields4)

		// THEN: The oldest entry should be evicted and the new entry should be present.
		_, found1 := c.Get(type1)
		assert.False(t, found1, "oldest entry should have been evicted")

		retrieved2, found2 := c.Get(type2)
		assert.True(t, found2)
		assert.Equal(t, fields2, retrieved2)

		retrieved3, found3 := c.Get(type3)
		assert.True(t, found3)
		assert.Equal(t, fields3, retrieved3)

		retrieved4, found4 := c.Get(type4)
		assert.True(t, found4)
		assert.Equal(t, fields4, retrieved4)
	})
}

func TestTypeCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("get existing value", func(t *testing.T) {
		t.Parallel()

		// GIVEN: A type cache with a pre-stored value.
		c := NewTypeCache[[]TestField](10)
		type TestStruct struct{ Field string }
		typ := reflect.TypeOf(TestStruct{})
		fields := []TestField{{Name: "Field", IsMasked: true}}
		c.Set(typ, fields)

		// WHEN: Get is called for the existing type.
		retrieved, found := c.Get(typ)

		// THEN: The value should be returned with found=true.
		assert.True(t, found)
		assert.Equal(t, fields, retrieved)
	})

	t.Run("get non-existing value", func(t *testing.T) {
		t.Parallel()

		// GIVEN: An empty type cache.
		c := NewTypeCache[[]TestField](10)
		type TestStruct struct{ Field string }
		typ := reflect.TypeOf(TestStruct{})

		// WHEN: Get is called for a non-existing type.
		retrieved, found := c.Get(typ)

		// THEN: A zero value should be returned with found=false.
		assert.False(t, found)
		assert.Nil(t, retrieved)
	})
}

func TestTypeCache_DifferentPackagesSameTypeName(t *testing.T) {
	t.Parallel()

	t.Run("distinguish types with same name from different packages", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Two different struct types with the same field layout but different package paths.
		c := NewTypeCache[[]TestField](10)

		type User1 struct {
			Name  string `censor:"display"`
			Email string
		}

		type User2 struct {
			Name  string `censor:"display"`
			Email string
		}

		type1 := reflect.TypeOf(User1{})
		type2 := reflect.TypeOf(User2{})

		fields1 := []TestField{
			{Name: "Name", IsMasked: false},
			{Name: "Email", IsMasked: true},
		}
		fields2 := []TestField{
			{Name: "Name", IsMasked: true},
			{Name: "Email", IsMasked: false},
		}

		// WHEN: Values are set for both types.
		c.Set(type1, fields1)
		c.Set(type2, fields2)

		// THEN: Each type should maintain its own separate cache entry.
		assert.NotEqual(t, type1, type2)
		assert.NotEqual(t, type1.String(), type2.String())

		retrieved1, found1 := c.Get(type1)
		assert.True(t, found1)
		assert.Equal(t, fields1, retrieved1)

		retrieved2, found2 := c.Get(type2)
		assert.True(t, found2)
		assert.Equal(t, fields2, retrieved2)

		assert.NotEqual(t, retrieved1, retrieved2)
	})
}

func TestTypeCache_AnonymousStructs(t *testing.T) {
	t.Parallel()

	t.Run("treat identical anonymous structs as same type", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Two anonymous struct values with identical structure.
		c := NewTypeCache[[]TestField](10)

		val1 := struct {
			Name  string
			Email string
		}{Name: "John", Email: "john@example.com"}

		val2 := struct {
			Name  string
			Email string
		}{Name: "Jane", Email: "jane@example.com"}

		type1 := reflect.TypeOf(val1)
		type2 := reflect.TypeOf(val2)

		fields := []TestField{
			{Name: "Name", IsMasked: false},
			{Name: "Email", IsMasked: true},
		}

		// WHEN: A value is set using the first type.
		c.Set(type1, fields)

		// THEN: The same value should be retrievable using the second type.
		assert.Equal(t, type1, type2)

		retrieved1, found1 := c.Get(type1)
		assert.True(t, found1)
		assert.Equal(t, fields, retrieved1)

		retrieved2, found2 := c.Get(type2)
		assert.True(t, found2)
		assert.Equal(t, fields, retrieved2)
	})
}

func TestTypeCache_PointerVsValue(t *testing.T) {
	t.Parallel()

	t.Run("distinguish pointer type from value type", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Both pointer and value types of the same struct.
		c := NewTypeCache[[]TestField](10)

		type TestStruct struct{ Field string }

		typeValue := reflect.TypeOf(TestStruct{})
		typePointer := reflect.TypeOf(&TestStruct{})

		fieldsValue := []TestField{{Name: "Field", IsMasked: false}}
		fieldsPointer := []TestField{{Name: "Field", IsMasked: true}}

		// WHEN: Both types are added to the cache.
		c.Set(typeValue, fieldsValue)
		c.Set(typePointer, fieldsPointer)

		// THEN: Both should be stored independently with different values.
		assert.NotEqual(t, typeValue, typePointer)

		retrievedValue, foundValue := c.Get(typeValue)
		assert.True(t, foundValue)
		assert.Equal(t, fieldsValue, retrievedValue)

		retrievedPointer, foundPointer := c.Get(typePointer)
		assert.True(t, foundPointer)
		assert.Equal(t, fieldsPointer, retrievedPointer)
	})
}

func BenchmarkTypeCache_Get(b *testing.B) {
	c := NewTypeCache[[]TestField](100)

	type TestStruct struct {
		Name  string
		Email string
		Age   int
	}

	typ := reflect.TypeOf(TestStruct{})
	fields := []TestField{
		{Name: "Name", IsMasked: false},
		{Name: "Email", IsMasked: true},
		{Name: "Age", IsMasked: false},
	}
	c.Set(typ, fields)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = c.Get(typ)
	}
}

func BenchmarkTypeCache_Set(b *testing.B) {
	type TestStruct struct {
		Name  string
		Email string
		Age   int
	}

	typ := reflect.TypeOf(TestStruct{})
	fields := []TestField{
		{Name: "Name", IsMasked: false},
		{Name: "Email", IsMasked: true},
		{Name: "Age", IsMasked: false},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		c := NewTypeCache[[]TestField](100)
		c.Set(typ, fields)
	}
}
