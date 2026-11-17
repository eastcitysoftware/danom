# Changelog

## v0.1.2 - 2026-11-17

### Breaking changes

- The package-level `option.Do` and `option.Iter` functions are now methods on `Option`; replace `option.Do(o, some, none)` and `option.Iter(o, fn)` with `o.Do(some, none)` and `o.Iter(fn)`. Likewise, replace `result.Do(r, ok, error)` and `result.Iter(r, fn)` with `r.Do(ok, error)` and `r.Iter(fn)`.
- The type parameters for `result.Ok` are now ordered as `[E, T]`, allowing Go to infer `T` from the value argument. Replace calls such as `result.Ok[int, error](value)` with `result.Ok[error](value)`.

### Added

- Added `Do` methods to handle both states of an `Option` or `Result`.
- Added `Iter` methods to run a callback only for a present `Option` value or successful `Result`.

## v0.1.1 - 2026-11-16

### Added

- Add JSON marshaling and unmarshaling for `option.Option[T]`. `None` encodes as
  `null`, `Some` encodes as its value, and JSON `null` decodes to `None`.
- Add `Option.IsZero` for use with `encoding/json`'s `omitzero` tag, and add
  `Option.String` for readable formatting.
- Add native-width `parse.Int` and `parse.Uint` functions, with matching
  `convert.Int` and `convert.Uint` functions.

### Changed

- Lower the minimum Go version to 1.21.
- Change `check.MapAll` and `check.MapAny` to accept predicates on map values
  instead of lists of map facts.
- Organize validation facts into type-specific source files.

### Removed

- Remove the `check.StrIsEmail` and `check.StrIsE164` validation facts.

## v0.1.0 - 2026-11-11

Hello world