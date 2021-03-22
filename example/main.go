package main

import (
	_ "embed"
	"fmt"
	"log"
	"os/exec"

	"github.com/fcjr/lazyexe"
)

//go:embed hello.com
var helloworld []byte

func main() {
	exe := lazyexe.New(helloworld)
	defer exe.Cleanup()

	exePath, err := exe.Path()
	if err != nil {
		log.Fatal(err)
	}

	out, err := exec.Command(exePath).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(out))
}
