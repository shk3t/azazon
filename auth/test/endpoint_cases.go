package authtest

import (
	"auth/internal/model"

	"google.golang.org/grpc/codes"
)

type decodedAuthResponse struct {
	Login string
}

var registerTestCases = []struct {
	request    model.User
	response   decodedAuthResponse
	statusCode codes.Code
}{
	{
		request: model.User{
			Login:    "man",
			Password: "somepassword",
		},
		response: decodedAuthResponse{
			Login: "man",
		},
		statusCode: codes.OK,
	},
	{
		request: model.User{
			Login:    "shortman",
			Password: "short",
		},
		response:   decodedAuthResponse{},
		statusCode: codes.InvalidArgument,
	},
}

var loginTestCases = []struct {
	request    model.User
	statusCode codes.Code
}{
	{
		request: model.User{
			Login:    "man2",
			Password: "somepassword2",
		},
		statusCode: codes.OK,
	},
}

var validateTokenTestCases = []struct {
	registerRequest model.User
	statusCode      codes.Code
}{
	{
		registerRequest: model.User{
			Login:    "man3",
			Password: "somepassword3",
		},
		statusCode: codes.OK,
	},
}

var updateUserTestCases = []struct {
	oldUser    model.User
	newUser    model.User
	statusCode codes.Code
}{
	{
		oldUser: model.User{
			Login:    "man4",
			Password: "somepassword4",
		},
		newUser: model.User{
			Login:    "newman4",
			Password: "newpassword4",
			Role:     model.UserRoles.Admin,
		},
		statusCode: codes.OK,
	},
}