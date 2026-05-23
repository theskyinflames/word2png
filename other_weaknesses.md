# Other Weaknesses

- Non-canonical module path + self-referencing dependency in `go.mod`
- No end-to-end test exercising the *real* crypto (only mocked)
- `os.Exit()` skips deferred cleanup in both CLIs
- Minor: exit code `-1` (non-standard), inconsistent error naming
