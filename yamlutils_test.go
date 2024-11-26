package yamlutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var fooString = `apiVersion: v1
kind: Foo
metadata:
  annotations: {}
  labels:
    foo: bar
  name: test
  strings:
  - one
  - two
  stringsEmpty: []
`

func TestNestedStringMap(t *testing.T) {
	d := yaml.NewDecoder(strings.NewReader(fooString))
	node := &yaml.Node{}
	err := d.Decode(node)
	require.NoError(t, err)

	l, found, err := NestedStringMap(node, "metadata", "does-not-exist")
	require.NoError(t, err)
	require.Equal(t, map[string]string(nil), l)
	require.True(t, found)

	l, found, err = NestedStringMap(node, "metadata", "labels")
	require.NoError(t, err)
	require.Equal(t, map[string]string{"foo": "bar"}, l)
	require.False(t, found)

	l, found, err = NestedStringMap(node, "metadata", "annotations")
	require.NoError(t, err)
	require.Equal(t, map[string]string(nil), l)
	require.True(t, found)

	l, found, err = NestedStringMap(node, "metadata", "name")
	require.NotNil(t, err)
	require.Equal(t, map[string]string(nil), l)
	require.False(t, found)
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

	l, found, err = NestedStringSlice(node, "metadata", "stringsEmpty")
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

func TestNestedString(t *testing.T) {
	d := yaml.NewDecoder(strings.NewReader(fooString))
	node := &yaml.Node{}
	err := d.Decode(node)
	require.NoError(t, err)

	l, found, err := NestedString(node, "metadata", "name")
	require.NoError(t, err)
	require.Equal(t, "test", l)
	require.True(t, found)

	l, found, err = NestedString(node, "metadata", "does-not-exist")
	require.NoError(t, err)
	require.Equal(t, "", l)
	require.False(t, found)

	l, found, err = NestedString(node, "metadata", "labels")
	require.NotNil(t, err)
	require.Equal(t, "", l)
	require.False(t, found)
}
