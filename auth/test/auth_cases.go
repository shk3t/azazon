package authtest

import (
	"auth/internal/model"
)

var registerTestCases = []struct {
	payload    model.User
	response   model.AuthResponse
	statusCode int
}{
	{
		payload: model.User{
			Login:    "man",
			Password: "somepassword",
		},
		response: model.AuthResponse{
			User: &model.User{
				Login: "man",
			},
		},
		statusCode: 200,
	},
}