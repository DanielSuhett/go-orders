package order

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewOrder(productCategory string, productValue int, paymentMethod string, paymentValue int) (Order, error) {
	id := uuid.New()

	if productCategory == "" || paymentMethod == "" || paymentValue == 0 || productValue == 0 {
		return *new(Order), fmt.Errorf("bad request args to new order")
	}

	payment := &Payment{
		Method: paymentMethod,
		Value:  paymentValue,
	}

	product := &Product{
		Category: productCategory,
		Value:    productValue,
	}

	order := &Order{
		Id:      id,
		Product: *product,
		Payment: *payment,
		Labels:  make([]string, 0, 10),
	}

	return *order, nil
}

func Save(pool *pgxpool.Pool, order Order) (uuid.UUID, error) {
	orderRepository := newOrderRepository(pool)
	orderService := newOrderService(orderRepository)

	err := orderService.save(order)

	if err != nil {
		return order.Id, err
	}

	return order.Id, nil
}

func List(pool *pgxpool.Pool) ([]Order, error) {
	orderRepository := newOrderRepository(pool)
	orderService := newOrderService(orderRepository)

	return orderService.list()
}

func Get(pool *pgxpool.Pool, id uuid.UUID) (Order, error) {
	orderRepository := newOrderRepository(pool)
	orderService := newOrderService(orderRepository)

	return orderService.get(id)
}
