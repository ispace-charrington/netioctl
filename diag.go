package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ispace-charrington/netioctl/bridge"
	"github.com/ispace-charrington/netioctl/tuntap"
)

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

func abort(a *chan bool) {
	if *a == nil {
		return
	}
	close(*a)
	*a = nil
}

func main() {
	b, err := bridge.Create("testbr0")
	must("create a new bridge", err)
	defer func() { must("destroy bridge", b.Destroy()) }()

	t1, err := tuntap.CreateTapNamed("tap_t1")
	must("create first tap interface", err)
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

	//go dump("t1", t1)
	//go dump("t2", t2)
	//time.Sleep(30 * time.Second)

	abortChan := make(chan bool)
	go func() {
		// reader thread that watches for a packet sent on another thread
		defer abort(&abortChan)
		k := make([]byte, 1550)
		for {
			n, err := t2.Read(k)
			if err != nil {
				fmt.Fprintf(os.Stderr, "t2 read error %v\n", err)
				return
			}
			if n == 0 {
				continue
			}
			if n < 14 {
				fmt.Fprintf(os.Stderr, "t2 short read len %d\n", n)
				return
			}
			if bytes.Equal(k[12:14], []byte{0xC0, 0xC0}) {
				fmt.Fprintf(os.Stderr, "t2 successfully read signal frame (received from 0x%012x)\n", k[6:12])
				return
			}

		}
	}()
	go func() {
		// writer thread that sends a frame with a specific ethertype
		defer abort(&abortChan)
		// frame header
		k := []byte{
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // destination
			0xde, 0xad, 0xde, 0xad, 0xde, 0xad, // source
			0xc0, 0xc0, // ethertype
		}
		// 64 byte payload
		k = append(k, "payload payload payload payload payload payload payload payload."...)
		n, err := t1.Write(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "t1 write error %v\n", err)
			return
		}
		if n < len(k) {
			fmt.Fprintf(os.Stderr, "t1 short write (%d < %d)\n", n, len(k))
			return
		}
		fmt.Fprintln(os.Stderr, "t1 write success, waiting on abort channel")
		<-abortChan
	}()

	select {
	case <-abortChan:
		fmt.Fprintln(os.Stderr, "main thread abort")
		return
	case <-time.Tick(10 * time.Second):
		fmt.Fprintln(os.Stderr, "main thread timeout")
		return
	}
}
