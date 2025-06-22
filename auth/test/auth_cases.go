package authtest

import (
	m "auth/internal/model"
)

var registerTestCases = []struct {
	payload      m.User
	response   m.AuthResponse
	statusCode int
}{
	{
		payload: m.User{
			Login:    "man",
			Password: "somepassword",
		},
		response: m.AuthResponse{
			User: &m.User{
				Login: "man",
			},
		},
		statusCode: 200,
	},
}