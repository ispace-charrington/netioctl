package bridge

import "golang.org/x/sys/unix"
import "github.com/ispace-charrington/netioctl/ioctl"
import "github.com/ispace-charrington/netioctl/netif"

// Bridge is an abstract representation of a 802.1d Ethernet Bridge.
type Bridge [unix.IFNAMSIZ]byte

// Create requests the OS create a named bridge interface.
func Create(name string) (b Bridge, err error) {
	var br Bridge
	copy(br[:], name)
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRADDBR, &br)
	if err != nil {
		return
	}
	b = br
	return
}

// NetIf returns a NetIf for the bridge. For example:
//    b := bridge.Create("testbr0")
//    b.NetIf().Up()
func (b Bridge) NetIf() (n netif.NetIf) {
	copy(n[:], b[:])
	return
}

// Destroy takes a Bridge and requests the OS remove it.
func (b Bridge) Destroy() (err error) {
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRDELBR, &b)
	return
}

// Add attempts to add the interface to the bridge.
func (b Bridge) Add(n netif.NetIf) (err error) {
	idx, err := n.Index()
	if err != nil {
		return
	}
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	r := struct {
		n Bridge
		i int32
	}{b, idx}
	err = ioctl.Ioctl(s, unix.SIOCBRADDIF, &r)
	return
}

// Remove attempts to remove the interface to the bridge.
func (b Bridge) Remove(n netif.NetIf) (err error) {
	idx, err := n.Index()
	if err != nil {
		return
	}
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	r := struct {
		n Bridge
		i int32
	}{b, idx}
	err = ioctl.Ioctl(s, unix.SIOCBRDELIF, &r)
	return
}

// AddInterface attempts to add the named interface to the
// bridge.
func (b Bridge) AddInterface(ifname string) (err error) {
	n, err := netif.GetByName(ifname)
	if err != nil {
		return
	}
	n.Name = b
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRADDIF, n)
	return
}

// RemoveInterface attempts to remove the named interface
// from the bridge.
func (b Bridge) RemoveInterface(ifname string) (err error) {
	n, err := netif.GetByName(ifname)
	if err != nil {
		return
	}
	n.Name = b
	s := netif.SocketFd()
	defer netif.SocketClose(s)
	err = ioctl.Ioctl(s, unix.SIOCBRDELIF, n)
	return
}
