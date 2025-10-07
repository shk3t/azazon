return {
  debug = {
    go = {
      {
        name = "Auth",
        type = "delve",
        request = "launch",
        program = "auth/cmd/main.go",
      },
      {
        name = "Test auth",
        type = "delve",
        request = "launch",
        mode = "test",
        program = "./auth/test",
      },
      {
        name = "Test notification",
        type = "delve",
        request = "launch",
        mode = "test",
        program = "./notification/test",
      },
    },
  },
}
