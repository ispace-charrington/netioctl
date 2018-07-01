package main

import "fmt"
import "os"
import "io"
import "time"
import "github.com/ispace-charrington/netioctl/bridge"
import "github.com/ispace-charrington/netioctl/tuntap"

func must(msg string, err error) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "ok: %s\n", msg)
		return
	}
	fmt.Fprintf(os.Stderr, "failed to %s: %v\n", msg, err)
	panic(err)
}

func dump(label string, r io.Reader) {
	defer fmt.Fprintf(os.Stderr, "%s dump returning\n", label)

	k := make([]byte, 1550)
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
	defer func() { must("destroy bridge", b.Destroy()) }()

	t1, err := tuntap.CreateTapNamed("tap_t1")
	must("create first tap interface", err)
	fmt.Fprintf(os.Stderr, "t1=%+v\n", t1)
	defer func() { must("close t1", t1.Close()) }()

	t2, err := tuntap.CreateTapNamed("tap_t2")
	must("create second tap interface", err)
	fmt.Fprintf(os.Stderr, "t2=%+v\n", t2)
	defer func() { must("close t2", t2.Close()) }()

	err = b.Add(t1.NetIf())
	must("add t1 to bridge", err)
	defer func() { must("remove t1 from bridge", b.Remove(t1.NetIf())) }()

	err = b.Add(t2.NetIf())
	must("add t2 to bridge", err)
	defer func() { must("remove t2 from bridge", b.Remove(t2.NetIf())) }()

	err = t1.NetIf().Up()
	must("bring up interface t1", err)
	defer func() { must("shut down t1", t1.NetIf().Down()) }()

	err = t2.NetIf().Up()
	must("bring up interface t2", err)
	defer func() { must("shut down t2", t2.NetIf().Down()) }()

	err = b.NetIf().Up()
	must("bring up bridge interface", err)
	defer func() { must("shut down bridge", b.NetIf().Down()) }()

	go dump("t1", t1)
	go dump("t2", t2)
	time.Sleep(30 * time.Second)

}
