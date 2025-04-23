package configfile_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/michaellenaghan/go-configfile"
)

func Example() {
	// Create a temporary config file for the example
	tmpfile, err := os.CreateTemp("", "config_example")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte(`
		# Database configuration
		db-url = example.com:5432

		# Server configuration
		server-port = 9090

		# Enable debugging
		debug = true
	`)); err != nil {
		fmt.Printf("Failed to write to temp file: %v\n", err)
		return
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Printf("Failed to close temp file: %v\n", err)
		return
	}

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
	os.Args = []string{"program", "--config-file=" + tmpfile.Name()}
	flag.Parse()

	// Print the flags
	fmt.Printf("Server running on host: %s\n", *serverHost)
	fmt.Printf("Server running on port: %d\n", *serverPort)
	fmt.Printf("Connected to database: %s\n", *dbURL)
	fmt.Printf("Debug mode: %t\n", *debug)

	// Output:
	// Server running on host: localhost
	// Server running on port: 9090
	// Connected to database: example.com:5432
	// Debug mode: true
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
		wantErr  bool
	}{
		{
			name: "Empty file",
			content: `
				# Just a comment
			`,
			expected: map[string]string{},
			wantErr:  false,
		},
		{
			name: "Empty name",
			content: `
				= value1
				key2 = value2
			`,
			wantErr: true,
		},
		{
			name: "Empty value",
			content: `
				key1 = 
				key2 = value2
			`,
			expected: map[string]string{
				"key1": "",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "Flag overriding",
			content: `
				key1 = value1
				key1 = value2
				key1 = value3
			`,
			expected: map[string]string{
				"key1": "value3",
			},
			wantErr: false,
		},
		{
			name: "Invalid line",
			content: `
				key1 = value1
				invalid_line
				key2 = value2
			`,
			wantErr: true,
		},
		{
			name: "Unknown flag",
			content: `
				key1 = value1
				unknown_flag = value
			`,
			wantErr: true,
		},
		{
			name: "Valid config",
			content: `
				key1 = value1
				key2 = value2
				# Comment
				key3 = value3
			`,
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			wantErr: false,
		},
		{
			name: "Whitespace handling",
			content: `
				  key1    =    value1    
				key2=value2
				   key3 = value3   
			`,
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			wantErr: false,
		},
	}

	// Define normal string flags
	_ = flag.String("key1", "", "key1 flag")
	_ = flag.String("key2", "", "key2 flag")
	_ = flag.String("key3", "", "key3 flag")

	// Define an int flag to test invalid value parsing
	_ = flag.Int("int_flag", 0, "int flag")

	t.Run("Invalid flag value", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "config_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		content := "int_flag = not_an_int"
		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		if err := configfile.Load(tmpfile.Name()); err == nil {
			t.Error("Load() expected error for invalid flag value, got nil")
		}
	})

	t.Run("Non-existent file", func(t *testing.T) {
		if err := configfile.Load("Non-existent file"); err == nil {
			t.Errorf("Load() error = %v, wantErr <not-nil>", err)
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "config_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if tt.content != "" {
				if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
					t.Fatal(err)
				}
				if err := tmpfile.Close(); err != nil {
					t.Fatal(err)
				}
			}

			err = configfile.Load(tmpfile.Name())

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for k, v := range tt.expected {
					if flagValue := flag.Lookup(k).Value.String(); flagValue != v {
						t.Errorf("Expected flag %s to be set to %s, but got %s", k, v, flagValue)
					}
				}
			}
		})
	}
}
