package query

import (
	db "auth/internal/database"
	"auth/internal/model"
	"context"
)

const userBaseSelectQuery = "SELECT * FROM \"user\" "

func GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	user := model.User{}
	err := db.ConnPool.QueryRow(
		ctx, userBaseSelectQuery+"WHERE login = $1", login,
	).Scan(&user.Id, &user.Login, &user.Password)
	return user, err
}

func CreateUser(ctx context.Context, u model.User) (int, error) {
	var id int
	err := db.ConnPool.QueryRow(
		ctx, `
        INSERT INTO "user" (login, password)
        VALUES ($1, $2)
        RETURNING id`,
		u.Login, u.Password,
	).Scan(&id)
	return id, err
}

func UpdateUser(ctx context.Context, id int, u model.User) error {
	_, err := db.ConnPool.Exec(
		ctx,
		"UPDATE \"user\" SET login = $1, password = $2 WHERE id = $3",
		u.Login, u.Password,
		id,
	)
	return err
}

func DeleteUser(ctx context.Context, id int) {
	db.ConnPool.Exec(ctx, "DELETE FROM user WHERE id = $1", id)
}