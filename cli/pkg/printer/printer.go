package printer

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/pretty"
)

// Print takes any data and performs a pretty print
func Print(data any, color, raw bool) error {
	if raw {
		fmt.Printf("%#+v\n", data)
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
