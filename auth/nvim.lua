return {
  debug = {
    go = {
      {
        type = "delve",
        request = "launch",
        name = "Default",
        program = "cmd/authmain.go",
      },
    },
  },
}
