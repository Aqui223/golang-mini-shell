package main

import (
    //"bufio"
    "syscall"
    "fmt"
//    "strings"
    "os"
    "os/exec"
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
    return fields
}

func main() {
    var command string
    var execable_fp string
    var fields []string
    var arg_vars []string
    var env_vars []string
    var status syscall.WaitStatus

    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    exec.Command("stty", "-F", "-echo").Run()

    for {
        command = ""
        if syscall.Getuid() == 0 {
            fmt.Print("# ")
        } else {
            fmt.Print("$ ")
        }
        var b []byte = make([]byte, 1)
        for {
            os.Stdin.Read(b)
            if b[0] == 10 { // was 13 for some reason
                break
            } else if b[0] == 127 && len(command) != 0 {
                command = command[:len(command)-1]
                fmt.Print("\x1b[3D   \x1b[3D")
            } else {
                command += string([]byte{b[0]})
            }
        }
        fmt.Print("\x1b[2D  \n")
        if command == "exit" {
            return
        }
        fields = Fields(command)
        arg_vars = fields
        execable_fp = fields[0]
        if execable_fp == "cd" {
            syscall.Chdir(fields[1])
            continue
        }
        err := syscall.Access("/bin/"+execable_fp, syscall.F_OK)
        if err == nil {
            execable_fp = "/bin/"+execable_fp
        } else if syscall.Access(execable_fp, syscall.F_OK) != nil {
            fmt.Println("Invalid executable path")
            continue
        }
        id, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
        if id == 0 {
            syscall.Exec(execable_fp, arg_vars, env_vars)
            return;
        } else {
            _, err := syscall.Wait4(int(id), &status, 0, nil);
            if err != nil {
                fmt.Print("Wait4 syscall returned error: ")
                fmt.Print(err)
                fmt.Print("\n")
            }
        }
    }
}
