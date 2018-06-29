package bridge

import "golang.org/x/sys/unix"
import "github.com/ispace-charrington/netioctl/ioctl"
import "github.com/ispace-charrington/netioctl/netif"

// Bridge is an abstract representation of a 802.1d Ethernet Bridge.
type Bridge struct {
	Name string
}

// Destroy takes a Bridge and requests the OS remove it.
func (b *Bridge) Destroy() error {
	return DestroyByName(b.Name)
}

// DestroyByName requests the OS remove a named bridge interface.
func DestroyByName(name String) (err error) {
	var brnam [unix.IFNAMSIZ]byte
	copy(brnam[:], name)
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRDELBR, &brnam)
	return
}

// Create requests the OS create a named bridge interface.
func Create(name String) (b *Bridge, err error) {
	var brnam [unix.IFNAMSIZ]byte
	copy(brnam[:], name)
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRADDBR, &brnam)
	if err != nil {
		return
	}
	b = &Bridge{Name: name}
	return
}

func (b *Bridge) AddInterface(ifname String) (err error) {
	// stub
}

func (b *Bridge) RemoveInterface(ifname String) (err error) {
	// stub
}
