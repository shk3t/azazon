package conversion

import (
	"auth/internal/model"
	"common/api/auth"
	"common/pkg/sugar"
)

func UserModel[R *auth.RegisterRequest | *auth.LoginRequest | *auth.UpdateUserRequest](
	r R,
) *model.User {
	switch r := any(r).(type) {
	case *auth.RegisterRequest:
		return &model.User{
			Login:    r.Login,
			Password: r.Password,
		}
	case *auth.LoginRequest:
		return &model.User{
			Login:    r.Login,
			Password: r.Password,
		}
	case *auth.UpdateUserRequest:
		return &model.User{
			Login:    sugar.Value(r.NewLogin),
			Password: sugar.Value(r.NewPassword),
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
