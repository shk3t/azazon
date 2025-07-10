package authtest

import (
	"auth/internal/model"

	"google.golang.org/grpc/codes"
)

var registerTestCases = []struct {
	payload    model.User
	response   model.AuthResponse
	statusCode codes.Code
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
		statusCode: codes.OK,
	},
}