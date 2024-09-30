package utils

import (
	"fmt"
	"os"
)

func CreateLogs(port string) error {
	file, err := os.Create("netcat-connection_" + port + ".log")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = file.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	file, err = os.Create("netcat-chat_" + port + ".log")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = file.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
