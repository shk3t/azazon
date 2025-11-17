package conversion

import (
	"auth/internal/model"
	"auth/internal/query"
	"common/api/auth"
	"common/pkg/sugar"
)

func UserModel[T *query.User | *auth.RegisterRequest | *auth.LoginRequest | *auth.UpdateUserRequest](
	in T,
) *model.User {
	switch in := any(in).(type) {
	case *query.User:
		return &model.User{
			Id:           int(in.ID),
			Login:        in.Login,
			PasswordHash: in.PasswordHash,
			Role:         model.UserRole(in.Role),
		}
	case *auth.RegisterRequest:
		return &model.User{
			Login:    in.Login,
			Password: in.Password,
		}
	case *auth.LoginRequest:
		return &model.User{
			Login:    in.Login,
			Password: in.Password,
		}
	case *auth.UpdateUserRequest:
		return &model.User{
			Login:    sugar.Value(in.NewLogin),
			Password: sugar.Value(in.NewPassword),
		}
	}
	return nil
}

func RegisterRequest(u *model.User) *auth.RegisterRequest {
	return &auth.RegisterRequest{
		Login:    u.Login,
		Password: u.Password,
	}
}

func LoginRequest(u *model.User) *auth.LoginRequest {
	return &auth.LoginRequest{
		Login:    u.Login,
		Password: u.Password,
	}
}