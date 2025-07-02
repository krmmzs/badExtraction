# badExtraction

A Go tool for cleaning up accidentally extracted archive files.

## Why does this project exist?

We often encounter this situation: accidentally extracting an archive in a directory that already contains many files, causing the archive's contents to scatter and mix with existing files. It becomes difficult to distinguish between original files and newly extracted ones, making manual cleanup tedious and error-prone.

The `cleanup` tool analyzes the first-level contents of an archive, automatically identifies files and folders in the current directory that may have come from that archive, and provides safe deletion functionality.

## Features

- **Multi-format support**: Supports ZIP, TAR, TAR.GZ, TAR.BZ2, TAR.XZ and other common archive formats
- **Smart detection**: Automatically analyzes archive's first-level contents and matches corresponding files/folders in current directory
- **Safe confirmation**: Lists all items to be deleted before execution, requires user confirmation
- **Error handling**: Comprehensive error messages and exception handling

## Installation

### Using Makefile

```bash
# Clone the project
git clone <repository-url>
cd badExtraction

# Build and install to $HOME/bin
make install

# Or just build
make build

# Clean build artifacts
make clean

# Uninstall the program
make uninstall
```

### Manual Build and Install

```bash
# Build
go build -o cleanup cleanup.go

# Install to system path
sudo cp cleanup /usr/local/bin/
# Or install to user path
cp cleanup $HOME/bin/
```

## Usage

```bash
cleanup <archive-filename>
```

### Examples

```bash
# Handle ZIP files
cleanup archive.zip

# Handle TAR.XZ files
cleanup software-1.0.tar.xz

# Handle TAR.GZ files
cleanup project.tar.gz
```

### Workflow

1. Run the `cleanup` command in the directory containing accidentally extracted files
2. The program analyzes the first-level contents of the specified archive
3. Displays a list of matching files and folders in the current directory
4. Confirm deletion (enter `y` or `yes` to confirm, any other input cancels)
5. The program executes deletion and shows results

### Example Output

```
$ cleanup example.tar.xz
Found the following files/folders that may have come from archive extraction:
  - src/
  - README.md
  - Makefile
Confirm deletion of these files and folders? (y/N): y
Deleted: src/
Deleted: README.md
Deleted: Makefile
Cleanup complete
```

## Supported Archive Formats

- `.zip` - ZIP archive files
- `.tar` - TAR archive files
- `.tar.gz` / `.tgz` - TAR + GZIP compressed files
- `.tar.bz2` - TAR + BZIP2 compressed files
- `.tar.xz` - TAR + XZ compressed files

## Safety Notes

- The program only deletes files and folders that match the archive's first-level contents
- Shows complete list and requires user confirmation before deletion
- Provides clear error messages if archive doesn't exist or format is unsupported
- Deletion operations are irreversible, please confirm carefully

## System Requirements

- Go 1.16 or higher (for compilation)
- Linux/macOS/Windows systems
- For `.tar.xz` format, requires `xz` tool to be installed on the system

## Development

```bash
# Get help
make help

# Build
make build

# Install
make install

# Clean
make clean
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.