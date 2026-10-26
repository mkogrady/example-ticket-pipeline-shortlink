Build: `go build ./...`. Test: `go test ./...`. Format: `gofmt -w .`.
Standard library only. Do not add dependencies.
Each package has its own tests next to it; add a test for every change in behavior.
Keep handlers in `server`, storage in `store`, input rules in `validate`.
Commit messages are one short imperative line.
