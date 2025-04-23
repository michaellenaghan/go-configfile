# go-configfile

Package `configfile` provides functionality for loading flag-based configuration values from a file. It integrates directly with the standard `flag` package, allowing configuration values to be defined in files with exactly the same semantics as command-line flags.

- Uses the standard `flag` package - no new semantics to learn
- Supports `#` comments
- Supports `key=value` lines
- Supports layering configurations from multiple files

The code is available at [github.com/michaellenaghan/go-configfile](https://github.com/michaellenaghan/go-configfile).

The documentation is available at [pkg.go.dev/github.com/michaellenaghan/go-configfile](https://pkg.go.dev/github.com/michaellenaghan/go-configfile).

## Installation

```bash
go get github.com/michaellenaghan/go-configfile
```

## Quick Start

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/michaellenaghan/go-configfile"
)

func main() {
	// Create the flags
	dbURL := flag.String("db-url", "localhost:5432", "Database URL")
	serverHost := flag.String("server-host", "localhost", "Server host")
	serverPort := flag.Int("server-port", 8080, "Server port")
	debug := flag.Bool("debug", false, "Enable debug mode")

	// You can load config files programmatically:
	if _, err := os.Stat("./defaults.conf"); err == nil {
		err := configfile.Load("./defaults.conf")
		if err != nil {
			fmt.Printf("Error loading default config: %v\n", err)
		}
	}

	// You can load config files on the command-line:
	flag.Func("config-file", "Configuration file", configfile.Load)

	// Multiple config files can be loaded by passing the same flag multiple
	// times:
	//
	//   myapp --config-file=/etc/myapp/global.conf --config-file=./local.conf
	//
	// It's possible to combine the two approaches, loading default config
	// files first and then allowing users to load additional config files on
	// the command-line.
	//
	// The order of processing follows standard flag package behavior, with
	// later files (and values) overriding earlier ones.

	// Parse the flags
	flag.Parse()

	// Print the flags
	fmt.Printf("Server host: %s\n", *serverHost)
	fmt.Printf("Server port: %d\n", *serverPort)
	fmt.Printf("Database URL: %s\n", *dbURL)
	fmt.Printf("Debug mode: %t\n", *debug)
}
```

## License

MIT License