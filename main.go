//GO REPL is a simple application promising to write/compile/run code in terminal, inspired by python shell
//Author: Saddam H
//Email: thedevsaddam@gmail.com
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	APP_FULL_NAME    = "GO REPL"
	APP_SHORT_NAME   = "gorepl"
	APP_VERSION      = "1.0.0"
	EXT              = ".go" //go code file extension
	FILE_NAME_LENGTH = 10    //temporary file name length
	TIME_LAYOUT      = "Mon, 01/02/06, 03:04 PM"
)

func main() {
	//message to start writing code
	help := "Use :cr + enter to compile and run the code"
	fmt.Fprintf(os.Stdout, "%s (version %s) | %s | %s\n", APP_FULL_NAME, APP_VERSION, help, time.Now().Format(TIME_LAYOUT))
	fmt.Fprintln(os.Stdout, "Start writing your 'GO' code:\n  1|> package main")
	//read user codes
	body := readLines()
	//write the codes to a temp file for executing
	file, err := writeBody(randString(FILE_NAME_LENGTH)+EXT, body)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s ERROR: %s\n", APP_FULL_NAME, err.Error())
		os.Exit(1)
	}
	//run the temp file
	err = run(file)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s ERROR: %s\n", APP_FULL_NAME, err.Error())
		os.Exit(1)
	}
	//remove the temp file
	os.Remove(file)
}

//writeBody write the whole code body into a temp go file to execute
func writeBody(fileName, line string) (string, error) {
	//make the temp directory if not exist
	dir, err := ioutil.TempDir("", APP_SHORT_NAME)
	if err != nil {
		return fileName, err
	}
	//write the file inside the HOME/directory
	tempGoFile := filepath.Join(dir, fileName) + EXT
	return tempGoFile, ioutil.WriteFile(tempGoFile, []byte(line), 0644)
}

//readLine reads to standard input until it get the signal (:cr)
func readLines() string {
	body := "package main\n"
	reader := bufio.NewReader(os.Stdin)
	lineNo := 2
	for {
		//make a fancy look with line number
		if lineNo > 9 && lineNo < 100 {
			fmt.Fprintf(os.Stdout, " %d|> ", lineNo)
		} else if lineNo > 99 {
			fmt.Fprintf(os.Stdout, "%d|> ", lineNo)
		} else {
			fmt.Fprintf(os.Stdout, "  %d|> ", lineNo)
		}

		text, _ := reader.ReadString('\n')
		if strings.TrimRight(text, "\n") == ":cr" {
			return body
		}
		body += "\n" + text
		lineNo++
	}
	return body
}

//run execute the temp go code and output to the standard output
func run(file string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("go", "run", file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

//randString a random string with length
func randString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
