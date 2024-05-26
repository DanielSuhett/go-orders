package order_test

import (
	"slices"
	"testing"

	"github.com/DanielSuhett/orders/internal/order"
	"github.com/DanielSuhett/orders/pkg/database"
)

func TestSave(t *testing.T) {
	tests := []struct {
		name                 string
		category             string
		value                int
		paymentMethod        string
		expectedValue        int
		expectDiscount       bool
		expectedFreeShipping bool
	}{
		{
			name:                 "With discount for Boleto",
			category:             "appliances",
			value:                800,
			paymentMethod:        order.BrazilianBoletoPayment,
			expectDiscount:       true,
			expectedFreeShipping: false,
		},
		{
			name:                 "No discount for other methods",
			category:             "kids",
			value:                800,
			paymentMethod:        order.BrazilianPixPayment,
			expectDiscount:       false,
			expectedFreeShipping: false,
		},
		{
			name:                 "With free shipping",
			category:             "-",
			value:                1500,
			paymentMethod:        order.BrazilianPixPayment,
			expectDiscount:       false,
			expectedFreeShipping: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := database.Pool()

			if err != nil {
				t.Fatal(err)
			}

			no, err := order.NewOrder(tt.category, tt.value, tt.paymentMethod, tt.value)

			if err != nil {
				t.Fatal(err)
			}

			id, err := order.Save(pool, no)

			if err != nil {
				t.Fatal(err)
			}

			no, err = order.Get(pool, *id)

			if err != nil {
				t.Fatal(err)
			}

			if tt.expectDiscount && no.Payment.Value == tt.value {
				t.Errorf("Expected discount to be applied, but it was not")
			}

			if !tt.expectDiscount && no.Payment.Value != tt.value {
				t.Errorf("Expected no discount to be applied, but discount was applied")
			}

			if tt.expectDiscount && no.Payment.Value == tt.value {
				t.Errorf("Expected discount to be applied, but it was not")
			}

			if !tt.expectedFreeShipping && slices.Contains(no.Labels, "free-shipping") {
				t.Errorf("Unexpected free shipping applied")
			}

			if tt.expectedFreeShipping && !slices.Contains(no.Labels, "free-shipping") {
				t.Errorf("Expected free shipping to be applied, but it was not")
			}

			if tt.category == "kids" && !slices.Contains(no.Labels, "gift") {
				t.Errorf("Expected gift to be applied, but it was not")
			}

			if tt.category == "appliances" && !slices.Contains(no.Labels, "fragile") {
				t.Errorf("Expected fragile to be applied, but it was not")
			}

		})
	}
}
