package main

import (
    "bufio"
    "syscall"
    "fmt"
    "strings"
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
                was_quote := s[i-1] == "\""[0] && s[start] == "\""[0]
                var was_quote_int int
                if was_quote {
                    was_quote_int = 1
                } else {
                    was_quote_int = 0
                }
                fields = append(fields, s[start+was_quote_int:i-was_quote_int])
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
        was_quote := s[len(s)-1] == "\""[0] && s[start] == "\""[0]
        var was_quote_int int
        if was_quote {
            was_quote_int = 1
        } else {
            was_quote_int = 0
        }
        fields = append(fields, s[start+was_quote_int:len(s)-was_quote_int])
    }
    return fields
}

func main() {
    var command string
    var execable_fp string
    var fields []string
    var arg_vars []string
    var env_vars []string
    var bash_history []string
    var status syscall.WaitStatus

    //env_vars = append(env_vars, "A=1")
    file, err := os.Open("envs")
    if err == nil {
        fmt.Println("Loading envs")
    
        defer file.Close()

        scanner := bufio.NewScanner(file)

        for scanner.Scan() {
            text := scanner.Text()
            if strings.Contains(text, "=") {
                env_vars = append(env_vars, text)
            }
        }
    }

    file, err = os.Open("bash_history")
    if err == nil {
        fmt.Println("Loading bash_history")
        defer file.Close()
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            bash_history = append(bash_history, scanner.Text())
        }
    }

    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    exec.Command("stty", "-F", "-echo").Run()

    var score float64
    var best_match int = 0
    var best_match_score float64 = 0

    for {
        command = ""
        if syscall.Getuid() == 0 { // shouldn't even work
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
                fmt.Print("\r$ "+command+"\033[0J")
            } else if b[0] == 9 {
                command = bash_history[best_match]
                fmt.Print("\r$ "+command+"\033[0J")
            } else {
                command += string([]byte{b[0]})
                fmt.Print("\r$ "+command+"\033[0J")
                best_match = -1
                best_match_score = 0

                for history_line_number, history_line := range bash_history {
                    score = 0
                    if len(history_line) <= len(command) {
                        continue
                    }
                    for i, r := range command { // later use binary search
                        if history_line[i] == string(r)[0] {
                            score += 1 / float64(i+1)
                        }
                    }
                    if score > best_match_score {
                        best_match_score = score
                        best_match = history_line_number
                    }
                }
                if best_match != -1 {
                    fmt.Print("\x1b[2m")
                    fmt.Print(bash_history[best_match][len(command):])
                    fmt.Print("\x1b[0m")
                }
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
