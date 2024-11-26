package yamlutils

import (
	"fmt"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

func NestedString(n *yaml.Node, fields ...string) (s string, found bool, err error) {
	current := n

	if n.Kind == yaml.DocumentNode {
		if len(n.Content) != 1 {
			return "", false, fmt.Errorf("expected a single document node, but got %d", len(n.Content))
		}
		current = n.Content[0]
	}

	for _, field := range fields {
		if current.Kind != yaml.MappingNode {
			return "", false, fmt.Errorf("expected a mapping node at field '%s', but got kind %d", field, current.Kind)
		}

		found := false
		for i := 0; i < len(current.Content); i += 2 {
			keyNode := current.Content[i]
			valueNode := current.Content[i+1]
			if keyNode.Value == field {
				current = valueNode
				found = true
				break
			}
		}

		if !found {
			return "", false, nil
		}
	}

	if current.Kind != yaml.ScalarNode {
		return "", false, fmt.Errorf("expected a scalar node, but got kind %d", current.Kind)
	}

	return current.Value, true, nil
}

// NestedStringSlice traverses the given Node to find a nested string slice based on the provided fields.
func NestedStringSlice(n *yaml.Node, fields ...string) ([]string, bool, error) {
	current := n
	if n.Kind == yaml.DocumentNode {
		if len(n.Content) != 1 {
			return nil, false, fmt.Errorf("expected a single document node, but got %d", len(n.Content))
		}
		current = n.Content[0]
	}
	for _, field := range fields {
		if current.Kind != yaml.MappingNode {
			return nil, false, fmt.Errorf("expected a mapping node at field '%s', but got kind %d", field, current.Kind)
		}

		found := false
		for i := 0; i < len(current.Content); i += 2 {
			keyNode := current.Content[i]
			valueNode := current.Content[i+1]
			if keyNode.Value == field {
				current = valueNode
				found = true
				break
			}
		}

		if !found {
			return nil, false, nil
		}
	}

	if current.Kind != yaml.SequenceNode {
		return nil, false, fmt.Errorf("expected a sequence node, but got kind %d: %q", current.Kind, current.Value)
	}

	var result []string
	for _, item := range current.Content {
		if item.Kind != yaml.ScalarNode {
			return nil, false, fmt.Errorf("expected a scalar node in sequence, but got kind %d", item.Kind)
		}
		result = append(result, item.Value)
	}

	return result, true, nil
}

// NestedStringMap traverses the given Node to find a nested map[string]string based on the provided fields.
func NestedStringMap(n *yaml.Node, fields ...string) (m map[string]string, found bool, err error) {
	current := n
	if n.Kind == yaml.DocumentNode {
		if len(n.Content) != 1 {
			return nil, false, fmt.Errorf("expected a single document node, but got %d", len(n.Content))
		}
		current = n.Content[0]
	}

	for _, field := range fields {
		if current.Kind != yaml.MappingNode {
			return nil, false, fmt.Errorf("expected a mapping node at field '%s', but got kind %d", field, current.Kind)
		}

		found := false
		for i := 0; i < len(current.Content); i += 2 {
			keyNode := current.Content[i]
			valueNode := current.Content[i+1]
			// fmt.Printf("keyNode: %+v valueNode %+v\n", keyNode, valueNode)
			if keyNode.Value == field {
				current = valueNode
				found = true
				break
			}
		}

		if !found {
			return nil, true, nil
		}
	}

	if current.Kind != yaml.MappingNode {
		return nil, false, fmt.Errorf("expected a mapping node, but got kind %d", current.Kind)
	}
	if len(current.Content) == 0 {
		return nil, true, nil
	}
	result := make(map[string]string)
	for i := 0; i < len(current.Content); i += 2 {
		keyNode := current.Content[i]
		valueNode := current.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode || valueNode.Kind != yaml.ScalarNode {
			return nil, false, fmt.Errorf("expected scalar nodes for key-value pair, but got key kind %d and value kind %d", keyNode.Kind, valueNode.Kind)
		}

		result[keyNode.Value] = valueNode.Value
	}

	return result, false, nil
}
