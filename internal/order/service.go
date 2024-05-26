package order

import (
	"github.com/google/uuid"
)

const (
	BrazilianBoletoPayment = "BOLETO"
	BrazilianPixPayment    = "PIX"
)

type Product struct {
	Category string
	Value    int
}

type Payment struct {
	Method string
	Value  int
}

type Order struct {
	Id      uuid.UUID
	Product Product
	Payment Payment
	Labels  []string
}

type orderService struct {
	repository orderRepository
}

func newOrderService(repo orderRepository) orderService {
	return orderService{repo}
}

func markFragileProducts(order *Order) {
	if order.Product.Category == "appliances" {
		order.Labels = append(order.Labels, "fragile")
	}
}

func markFreeShipping(order *Order) {
	if order.Product.Value > 1000 {
		order.Labels = append(order.Labels, "free-shipping")
	}
}

func markGift(order *Order) {
	if order.Product.Category == "kids" {
		order.Labels = append(order.Labels, "gift")
	}
}

func setDiscounts(order *Order) {
	if order.Payment.Method == BrazilianBoletoPayment {
		p10off := float64(order.Payment.Value) * 0.10
		order.Payment.Value = order.Payment.Value - int(p10off)
	}
}

func (os orderService) save(order Order) error {
	markFragileProducts(&order)
	markFreeShipping(&order)
	markGift(&order)

	setDiscounts(&order)

	return os.repository.Save(&order)
}

func (os orderService) list() ([]Order, error) {
	return os.repository.list()
}

func (os orderService) get(id uuid.UUID) (Order, error) {
	return os.repository.get(id)
}
