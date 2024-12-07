package order

import "strconv"

func IsNumberValid(orderNumber string) bool {
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
