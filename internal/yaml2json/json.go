package yaml2json

import "bytes"
import "encoding/json"

func toJSON(data *interface{}, compact bool) (string, error) {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	if !compact {
		enc.SetIndent("", "  ")
	}

	err := enc.Encode(data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
