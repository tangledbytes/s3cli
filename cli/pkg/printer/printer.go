package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/tidwall/pretty"
)

// Print takes any data and performs a pretty print
func Print(data any, color bool, raw string) error {
	if raw != "" {
		printRaw(data, raw)
		return nil
	}

	byt, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if color {
		byt = pretty.Color(byt, nil)
	}

	fmt.Println(string(byt))

	return nil
}

func printRaw(data any, raw string) error {
	temp, err := template.New("raw").Parse(raw)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return temp.Execute(os.Stdout, data)
}
