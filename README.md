# golang-mini-shell
Minimalistic shell written in golang.

## Run
`go run main.go`

## Usage
Currently this can only run commands in /bin/ and whatever executables you enter, it also supports cd but doesn't support environment variables and argv input isn't finished.

## Uses
Interacting with BASH can cause security vulnerabilities because it is a big project and often it would be hard to understand all the parsing that it is doing, also this project is way easier to modify.
For an example the $PATH environment variable can be edited so that any command you run is a malicious one, potentially swapping sudo for something else, also there's .bashrc that might cause trouble if you break it.
