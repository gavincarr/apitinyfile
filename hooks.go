package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// checkPostHook checks that path exists and is executable
func checkPostHook(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("post-hook does not exist?")
		} else {
			return err
		}
	} else if stat.IsDir() {
		return errors.New("post-hook is directory?")
	}

	// Check that at least one executable bit is set
	perm := stat.Mode().Perm()
	if perm&0111 == 0 {
		return errors.New("post-hook is not executable?")
	}

	return nil
}

// postHook
func postHook(path, method, filepath string, verbose bool) {
	os.Setenv("HOOK_METHOD", method)
	os.Setenv("HOOK_PATH", filepath)
	cmd := exec.Command(path)
	out, err := cmd.CombinedOutput()
	if verbose {
		if err != nil {
			fmt.Println("post-hook error: ", err.Error())
		}
		if len(out) != 0 {
			fmt.Println("post-hook output:")
			fmt.Print(string(out))
		}
	}
}
