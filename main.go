package main

import (
    "bufio"
    "syscall"
    "fmt"
//    "strings"
    "os"
    "unicode"
)

func Fields(s string) []string {
    var fields []string

    start := 0
    var in_quote bool = false

    for i, r := range s {
        if unicode.IsSpace(r) {
            if (start < i) && (!in_quote) {
                fields = append(fields, s[start:i])
            }
            if !in_quote {
                start = i+1
            }
        }
        if r == '"' {
            in_quote = !in_quote
        }
    }
    if start != len(s) {
        fields = append(fields, s[start:])
    }
    fmt.Println(len(fields))
    return fields
}

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
        fields = Fields(command)
        arg_vars = fields
        execable_fp = fields[0]
        fmt.Println(arg_vars)

        syscall.Exec(execable_fp, arg_vars, env_vars)
    }
}
