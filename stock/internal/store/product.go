package store

import (
	errpkg "common/pkg/errors"
	"context"
	"errors"
	"stock/internal/model"
	"stock/internal/query"

	"github.com/jackc/pgx/v5"
)

type PostgreProductStore struct {
	HashPassword func(password string) (string, error)
}

func (s *PostgreProductStore) Get(ctx context.Context, id int) (model.Product, error) {
	product, err := query.GetProductById(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return product, errpkg.NotFound
	}

	return product, err
}

func (s *PostgreProductStore) Save(
	ctx context.Context,
	product model.Product,
) (model.Product, error) {
	var err error

	if product.Id == 0 {
		product.Id, err = query.CreateProduct(ctx, nil, product)
	} else {
		err = query.UpdateProduct(ctx, nil, product.Id, product)
	}

	return product, err
}

func (s *PostgreProductStore) Delete(ctx context.Context, id int) error {
	query.DeleteProduct(ctx, id)
	return nil
}