# jtl

A `cli` that can perform arbritrary manipulation of `json` payloads.

## Installation

`go get github.com/imjoshholloway/jtl/cmd/jtl`

## Example Usage

`jtl -spec=example/spec.yaml -source=example/input.json`


## Spec

`jtl` uses `specs` defined in `yaml` to decide what transformations need to be performed.

A basic `spec` has the following attributes:

- `sourcePath` - A `json-esque` path to the key in the source document. Keys within nested objects can be specified using dot notation. e.g `nested.key`.
- `targetPath` - A `json-esque` path to the key in the transformed document. To put the key within a nested object dot notation can be used.

**Note**: Keys in the source `json` payload are omitted from the output if they're not defined in the `spec`.

A `spec` can also have multiple other "child" `specs` within it:

```
sourcePath: contact_info
targetPath: contact
specs:
  - sourcePath: address.house_number
    targetPath: house_number

  - sourcePath: address.zip_code
    targetPath: zip_code

```

### Nested Specs

When a `spec` has child `specs` the functionality of the top-level `sourcePath` and `targetPath` change slightly. `sourcePath` becomes the start point in the source object and `targetPath` becomes the start point for transformations in the target object.

For example, given the following `spec`:
```
 - sourcePath: addresses
   targetPath: contact.addresses
   specs:
    - sourcePath: house_number
      targetPath: number

    - sourcePath: line1
      targetPath: line1
```

When applied to the following `json` payload:
```
{
  "addresses": [{
      "house_number": "10",
      "line1": "Main Street",
      "city": "London",
  }, {
      "house_number": "250",
      "line1": "Side street",
      "city": "Bristol"
  }],
}
```

The output would be:
```
{
  "contact": {
    "addresses": [{
        "number": "10",
        "line1": "Main Street"
    }, {
        "number": "250",
        "line1": "Side street"
    }]
  }
}

```

### Conditionals

TODO

 - `sourcePath`
 - `operator`
 - `value`