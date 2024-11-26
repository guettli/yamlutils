package yamlutils

import (
	"fmt"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var nodeKindToString = map[yaml.Kind]string{
	yaml.DocumentNode: "DocumentNode",
	yaml.SequenceNode: "SequenceNode",
	yaml.MappingNode:  "MappingNode",
	yaml.ScalarNode:   "ScalarNode",
	yaml.AliasNode:    "AliasNode",
}

// NestedString traverses the given Node to find a nested string slice based on the provided fields.
func NestedString(node *yaml.Node, fields ...string) (s string, found bool, err error) {
	n, found, err := NestedNode(node, fields...)
	if err != nil {
		return "", false, err
	}
	if !found {
		return "", false, nil
	}

	if n.Kind != yaml.ScalarNode {
		return "", false, fmt.Errorf("NestedString: expected a scalar node, but got kind %s: %q", nodeKindToString[n.Kind], n.Value)
	}

	return n.Value, true, nil
}

// NestedStringSlice traverses the given Node to find a nested string slice based on the provided fields.
func NestedStringSlice(node *yaml.Node, fields ...string) ([]string, bool, error) {
	n, found, err := NestedNode(node, fields...)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	if n.Kind != yaml.SequenceNode {
		return nil, false, fmt.Errorf("NestedStringSlice: expected a sequence node, but got kind %s: %q", nodeKindToString[n.Kind], n.Value)
	}

	var result []string
	for _, item := range n.Content {
		if item.Kind != yaml.ScalarNode {
			return nil, false, fmt.Errorf("NestedStringSlice: expected a scalar node in sequence, but got kind %s", nodeKindToString[item.Kind])
		}
		result = append(result, item.Value)
	}

	return result, true, nil
}

// NestedStringMap traverses the given Node to find a nested map[string]string based on the provided fields.
func NestedStringMap(node *yaml.Node, fields ...string) (m map[string]string, found bool, err error) {
	n, found, err := NestedNode(node, fields...)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	if n.Kind != yaml.MappingNode {
		return nil, false, fmt.Errorf("NestedStringMap: expected a mapping node, but got kind %s", nodeKindToString[n.Kind])
	}
	if len(n.Content) == 0 {
		return nil, true, nil
	}
	result := make(map[string]string)
	for i := 0; i < len(n.Content); i += 2 {
		keyNode := n.Content[i]
		valueNode := n.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode || valueNode.Kind != yaml.ScalarNode {
			return nil, false, fmt.Errorf("expected scalar nodes for key-value pair, but got key kind %s and value kind %s",
				nodeKindToString[keyNode.Kind], nodeKindToString[valueNode.Kind])
		}

		result[keyNode.Value] = valueNode.Value
	}

	return result, true, nil
}

// NestedNode traverses the given Node to find the yaml.Node based on the provided fields.
func NestedNode(n *yaml.Node, fields ...string) (node *yaml.Node, found bool, err error) {
	current := n

	if n.Kind == yaml.DocumentNode {
		if len(n.Content) != 1 {
			return nil, false, fmt.Errorf("NestedNode: expected a single document node, but got %d", len(n.Content))
		}
		current = n.Content[0]
	}

	for _, field := range fields {
		if current.Kind != yaml.MappingNode {
			return nil, false, fmt.Errorf("NestedNode: expected a mapping node at field '%s', but got kind %s", field, nodeKindToString[current.Kind])
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

	return current, true, nil
}
