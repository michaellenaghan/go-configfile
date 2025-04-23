// Package configfile provides functionality for loading flag-based
// configuration values from a file. It integrates directly with the standard
// flag package, allowing configuration values to be defined in files with
// exactly the same semantics as command-line flags.
//
// You can load config files programmatically:
//
//  if _, err := os.Stat("./defaults.conf"); err == nil {
//    err := configfile.Load("./defaults.conf")
//    if err != nil {
//      fmt.Printf("Error loading default config: %v\n", err)
//    }
//  }

// You can load config files on the command-line:
//
//	flag.Func("config-file", "Configuration file", configfile.Load)
//
// Multiple config files can be loaded by passing the same flag multiple
// times:
//
//	myapp --config-file=/etc/myapp/global.conf --config-file=./local.conf
//
// It's possible to combine the two approaches, loading default config
// files first and then allowing users to load additional config files on
// the command-line.
//
// The order of processing follows standard flag package behavior, with
// later files (and values) overriding earlier ones.
package configfile

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Load reads configuration values from the specified configuration file and
// sets them using the flag package.
//
// The configuration file should be in the format of "key=value" pairs, one per
// line. The keys correspond directly to flag names defined in the program.
//
// For example, if the program defines flags "--database-url", "--server-host"
// and "--server-port", the configuration file can contain:
//
//	database-url=localhost:5432
//	server-host=localhost
//	server-port=9090
//
// Empty lines are ignored, and lines starting with '#' are treated as
// comments.
//
// Whitespace around keys and values is trimmed.
//
// If any invalid lines are encountered (e.g., missing '=' separator), Load
// returns an error.
//
// If a key corresponds to a flag that hasn't been defined, or if the value
// isn't valid for the flag's type, Load returns an error, just as would happen
// with invalid command-line flags.
//
// If the file cannot be opened or read, Load returns an error with details.
func Load(configfile string) error {
	file, err := os.Open(configfile)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", configfile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("failed to split line (expected to find an '='): %s", line)
		}

		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)

		if err := flag.Set(name, value); err != nil {
			return fmt.Errorf("failed to set flag '%s' to value '%s': %w", name, value, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}
