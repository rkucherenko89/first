package data

import "testing"

func TestChecksValidation(t *testing.T) {
	product := &Product{
		Name:  "Tea",
		Price: 2,
	}

	err := product.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
