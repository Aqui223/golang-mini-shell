package main

import (
    "syscall"
    "fmt"
)

func main() {
    var command string = "none"
    var arg_vars []string
    var env_vars []string
    for {
        fmt.Scanln(&command)
        syscall.Exec(command, arg_vars, env_vars)
    }
}
