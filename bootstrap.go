package main

import (
	"os"
	"os/exec"
)

func main() {
	c := exec.Command("bash", "./bootstrap.sh")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}
