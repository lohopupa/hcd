# hcd

**hcd** (History Change Directory) is a simple Golang application designed to help you find directories from your shell history using Levenshtein distance for filtering. It enhances your command-line experience by allowing you to quickly access previously visited directories.

## Features

- Finds directories based on a substring match using Levenshtein distance.
- Stores your directory history in `~/.cd_history` for persistent access.
- Easy integration with your shell.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/lohopupa/hcd.git
   cd hcd
   ```

2. Build the application:

   ```bash
   go build
   ```

3. Move the binary to a directory in your PATH (e.g., `/usr/local/bin`):

   ```bash
   mv hcd /usr/local/bin/
   ```

## Usage

### Setting Up Directory History

To enable directory tracking, add the following function to your `~/.bashrc` file:

```bash
cd() {
    builtin cd "$@" && pwd >> ~/.cd_history
}
```

After adding this, restart your terminal or run `source ~/.bashrc` to apply the changes.

### Finding Directories

Run the `hcd` command without any arguments to open the TUI, where you can browse through directories in your history:

```bash
hcd
```

The TUI will display directories that match your input based on the Levenshtein distance, allowing you to select and get path to the desired directory.


## Contributing

Feel free to submit issues or pull requests for any features or improvements you'd like to see.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.