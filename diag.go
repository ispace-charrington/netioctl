package main

import "fmt"
import "os"
import "io"
import "github.com/ispace-charrington/netioctl/bridge"
import "github.com/ispace-charrington/netioctl/tuntap"

func must(msg string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "failed to %s: %v\n", msg, err)
	panic(err)
}

func main() {
	b, err := bridge.Create("testbr0")
	must("create a new bridge", err)
	fmt.Fprintf(os.Stderr, "b=%+v\n", b)

	t1, err := tuntap.CreateTap()
	must("create first tap interface", err)
	fmt.Fprintf(os.Stderr, "t1=%+v\n", t1)

	t2, err := tuntap.CreateTap()
	must("create second tap interface", err)
	fmt.Fprintf(os.Stderr, "t2=%+v\n", t2)

	err = b.AddInterface(t1.Name)
	must("add t1 to bridge", err)

	err = b.AddInterface(t2.Name)
	must("add t2 to bridge", err)

	// naive test to see if we can even see data
	io.Copy(os.Stdout, t2)

}
