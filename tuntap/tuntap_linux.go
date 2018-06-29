package tuntap

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"github.com/ispace-charrington/netioctl/ioctl"
	"github.com/ispace-charrington/netioctl/netif"
	"golang.org/x/sys/unix"
)

// TapIf represents a `tap` interface, which exchanges behaves like
// a full Ethernet-capable network interface to userspace, but also
// permits that network interface to be a io.ReadWriteCloser.
type TapIf struct {
	Name string
	fp   *os.File
}

func tuntapdev() (*os.File, error) {
	return os.OpenFile("/dev/net/tun", O_RDWR, 0600)
}

func createTap(n *netif.NetIf_Flags) (t *TapIf, err error) {
	// https://github.com/torvalds/linux/blob/fd3a88625844907151737fc3b4201676effa6d27/drivers/net/tap.c#L992
	fd, err := tuntapdev()
	if err != nil {
		return
	}
	err = ioctl.Ioctl(fd.Fd(), unix.TUNSETIFF, r)
	if err != nil {
		fd.Close() // don't leak a device if we can't configure it
		return
	}
	t = &TapIf{fp: fd}
	t.Name = r.Name[:bytes.IndexByte(r.Name, 0)]
	return
}

// CreateTap requests a new automatically named `tap` interface from
// the OS. This device is created with the NO_PI flag set, because
// essentially no sane users are interested in the alternative.
func CreateTap() (*TapIf, error) {
	r := &netif.NetIf_Flags{Flags: unix.IFF_NO_PI | unix.IFF_TAP}
	return createTap(r)
}

// CreateTapNamed requests a new `tap` interface from the OS, and
// requests that it be named with the provided string. It is otherwise
// identical to CreateTap(). The requested name may be up to IFNAMSIZ
// bytes, which technically can vary, but seems to be 16 everywhere.
func CreateTapNamed(name string) (*TapIf, error) {
	r := &netif.NetIf_Flags{Flags: unix.IFF_NO_PI | unix.IFF_TAP}
	if len(name) > unix.IFNAMSIZ {
		return nil, fmt.Errorf("'%s' is longer than maximum of %d", unix.IFNAMSIZ)
	}
	copy(r.Name[:], name)
	return createTap(r)
}

func (t *TapIf) GetHWAddress() (*net.HardwareAddr, error) {
	// https://golang.org/pkg/net/#HardwareAddr
	// https://github.com/torvalds/linux/blob/fd3a88625844907151737fc3b4201676effa6d27/drivers/net/tap.c#L1091
	// stub
}

func (t *TapIf) SetHWAddress(a *net.HardwareAddr) error {
	// https://github.com/torvalds/linux/blob/fd3a88625844907151737fc3b4201676effa6d27/drivers/net/tap.c#L1108
	// stub
}

// https://golang.org/pkg/io/#ReadWriteCloser ...

func (t *TapIf) Read(p []byte) (n int, err error) {
	// stub
}

func (t *TapIf) Write(p []byte) (n int, err error) {
	// stub
}

func (t *TapIf) Close() (err error) {
	// stub
}
