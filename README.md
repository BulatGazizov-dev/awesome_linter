## Build

```bash
make build 
```

Note: If building manually via golangci-lint custom -v, ensure your .golangci.yml points to the correct local module path.

## Usage
- Rename .golangci.example.yml to .golangci.yml.
- Configure your rules under linters-settings.custom.awesome_linter.`

```bash
custom-gcl run (path)
```

_If you built linter manually, you may have different name of executable file_

## Testing

```bash
make test
```