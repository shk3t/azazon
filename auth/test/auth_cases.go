package authtest

import (
	"base/pkg/model"

	"google.golang.org/grpc/codes"
)

type decodedAuthResponse struct {
	Login string
}

var registerTestCases = []struct {
	request    *model.User
	response   decodedAuthResponse
	statusCode codes.Code
}{
	{
		request: &model.User{
			Login:    "man",
			Password: "somepassword",
		},
		response: decodedAuthResponse{
			Login: "man",
		},
		statusCode: codes.OK,
	},
}

var loginTestCases = []struct {
	request    *model.User
	statusCode codes.Code
}{
	{
		request: &model.User{
			Login:    "man2",
			Password: "somepassword2",
		},
		statusCode: codes.OK,
	},
}

var validateTokenTestCases = []struct {
	request    *model.User
	statusCode codes.Code
}{
	{
		request: &model.User{
			Login:    "man3",
			Password: "somepassword3",
		},
		statusCode: codes.OK,
	},
}