package main

import (
    "bufio"
    "syscall"
    "fmt"
    "strings"
    "os"
)

func main() {
    var command string
    var execable_fp string
    var fields []string
    var arg_vars []string
    var env_vars []string

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("~$ ")
        if !scanner.Scan() {
            break
        }
        command = scanner.Text()
        fields = strings.Fields(command)
        arg_vars = fields[0:]
        execable_fp = fields[0]
        fmt.Println(arg_vars)

        syscall.Exec(execable_fp, arg_vars, env_vars)
    }
}
