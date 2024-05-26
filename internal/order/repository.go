package order

import (
	"context"
	"log/slog"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	pool *pgxpool.Pool
}

func newOrderRepository(pool *pgxpool.Pool) orderRepository {
	return orderRepository{pool}
}

func (or orderRepository) Save(order *Order) error {
	slog.Info("Saving", "order", order.Id)

	return crdbpgx.ExecuteTx(context.Background(), or.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO
				orders
					(
						id,
					 	product_category,
						product_value,
						payment_method,
						payment_value,
						labels
					)
			 VALUES ($1, $2, $3, $4, $5, $6)
			`,
			order.Id,
			order.Product.Category,
			order.Product.Value,
			order.Payment.Method,
			order.Payment.Value,
			order.Labels); err != nil {
			return err
		}
		return nil
	})
}

func (or orderRepository) list() ([]Order, error) {
	rows, err := or.pool.Query(context.Background(), "SELECT * FROM orders")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&order.Id,
			&order.Product.Category,
			&order.Product.Value,
			&order.Payment.Method,
			&order.Payment.Value,
			&order.Labels,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (or orderRepository) get(id uuid.UUID) (Order, error) {
	rows := or.pool.QueryRow(context.Background(), "SELECT * FROM orders WHERE id = $1", id)
	var order Order

	err := rows.Scan(
		&order.Id,
		&order.Product.Category,
		&order.Product.Value,
		&order.Payment.Method,
		&order.Payment.Value,
		&order.Labels,
	)

	if err != nil {
		return order, err
	}

	return order, nil
}
