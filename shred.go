package shred

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

const (
	numberPasses = 3
)

// Shred securely deletes a file by overwriting its contents with random data multiple times.
//
// Parameters:
//   - path: A string representing the file path to be shredded.
//
// Returns:
//   - An error if the file could not be opened, written to, or if any other issue occurs
//     during the shredding process. Otherwise, it returns nil.
//
// Usage example:
//
//	err := Shred("sensitive_data.txt")
//	if err != nil {
//	    log.Fatalf("Failed to shred the file: %v", err)
//	}
//
// Note:
//   - Effectiveness may vary depending on the filesystem and hardware, such as SSDs
//     with wear leveling, which may not overwrite the data physically on disk.
func Shred(path string) error {
	// Get information of the specified file.
	// - Special case "file doesn't exist" and return directly to the caller
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		} else {
			return fmt.Errorf("failed to get file information: %v", err)
		}
	}

	// Check if the file is a regular file, all other file types are NOT supported.
	if !info.Mode().IsRegular() {
		return fmt.Errorf("unsupported file type")
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	size := info.Size()
	// Note that files of size-0 are not treated as an error,
	// we will simply delete the file.

	for i := 0; i < numberPasses; i++ {
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("failed to seek file: %v", err)
		}

		// Note that CopyN internally will read/write in chunks
		if _, err := io.CopyN(file, rand.Reader, size); err != nil {
			return fmt.Errorf("failed to overwrite file: %v", err)
		}

		// Commits the current contents of the file to stable storage.
		// NB: This isn't guaranteed to write to physical storage, see README
		//     for caveats.
		if err := file.Sync(); err != nil {
			return fmt.Errorf("failed to sync file: %v", err)
		}
	}

	// Close the file before deletion
	err = file.Close()
	file = nil
	if err != nil {
		return fmt.Errorf("failed to close file: %v", err)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
