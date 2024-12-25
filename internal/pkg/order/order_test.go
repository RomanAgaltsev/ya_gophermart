package order_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/RomanAgaltsev/ya_gophermart/internal/pkg/order"
)

var _ = Describe("Order", func() {
	DescribeTable("Checking if order number is valid with Luhn algorithm",
		func(orderNumber string, expected bool) {
			Expect(order.IsNumberValid(orderNumber)).To(Equal(expected))
		},

		EntryDescription("When order number is %q, the result is %v"),
		Entry(nil, "12345678903", true),
		Entry(nil, "98765432103", true),
		Entry(nil, "79927398713", true),
		Entry(nil, "5555 5555 5555 4444", true),
		Entry(nil, "4111111111111111", true),
		Entry(nil, "378282246310005", true),
		Entry(nil, "", false),
		Entry(nil, "order number", false),
		Entry(nil, "order #123456", false),
	)
})
