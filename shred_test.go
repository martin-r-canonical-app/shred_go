package shred

import (
	"bytes"
	"os"
	"testing"
)

func TestShredDoesNotExist(t *testing.T) {
	err := Shred("")
	if err == nil {
		t.Fatalf("unexpected success on non-exsistant file")
	} else if !os.IsNotExist(err) {
		t.Fatalf("unexpected error on non-exsistant file: %v", err)
	}
}

func TestShredNonRegularFile(t *testing.T) {
	// Create a temporary directory
	tempDirName, err := os.MkdirTemp("", "shred-test-")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.Remove(tempDirName)

	err = Shred(tempDirName)
	if err == nil {
		t.Fatalf("unexpected success on non-regular file")
	}
}

func TestShredReadOnlyPermissions(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "shred-test-")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Close the file to manipulate permissions
	err = tempFile.Close()
	if err != nil {
		t.Fatalf("failed to close temporary file: %v", err)
	}

	// Remove read permissions
	err = os.Chmod(tempFile.Name(), 0400) // Read-only permissions
	if err != nil {
		t.Fatalf("failed to change file permissions: %v", err)
	}

	err = Shred(tempFile.Name())
	if err == nil {
		t.Fatalf("unexpected success on non-writable file")
	}
}

func TestShredEmptyFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "shred-test-")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Close the file before shredding
	err = tempFile.Close()
	if err != nil {
		t.Fatalf("failed to close temporary file: %v", err)
	}

	err = Shred(tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected failure shredding empty file: %v", err)
	}
}

func TestShredCheckRemoved(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "shred-test-")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Close the file before shredding
	err = tempFile.Close()
	if err != nil {
		t.Fatalf("failed to close temporary file: %v", err)
	}

	// Make the file 1KiB so that it's not trivially empty.
	err = os.Truncate(tempFile.Name(), 1024)
	if err != nil {
		t.Fatalf("failed to increase size of temporary file: %v", err)
	}

	err = Shred(tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected failure shredding: %v", err)
	}

	_, err = os.Stat(tempFile.Name())
	if !os.IsNotExist(err) {
		t.Fatalf("shredded file was NOT deleted")
	}
}

func TestShredCheckChanged(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "shred-test-")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Close the file before shredding
	err = tempFile.Close()
	if err != nil {
		t.Fatalf("failed to close temporary file: %v", err)
	}

	// Make the file 1KiB so that it's not trivially empty.
	err = os.Truncate(tempFile.Name(), 1024)
	if err != nil {
		t.Fatalf("failed to increase size of temporary file: %v", err)
	}

	// Get initial contents
	initialData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to read temporary file: %v", err)
	}

	// Create a hard link so we can read the contents after the
	// original tempfile has been deleted.
	hardLinkName := tempFile.Name() + "_hardlink"
	if err := os.Link(tempFile.Name(), hardLinkName); err != nil {
		t.Fatalf("failed to create hard link: %v", err)
	}
	defer os.Remove(hardLinkName)

	err = Shred(tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected failure shredding: %v", err)
	}

	// Get post-shredding contents
	postShredData, err := os.ReadFile(hardLinkName)
	if err != nil {
		t.Fatalf("failed to read temporary file: %v", err)
	}

	// Compare the initial and post-shredded contents
	if bytes.Equal(initialData, postShredData) {
		// Note: There is a (very) unlikely chance the random data
		// would be all zeros, but deemed unlikely enough to
		// not realistically affect test stability.
		t.Fatalf("file not shredded: %v", err)
	}
}
