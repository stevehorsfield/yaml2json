package yaml2json

import "fmt"

type Options struct {
	CompactJSON  bool
	MultipleYAML bool
}

func Process(content []byte, opts Options) (string, error) {
	var jsonData interface{}
	var err error
	if opts.MultipleYAML {
		var data [][]byte
		data, err := splitYAML(content)

		if err != nil {
			return "", fmt.Errorf("Error splitting YAML content: %s", err)
		}

		result := make([]interface{}, 0)
		for i, v := range data {

			o, err := yamlToObject(v)
			if err != nil {
				return "", fmt.Errorf("Unable to process multi-document YAML at document index %d: %s", i, err)
			}

			result = append(result, o)
		}

		jsonData = result

	} else {
		jsonData, err = yamlToObject(content)
		if err != nil {
			return "", err
		}
	}

	return toJSON(&jsonData, opts.CompactJSON)
}
