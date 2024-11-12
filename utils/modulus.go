package utils

import "errors"

var DivideByZeroError = errors.New("Divide by Zero")

func Modulus(a, m int) (int, error) {
	if m == 0 {
		return -1, DivideByZeroError
	}
	return (a%m + m) % m, nil
}
