package order

import (
	"strconv"
	"strings"
)

// IsNumberValid checks if a given order numder is valid with Luhn algorithm.
func IsNumberValid(orderNumber string) bool {
	if orderNumber == "" {
		return false
	}

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/*-+!@#$%^&*()-_=+{}[]<>?|"
	if strings.ContainsAny(orderNumber, letters) {
		return false
	}

	var sum int

	parity := len(orderNumber) % 2

	for i := 0; i < len(orderNumber); i++ {
		digit, _ := strconv.Atoi(string(orderNumber[i]))
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
