package tuntap

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/ispace-charrington/netioctl/ioctl"
	"github.com/ispace-charrington/netioctl/netif"
	"golang.org/x/sys/unix"
)

// TapIf represents a tap interface, which exchanges behaves like
// a full Ethernet-capable network interface to userspace, but also
// permits that network interface to be a io.ReadWriteCloser.
type TapIf struct {
	Name       string
	persistent bool
	fp         *os.File
}

func tuntapdev() (*os.File, error) {
	return os.OpenFile("/dev/net/tun", unix.O_RDWR, 0600)
}

func tapdevDetectLeak(t *TapIf) {
	if t.fp == nil {
		return // already closed 👍
	}
	t.fp.Close()

	if t.persistent {
		// A persistent interface doesn't need to be explicitly closed
		// since they'll continue to exist anyway, so don't crash, but
		// do ridicule and shame
		return
	}
	panic(fmt.Errorf("non-persistent TapIf leaked (Name=%q)", t.Name))
}

func createTap(n *netif.NetIf_Flags) (t *TapIf, err error) {
	// https://github.com/torvalds/linux/blob/fd3a88625844907151737fc3b4201676effa6d27/drivers/net/tap.c#L992
	fd, err := tuntapdev()
	if err != nil {
		return
	}
	err = ioctl.Ioctl(int(fd.Fd()), unix.TUNSETIFF, n)
	if err != nil {
		fd.Close() // don't leak a device if we can't configure it
		return
	}
	t = &TapIf{fp: fd}
	t.Name = string(n.Name[:bytes.IndexByte(n.Name[:], 0)])
	runtime.SetFinalizer(t, tapdevDetectLeak)
	return
}

// CreateTap requests a new automatically named tap interface from
// the OS. This device is created with the NO_PI flag set, because
// essentially no sane users are interested in the alternative. A
// TapIf must be Close()d before it is GC'd or we will panic.
func CreateTap() (*TapIf, error) {
	r := &netif.NetIf_Flags{Flags: unix.IFF_NO_PI | unix.IFF_TAP}
	return createTap(r)
}

// CreateTapNamed requests a new tap interface from the OS, and
// requests that it be named with the provided string. It is otherwise
// identical to CreateTap(). The requested name may be up to IFNAMSIZ
// bytes, which technically can vary, but seems to be 16 everywhere.
// A TapIf must be Close()d before it is GC'd or we will panic.
func CreateTapNamed(name string) (*TapIf, error) {
	r := &netif.NetIf_Flags{Flags: unix.IFF_NO_PI | unix.IFF_TAP}
	if len(name) > unix.IFNAMSIZ {
		return nil, fmt.Errorf("'%s' is longer than maximum of %d", name, unix.IFNAMSIZ)
	}
	copy(r.Name[:], name)
	return createTap(r)
}

// GetHWAddress returns the MAC of the tap interface.
func (t *TapIf) GetHWAddress() (*net.HardwareAddr, error) {
	// https://golang.org/pkg/net/#HardwareAddr
	// https://github.com/torvalds/linux/blob/fd3a88625844907151737fc3b4201676effa6d27/drivers/net/tap.c#L1091
	// stub
	return nil, nil
}

// SetHWAddress changes the MAC of the tap interface.
func (t *TapIf) SetHWAddress(a *net.HardwareAddr) error {
	// https://github.com/torvalds/linux/blob/fd3a88625844907151737fc3b4201676effa6d27/drivers/net/tap.c#L1108
	// stub
	return nil
}

// https://golang.org/pkg/io/#ReadWriteCloser ...

// Read reads ethernet frames that were "transmitted" on this
// tap interface.
func (t *TapIf) Read(p []byte) (n int, err error) {
	// stub
	return
}

// Write queues ethernet frames to be "received" on this tap
// interface.
func (t *TapIf) Write(p []byte) (n int, err error) {
	// stub
	return
}

// Close disposes the tap interface and frees any resources.
func (t *TapIf) Close() (err error) {
	// stub
	return
}
