package query

import (
	db "auth/internal/database"
	"common/pkg/model"
	"context"
)

func GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	user := model.User{}
	err := db.ConnPool.QueryRow(
		ctx, `
		SELECT id, login, password_hash, role
		FROM "user"
		WHERE login = $1`,
		login,
	).Scan(&user.Id, &user.Login, &user.PasswordHash, &user.Role)
	return user, err
}

func CreateUser(ctx context.Context, u model.User) (int, error) {
	var id int
	err := db.ConnPool.QueryRow(
		ctx, `
        INSERT INTO "user" (login, password_hash, role)
        VALUES ($1, $2, $3)
        RETURNING id`,
		u.Login, u.PasswordHash, u.Role,
	).Scan(&id)
	return id, err
}

func UpdateUser(ctx context.Context, id int, u model.User) error {
	_, err := db.ConnPool.Exec(
		ctx, `
		UPDATE "user"
		SET login = $1, password_hash = $2, role = $3
		WHERE id = $3`,
		u.Login, u.PasswordHash, u.Role, id,
	)
	return err
}

func DeleteUser(ctx context.Context, id int) {
	db.ConnPool.Exec(ctx, "DELETE FROM user WHERE id = $1", id)
}