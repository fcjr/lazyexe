# LazyExe

Lazily write an executable to disk.  Useful for embedding executables via go's new `//go:embed` directive.

A new LazyExe can be created by running `lazyexe.New(exeBytes)`. Don't call `Path()` until you need to run the embedded executable.

The executable is written out to your tmp folder and chmoded the first time the `Path()` is requested. Subsequent calls to `Path()`
will return the same executable.  When finished using the executable you are expected to call `Cleanup()` to ensure the temporary
file is removed from the system.

## Example

Try running `example/main.go` which embeds a [cosmopolitan libc](https://justine.lol/cosmopolitan/index.html) hello world, loads it, and runs it!

```go
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
```