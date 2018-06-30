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
func DestroyByName(name string) (err error) {
	var brnam [unix.IFNAMSIZ]byte
	copy(brnam[:], name)
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRDELBR, &brnam)
	return
}

// Create requests the OS create a named bridge interface.
func Create(name string) (b *Bridge, err error) {
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

// AddInterface attempts to add the named interface to the
// bridge.
func (b *Bridge) AddInterface(ifname string) (err error) {
	n, err := netif.GetByName(ifname)
	if err != nil {
		return
	}
	copy(n.Name[:], b.Name)
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRADDIF, n)
	return
}

// RemoveInterface attempts to remove the named interface
// from the bridge.
func (b *Bridge) RemoveInterface(ifname string) (err error) {
	n, err := netif.GetByName(ifname)
	if err != nil {
		return
	}
	copy(n.Name[:], b.Name)
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRDELIF, n)
	return
}
