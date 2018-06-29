package netif

import "golang.org/x/sys/unix"
import "github.com/ispace-charrington/netioctl/ioctl"

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

// SocketFd is a convenience method to get a socket file descriptor on which
// to execute ioctls. It's implied that this call cannot fail, and if the OS
// disagrees, we panic. A fd returned in this way must be closed by SocketClose().
func SocketFd() (fd int) {
	fd, err := unix.Socket(unix.AF_UNIX, unix.SOCK_DGRAM, 0)
	if err != nil {
		panic(err)
	}
	return
}

// SocketClose is a convenience method to clean up a socket file descriptor
// after use. It's implied that this call cannot fail, and if the OS
// disagrees, we panic.
func SocketClose(fd int) {
	err := unix.Close(fd)
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
