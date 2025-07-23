package utils

import "os"

// EnsureDetinationPAth is responsible for ensuring a path that is to be written
// to exists.
func EnsureDestinationPath(dest string) error {
	return os.MkdirAll(dest, os.ModePerm)
}