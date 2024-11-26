# yamlutils: Go package for working with yaml

## Introduction

Imagine you have parsed a yaml like this with [sigs.k8s.io/yaml/goyaml.v3](https://pkg.go.dev/sigs.k8s.io/yaml/goyaml.v3):

```yaml
apiVersion: example.com/v1
kind: Foo
metadata:
  labels:
    foo: bar
spec:
  level1:
    level2:
      level3:
        field: myValue

```

You can easily get the value of the deeply nested field like this:

```go
s, found, err := NestedString(node, "spec", "level1", "level2", "level3", "field")
```

In above example `s` contains "myValue", found is true, and err is nil.

If you try to access a field which does not exit, `found` will be false.

If the value you try to get exists, but the data types does not match, `err` will be non-nil.

## API Docs

<https://pkg.go.dev/github.com/guettli/yamlutils>

## Feedback is welcome

Please create an issue if you find a typo or if you have other ideas how to improve this package.

## Releasing

```terminal
go test ./...

RELEASE_TAG=v0.0.X

git tag $RELEASE_TAG

git push origin $RELEASE_TAG
```
