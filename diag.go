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

func dump(label string, r io.Reader, abort chan bool) {
	defer func() { abort <- true }()
	k := make([]byte, 0, 1550)
	for {
		n, err := r.Read(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s.Read() returned error %v\n", label, err)
			return
		}
		if n == 0 {
			continue
		}
		if n < 14 {
			fmt.Fprintf(os.Stderr, "%s.Read() returned less than 14 bytes (%d)??\n", label, n)
			continue
		} else {
			fmt.Printf("%s %012x > %012x %04x\n", label, k[6:12], k[0:6], k[12:14])
		}
	}
}

func main() {
	b, err := bridge.Create("testbr0")
	must("create a new bridge", err)
	fmt.Fprintf(os.Stderr, "b=%+v\n", b)

	t1, err := tuntap.CreateTap()
	must("create first tap interface", err)
	fmt.Fprintf(os.Stderr, "t1=%+v\n", t1)
	defer t1.Close()

	t2, err := tuntap.CreateTap()
	must("create second tap interface", err)
	fmt.Fprintf(os.Stderr, "t2=%+v\n", t2)
	defer t2.Close()

	err = b.Add(t1.NetIf())
	must("add t1 to bridge", err)

	err = b.Add(t2.NetIf())
	must("add t2 to bridge", err)

	err = t1.NetIf().Up()
	must("bring up interface t1", err)

	err = t2.NetIf().Up()
	must("bring up interface t2", err)

	err = b.NetIf().Up()
	must("bring up bridge interface", err)

	abort := make(chan bool)
	go dump("t1", t1, abort)
	go dump("t2", t2, abort)
	<-abort
}
