package authtest

import (
	"base/api/auth"

	"google.golang.org/grpc/codes"
)

var registerTestCases = []struct {
	payload    *auth.User
	response   *auth.AuthResponse
	statusCode codes.Code
}{
	{
		payload: &auth.User{
			Login:    "man",
			Password: "somepassword",
		},
		response: &auth.AuthResponse{
			User: &auth.User{
				Login: "man",
			},
		},
		statusCode: codes.OK,
	},
}