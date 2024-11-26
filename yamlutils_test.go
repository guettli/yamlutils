package yamlutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var fooString = `apiVersion: example.com/v1
kind: Foo
metadata:
  labels:
    foo: bar
  name: test
  strings:
  - one
  - two
  emptyString: ""
  emptySequence: []
  emptyMap: {}
spec:
  level1:
    level2:
      level3:
        field: myValue
foo:
  foo:
    foo: bar
`

func TestNestedString(t *testing.T) {
	d := yaml.NewDecoder(strings.NewReader(fooString))
	node := &yaml.Node{}
	err := d.Decode(node)
	require.NoError(t, err)

	s, found, err := NestedString(node, "metadata", "name")
	require.NoError(t, err)
	require.Equal(t, "test", s)
	require.True(t, found)

	s, found, err = NestedString(node, "metadata", "emptyString")
	require.NoError(t, err)
	require.Equal(t, "", s)
	require.True(t, found)

	s, found, err = NestedString(node, "metadata", "does-not-exist")
	require.NoError(t, err)
	require.Equal(t, "", s)
	require.False(t, found)

	s, found, err = NestedString(node, "metadata", "labels")
	require.NotNil(t, err)
	require.Equal(t, "", s)
	require.False(t, found)

	s, found, err = NestedString(node, "spec", "level1", "level2", "level3", "field")
	require.NoError(t, err)
	require.Equal(t, "myValue", s)
	require.True(t, found)

	s, found, err = NestedString(node, "foo", "foo", "foo")
	require.NoError(t, err)
	require.Equal(t, "bar", s)
	require.True(t, found)
}

func TestNestedStringSlice(t *testing.T) {
	d := yaml.NewDecoder(strings.NewReader(fooString))
	node := &yaml.Node{}
	err := d.Decode(node)
	require.NoError(t, err)

	l, found, err := NestedStringSlice(node, "metadata", "strings")
	require.NoError(t, err)
	require.Equal(t, []string{"one", "two"}, l)
	require.True(t, found)

	l, found, err = NestedStringSlice(node, "metadata", "emptySequence")
	require.NoError(t, err)
	require.Equal(t, []string(nil), l)
	require.True(t, found)

	l, found, err = NestedStringSlice(node, "metadata", "does-not-exist")
	require.NoError(t, err)
	require.Equal(t, []string(nil), l)
	require.False(t, found)

	l, found, err = NestedStringSlice(node, "metadata", "labels")
	require.NotNil(t, err)
	require.Equal(t, []string(nil), l)
	require.False(t, found)
}

func TestNestedNode(t *testing.T) {
	d := yaml.NewDecoder(strings.NewReader(fooString))
	node := &yaml.Node{}
	err := d.Decode(node)
	require.NoError(t, err)

	n, found, err := NestedNode(node, "spec", "level1", "level2")
	require.NoError(t, err)
	b := &bytes.Buffer{}
	e := yaml.NewEncoder(b)
	err = e.Encode(n)
	require.NoError(t, err)
	require.Equal(t, "level3:\n    field: myValue\n", b.String())
	require.True(t, found)
}

func TestNestedStringMap(t *testing.T) {
	d := yaml.NewDecoder(strings.NewReader(fooString))
	node := &yaml.Node{}
	err := d.Decode(node)
	require.NoError(t, err)

	m, found, err := NestedStringMap(node, "metadata", "does-not-exist")
	require.NoError(t, err)
	require.Equal(t, map[string]string(nil), m)
	require.False(t, found)

	m, found, err = NestedStringMap(node, "metadata", "labels")
	require.NoError(t, err)
	require.Equal(t, map[string]string{"foo": "bar"}, m)
	require.True(t, found)

	m, found, err = NestedStringMap(node, "metadata", "emptyMap")
	require.NoError(t, err)
	require.Equal(t, map[string]string(nil), m)
	require.True(t, found)

	m, found, err = NestedStringMap(node, "metadata", "name")
	require.NotNil(t, err)
	require.Equal(t, map[string]string(nil), m)
	require.False(t, found)
}
