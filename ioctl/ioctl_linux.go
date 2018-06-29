package ioctl

import "golang.org/x/sys/unix"
import "reflect"

// Ioctl serves as a simple convenience function for invoking the
// standard ioctl syscall using the common C convention
func Ioctl(fd, fn uint, d interface{}) (err error) {
	// if nil is passed in, use a NULL pointer
	if d == nil {
		_, _, err = unix.Syscall(
			unix.SYS_IOCTL,
			uintptr(fd),
			uintptr(fn),
			uintptr(0)
		)	
	} else {
		vo := reflect.ValueOf(d)
		if vo.Kind() != reflect.Ptr {
			panic("invalid argument: d must be a pointer")
		}
		_, _, err = unix.Syscall(
			unix.SYS_IOCTL,
			uintptr(fd),
			uintptr(fn),
			vo.Pointer()
	)
	return
}
