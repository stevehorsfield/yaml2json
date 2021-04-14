package yaml2json

import "fmt"
import "strconv"

import "gopkg.in/yaml.v3"

func yamlToObject(content []byte) (interface{}, error) {
	var n yaml.Node
	err := yaml.Unmarshal(content, &n)
	if err != nil {
		return nil, fmt.Errorf("Error translating YAML to JSON: %s", err)
	}

	return yamlNode(&n)
}

func yamlNode(n0 *yaml.Node) (interface{}, error) {
	switch n0.Kind {
	case yaml.DocumentNode:
		return yamlDocument(n0)
	case yaml.SequenceNode:
		return yamlSequence(n0)
	case yaml.MappingNode:
		return yamlMapping(n0)
	case yaml.ScalarNode:
		return yamlScalar(n0)
	case yaml.AliasNode:
		return nil, fmt.Errorf("No translation to JSON for AliasNode")
	default:
		return nil, fmt.Errorf("Unsupported YAML node type: %v", n0.Kind)
	}
}

func yamlDocument(n0 *yaml.Node) (interface{}, error) {
	if len(n0.Content) != 1 {
		return nil, fmt.Errorf("Unexpected YAML Document node content length: %d", len(n0.Content))
	}
	return yamlNode(n0.Content[0])
}

func yamlMapping(n0 *yaml.Node) (interface{}, error) {
	m := make(map[string]interface{})

	for i := 0; i < len(n0.Content); i += 2 {

		k, err := yamlStringScalar(n0.Content[i])
		if err != nil {
			return nil, fmt.Errorf("Unable to decode YAML map key: %s", err)
		}
		v, err := yamlNode(n0.Content[i+1])
		if err != nil {
			return nil, fmt.Errorf("Unable to process YAML map value for key '%s': %s", k, err)
		}
		m[k] = v
	}
	return m, nil
}

func yamlSequence(n0 *yaml.Node) (interface{}, error) {
	s := make([]interface{}, 0)

	for i := 0; i < len(n0.Content); i++ {

		v, err := yamlNode(n0.Content[i])
		if err != nil {
			return nil, fmt.Errorf("Unable to decode YAML sequence value: %s", err)
		}
		s = append(s, v)
	}
	return s, nil
}

func yamlScalar(n0 *yaml.Node) (interface{}, error) {
	switch n0.LongTag() { // See https://yaml.org/type/
	case "tag:yaml.org,2002:str":
		return n0.Value, nil
	case "tag:yaml.org,2002:bool":
		b, err := strconv.ParseBool(n0.Value)
		if err != nil {
			return nil, fmt.Errorf("Unable to process scalar node. Got '%s'. Expecting bool content: %s", n0.Value, err)
		}
		return b, nil
	case "tag:yaml.org,2002:int":
		i, err := strconv.ParseInt(n0.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to process scalar node. Got '%s'. Expecting integer content: %s", n0.Value, err)
		}
		return i, nil
	case "tag:yaml.org,2002:float":
		f, err := strconv.ParseFloat(n0.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to process scalar node. Got '%s'. Expecting float content: %s", n0.Value, err)
		}
		return f, nil
	case "tag:yaml.org,2002:null":
		return nil, nil
	default:
		return nil, fmt.Errorf("Error: YAML tag '%s' is not supported", n0.LongTag())
	}
}

func yamlStringScalar(n0 *yaml.Node) (string, error) {
	if n0.Kind != yaml.ScalarNode {
		return "", fmt.Errorf("Expecting a string scalar but got %q", n0.Kind)
	}
	if n0.LongTag() != "tag:yaml.org,2002:str" {
		return "", fmt.Errorf("Unable to process scalar node. Expecting string content")
	}
	return n0.Value, nil
}
