# jtl

A `cli` that can perform arbritrary manipulation of `json` payloads.

## Example Usage

`jtl -spec=example/spec.yaml -source=example/input.json`


## Spec

`jtl` uses specs defined in `yaml` to understand the transformations that need
to be performed. Keys in the source `json` payload are omitted from the output
if they're not defined in the spec.


A basic `Spec` has the following attributes:

- `sourcePath` - A `json-esque` path to the key in the source document. Keys
  within nested objects can be specified using dot notation. e.g `nested.key`.
- `targetPath` - A `json-esque` path to the key in the transformed document. To
  put the key within a nested object dot notation can be used.

A `Spec` can also have multiple other `specs` within it:

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

When a `Spec` has child `Specs` the functionality of the top-level `sourcePath`
and `targetPath` change slightly. `sourcePath` becomes the start point in the
source object and `targetPath` becomes the "start" point for transformations
in the target object.

### Conditionals

TODO

 - `sourcePath`
 - `operator`
 - `value`
