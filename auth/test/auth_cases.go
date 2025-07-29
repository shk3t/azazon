package authtest

import (
	"base/api/auth"

	"google.golang.org/grpc/codes"
)

type decodedRegisterResponse struct {
	Login string
}

var registerTestCases = []struct {
	payload    *auth.RegisterRequest
	response   decodedRegisterResponse
	statusCode codes.Code
}{
	{
		payload: &auth.RegisterRequest{
			Login:    "man",
			Password: "somepassword",
		},
		response: decodedRegisterResponse{
			Login: "man",
		},
		statusCode: codes.OK,
	},
}