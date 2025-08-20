return {
  debug = {
    go = {
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
