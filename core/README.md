# Core

The core of the system is a server that listens for connections from the libraries and tracks the requests.

## Configuration file

You configure the core using a JSON configuration file.  By default, core looks for a file in the working directory named "config.json", but you can specify the path to the configuration file with the `-c` command line option.  There is an example configuration file included in this repository.  For a description of all of the configuration options, see the package comment in `config/config.go`.
