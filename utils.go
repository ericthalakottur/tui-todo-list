package main

import "errors"

func modulus(a, m int) (int, error) {
	if m == 0 {
		return -1, errors.New("Division by 0")
	}
	return (a%m + m) % m, nil
}
