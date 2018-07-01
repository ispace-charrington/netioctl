package netif

import (
	"fmt"
	"net"

	"github.com/ispace-charrington/netioctl/ioctl"
	"golang.org/x/sys/unix"
)

// NetIf_Index is a decomposed type, representing one identity of `struct ifreq`
// (which is a C union type - cannot be represented in Go). See `man 7 netdevice`
type NetIf_Index struct {
	Name  [unix.IFNAMSIZ]byte
	Index int32
}

// NetIf_Flags is a decomposed type, representing one identity of `struct ifreq`
// (which is a C union type - cannot be represented in Go). See `man 7 netdevice`
type NetIf_Flags struct {
	Name  [unix.IFNAMSIZ]byte
	Flags int16
}

// NetIf_Sockaddr is a decomposed type, representing one identity of struct ireq
// (which is a C union type - cannot be represented in Go). See man 7 netdevice
type NetIf_Sockaddr struct {
	Name     [unix.IFNAMSIZ]byte
	SockAddr unix.RawSockaddr
}

// NetIf represents a network interface.
type NetIf [unix.IFNAMSIZ]byte

type netifFlags struct {
	ifname [unix.IFNAMSIZ]byte
	flags  int16
}

type netifsockaddr struct {
	ifname [unix.IFNAMSIZ]byte
	fam    int16
	data   [14]byte
}

type netifindex struct {
	ifname [unix.IFNAMSIZ]byte
	index  int32
}

// Up sets the interface flags "UP" and "RUNNING"
func (n NetIf) Up() (err error) {
	s := SocketFd()
	defer SocketClose(s)
	f := &netifFlags{ifname: n}
	err = ioctl.Ioctl(s, unix.SIOCGIFFLAGS, f)
	if err != nil {
		return
	}
	f.flags = f.flags | unix.IFF_UP | unix.IFF_RUNNING
	err = ioctl.Ioctl(s, unix.SIOCSIFFLAGS, f)
	return
}

// Down clears the interface flags "UP" and "RUNNING"
func (n NetIf) Down() (err error) {
	s := SocketFd()
	defer SocketClose(s)
	f := &netifFlags{ifname: n}
	err = ioctl.Ioctl(s, unix.SIOCGIFFLAGS, f)
	if err != nil {
		return
	}
	f.flags = f.flags & ^(unix.IFF_UP | unix.IFF_RUNNING)
	err = ioctl.Ioctl(s, unix.SIOCSIFFLAGS, f)
	return
}

// GetHWAddress returns the MAC of the interface.
func (n NetIf) GetHWAddress() (hwa net.HardwareAddr, err error) {
	s := SocketFd()
	defer SocketClose(s)
	r := &netifsockaddr{ifname: n}

	err = ioctl.Ioctl(s, unix.SIOCGIFHWADDR, r)
	if err != nil {
		return
	}

	switch r.fam {
	case unix.ARPHRD_ETHER:
		hwa = r.data[:6]
	default:
		err = fmt.Errorf("unknown address family 0x%x", r.fam)
	}
	return
}

// SetHWAddress changes the MAC of the interface.
func (n NetIf) SetHWAddress(hwa net.HardwareAddr) (err error) {
	s := SocketFd()
	defer SocketClose(s)
	r := netifsockaddr{ifname: n}
	// TODO - not setting family here, do we need to?
	copy(r.data[:], hwa)

	err = ioctl.Ioctl(s, unix.SIOCSIFHWADDR, r)
	return
}

// Index gets the index of the interface.
func (n NetIf) Index() (idx int32, err error) {
	s := SocketFd()
	defer SocketClose(s)
	r := netifindex{ifname: n}
	err = ioctl.Ioctl(s, unix.SIOCGIFINDEX, &r)
	if err != nil {
		return
	}
	idx = r.index
	return
}

// SocketFd is a convenience method to get a socket file descriptor on which
// to execute ioctls. It's implied that this call cannot fail, and if the OS
// disagrees, we panic. A fd returned in this way must be closed by SocketClose().
func SocketFd() (fd uintptr) {
	n, err := unix.Socket(unix.AF_UNIX, unix.SOCK_DGRAM, 0)
	if err != nil {
		panic(err)
	}
	fd = uintptr(n)
	return
}

// SocketClose is a convenience method to clean up a socket file descriptor
// after use. It's implied that this call cannot fail, and if the OS
// disagrees, we panic.
func SocketClose(fd uintptr) {
	err := unix.Close(int(fd))
	if err != nil {
		panic(err)
	}
}

// GetByName returns a completed NetIf_Index struct for the interface with
// the specified name.
func GetByName(n string) (r *NetIf_Index, err error) {
	r = &NetIf_Index{}
	copy(r.Name[:], n)
	s := SocketFd()
	defer SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCGIFINDEX, r)
	return
}

// GetByIndex returns a completed NetIf_Index struct for the interface at
// the specified index.
func GetByIndex(i int32) (r *NetIf_Index, err error) {
	r = &NetIf_Index{Index: i}
	s := SocketFd()
	defer SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCGIFNAME, r)
	return
}
