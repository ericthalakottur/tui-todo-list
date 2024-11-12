package utils

import (
	"errors"
	"testing"
)

func TestModulus(t *testing.T) {
	input := []int{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5}
	expected := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0}
	m := 5

	t.Run("modulus with positive integer", func(t *testing.T) {
		for i := range input {
			result, err := Modulus(input[i], m)

			if result != expected[i] || err != nil {
				t.Errorf("Expected %d as result of %d %% %d", result, input[i], m)
			}
		}
	})

	t.Run("modulus with 0", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Expected the function to return an error")
			}
		}()

		_, err := Modulus(10, 0)

		if !errors.Is(err, DivideByZeroError) {
			t.Errorf("Expected the function to return a DivideByZeroError")
		}

	})
}
