# direnv-to-dotenv

This tool translates .envrc formats (direnv) to .env formats. The motivation is to mostly retain a single source of truth for env vars with Jetbrain IDEs - the run configuration does env formats.

To install the tool:

```
go get github.com/veedubyou/direnv-to-dotenv
```

To use the tool, run:

```
direnv allow
go run github.com/veedubyou/direnv-to-dotenv
```

And the environment variables should be copied to the clipboard.
