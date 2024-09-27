package srm

import (
	"os"
	"syscall"
)

func Open(path string) (int, error) {
	fd, err := syscall.Open(path, os.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	return fd, nil
}

func Read(path string) ([]byte, error) {
	f, err := Open(path)
	fd := os.NewFile(uintptr(f), path)
	if err != nil {
		return nil, err
	}
	fileState, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileState.Size())

	fd.Read(buffer)
	return buffer, nil
}
