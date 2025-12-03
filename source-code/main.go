package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage:")
		fmt.Println(" toml pretty <file or ->")
		fmt.Println(" toml get <path> <file or ->")
		fmt.Println("Path example: .key.subkey[0]")
		os.Exit(0)
	}
	command := args[0]
	var input string
	if command == "get" {
		if len(args) < 3 {
			fmt.Println("Usage: toml get <path> <file or ->")
			os.Exit(0)
		}
		path := args[1]
		input = args[2]
		data := readInput(input)
		var value map[string]interface{}
		if _, err := toml.Decode(data, &value); err != nil {
			fmt.Printf("Error parsing TOML: %v\n", err)
			os.Exit(1)
		}
		result := getValue(value, path)
		printJSON(result)
	} else if command == "pretty" {
		if len(args) < 2 {
			fmt.Println("Usage: toml pretty <file or ->")
			os.Exit(0)
		}
		input = args[1]
		data := readInput(input)
		var value map[string]interface{}
		if _, err := toml.Decode(data, &value); err != nil {
			fmt.Printf("Error parsing TOML: %v\n", err)
			os.Exit(1)
		}
		pretty, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			fmt.Printf("Error pretty-printing: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(pretty))
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func readInput(input string) string {
	if input == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		return string(data)
	} else {
		data, err := os.ReadFile(input)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", input, err)
			os.Exit(1)
		}
		return string(data)
	}
}

func getValue(data interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := data
	for _, part := range parts {
		if part == "" {
			continue
		}
		if strings.HasSuffix(part, "]") {
			openBracket := strings.Index(part, "[")
			if openBracket == -1 {
				return nil
			}
			key := part[:openBracket]
			indexStr := part[openBracket+1 : len(part)-1]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil
			}
			if key != "" {
				if m, ok := current.(map[string]interface{}); ok {
					current = m[key]
				} else {
					return nil
				}
			}
			if arr, ok := current.([]interface{}); ok {
				if index < len(arr) {
					current = arr[index]
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else {
			if m, ok := current.(map[string]interface{}); ok {
				current = m[part]
			} else {
				return nil
			}
		}
	}
	return current
}

func printJSON(value interface{}) {
	if value == nil {
		fmt.Println("null")
		return
	}
	pretty, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("Error printing: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(pretty))
}
