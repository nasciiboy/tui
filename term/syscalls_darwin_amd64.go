// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs syscalls.go

package term

type syscall_Termios struct {
	Iflag     uint64
	Oflag     uint64
	Cflag     uint64
	Lflag     uint64
	Cc        [20]uint8
	Pad_cgo_0 [4]byte
	Ispeed    uint64
	Ospeed    uint64
}

const (
	syscall_IGNBRK = 0x1
	syscall_BRKINT = 0x2
	syscall_PARMRK = 0x8
	syscall_ISTRIP = 0x20
	syscall_INLCR  = 0x40
	syscall_IGNCR  = 0x80
	syscall_ICRNL  = 0x100
	syscall_IXON   = 0x200
	syscall_OPOST  = 0x1
	syscall_ECHO   = 0x8
	syscall_ECHONL = 0x10
	syscall_ICANON = 0x100
	syscall_ISIG   = 0x80
	syscall_IEXTEN = 0x400
	syscall_CSIZE  = 0x300
	syscall_PARENB = 0x1000
	syscall_CS8    = 0x300
	syscall_VMIN   = 0x10
	syscall_VTIME  = 0x11

	syscall_TCGETS = 0x40487413
	syscall_TCSETS = 0x80487414
)
