# Shred tool in go

## Provided specification

Shred tool in Go.

Implement a Shred(path) function that will overwrite the given file (e.g. “randomfile”) 3 times with random data and delete the file afterwards. Note that the file may contain any type of data.

You are expected to give information about the possible test cases for your Shred function, including the ones that  you don’t implement, and implementing the full test coverage is a bonus :)

In a few lines briefly discuss the possible use cases for such a helper function as well as advantages and drawbacks of addressing them with this approach.

## Prior art

- GNU coreutils implementation https://github.com/wertarbyte/coreutils/blob/master/src/shred.c 
- https://en.wikipedia.org/wiki/Data_erasure, including various standards

The approaches here are significantly more complex than a basic implementation, and offers some insight into both what can be easily added as well as what is missing.

## Approach

- Basic parameter checking for file existance, permissions, and file type (only supporting regular files).
- The file is opened once for all 3 passes, but operations are synced to stable storage (using fsync under the covers) after each pass.
- cryto/rand is used as a cryptographically secure random number generator for the contents on each pass.


The following features are deemed out of scope for this exercise:

- Writing to the file is chunked under the covers, as the chunk size performance will have a large variance between storage media there is little advantage trying to overoptimize this without a specific use case.
- Only support regular files and not symbolic links + testing lots of variants.
- Creating a wrapping CLI executable
- Alternative function APIs, e.g. accept a file pointer
- Using a more formal repository structure, e.g. https://github.com/golang-standards/project-layout
- Developer infrastructure e.g. for integrating tests with CI
- Specific types of random data, e.g. to conform to DoD 5220.22-M standard
- Progress indicator
- Using custom hardware API calls for erasure

## Testing

Test cases implemented cover:

- Error case: File does not exist
- Error case: Non-regular file (directory)
- Error case: Read-only permissions
- Edge case success: Empty file
- Mainline: Checking file is removed
- Mainline: Checking file is overwritten

Further test cases could be added to include:

- Failure cases for os functions (e.g. `Seek`, `Sync`, `Close`, `Remove`)
- Testing on different file sizes
- Check the file is overwritten three times
- Check all of the file is written
- Performance testing on different hardware

## Use cases, pros, and cons

The expected use case of the Shred(path) function is security, so that data contents are non-recoverable from the underlying software of storage medium.

The main pros of using this approach are:
- It offers a (relatively) quick way of (somewhat) securely erasing data with out needing to destroy the physcial machine.

Main cons of using for this scenario are:
- Shredding the file does NOT guarantee erasure of the original file, and can happen for a variety of reasons, including but not limited to:
  - SSDs with wear leveling, which may not overwrite the data physically on disk
  - Copy-on-write filesystems
  - Journaled filesystems
- Shredding does NOT erase file metadata
  - Basic attributes such as timestamps could be used to derive information from other data sources
  - Extended attributes could contain more sensitive information
- Performance
  - Shredding involves multiple overwrites, which can be time-consuming for large files.

# Usage
```
package main

import (
    "github.com/martin-r-canonical-app/shred_go"
    "log"
    )

func main() {
    err := shred.Shred("randomfile")
    if err != nil {
        log.Fatalf("Failed to shred the file: %v", err)
    }
}
```

# Installation
```
go get -u github.com/martin-r-canonical-app/shred_go
```
