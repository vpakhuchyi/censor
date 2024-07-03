package censor

import (
	"fmt"
	"testing"
)

func TestProcessor_Format(t *testing.T) {
	type Address struct {
		Country string `censor:"display"`
		City    string
	}

	type User struct {
		Name  string `censor:"display"`
		Email string
		Map   map[string]string `censor:"display"`
		Address
	}

	u := User{
		Name:  "John Doe",
		Email: "example@gmail.com",
		Map: map[string]string{
			"mapKey": "mapVal",
			"2":      "two",
		},
		Address: Address{
			Country: "Ukraine",
			City:    "Kharkiv",
		},
	}

	m := map[string]User{"1": u}

	cfg := Config{
		PrintConfigOnInit: true,
		MaskValue:         "[CENSORED]",
		ExcludePatterns:   []string{"[0-9]"},
	}

	p, err := NewWithOpts(WithConfig(&cfg))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p.Format(m))
}
