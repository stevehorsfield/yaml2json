package main

import "os"
import "flag"
import "fmt"
import "bytes"

import "github.com/stevehorsfield/yaml2json/internal/yaml2json"

func customUsage(fs *flag.FlagSet) {
	fmt.Fprintf(fs.Output(), "yaml2json - Convert YAML to equivalent JSON\n")
	fmt.Fprintf(fs.Output(), "(c) 2021 Stephen Horsfield\n\n")

	fs.PrintDefaults()

	fmt.Fprintf(fs.Output(), `
Supports basic YAML primitives with equivalent content in JSON.
Does not support Alias nodes or custom types.

This project is available at https://github.com/stevehorsfield/yaml2json
YAML processing provided by https://gopkg.in/yaml.v3
`)
}

var flagCompact bool
var flagMultiYAML bool

func initf() {
	fs := flag.NewFlagSet("yamltojson", flag.ExitOnError)

	fs.BoolVar(&flagCompact, "c", false, "compact JSON output")
	fs.BoolVar(&flagMultiYAML, "s", false, "process multiple YAML inputs as a JSON array")

	fs.Usage = func() { customUsage(fs) }
	_ = fs.Parse(os.Args[1:])
}

func main() {

	initf()

	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %s", err)
		os.Exit(1)
	}
	content := buf.Bytes()

	s, err := yaml2json.Process(content, yaml2json.Options{CompactJSON: flagCompact, MultipleYAML: flagMultiYAML})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(s)
}
