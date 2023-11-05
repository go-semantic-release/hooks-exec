# :mechanical_arm: hooks-exec
[![CI](https://github.com/go-semantic-release/hooks-exec/workflows/CI/badge.svg?branch=main)](https://github.com/go-semantic-release/hooks-exec/actions?query=workflow%3ACI+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-semantic-release/hooks-exec)](https://goreportcard.com/report/github.com/go-semantic-release/hooks-exec)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-semantic-release/hooks-exec)](https://pkg.go.dev/github.com/go-semantic-release/hooks-exec)

This plugin allows to execute custom scripts during the semantic-release process.

## Usage
Use this plugin by enabling it via `--hooks exec`.

## Configuration
There are two possible scripts that can be executed during the release process: (i) the success script, if the release was successful, and (ii) the no release script, if the was not created for a specific reason.

The scripts are configured via the `--hooks-opt exec_on_success="echo v{{.NewRelease.Version}} "` and `--hooks-opt exec_on_no_release="echo {{.Reason}}"` flags. Additionally, they `.semrelrc` file can be used to configure the scripts:
```json
{
  "plugins": {
    "hooks": {
      "names": [
        "exec"
      ],
      "options": {
        "exec_on_success": "echo v{{.PrevRelease.Version}} '->' v{{.NewRelease.Version}}",
        "exec_on_no_release": "echo {{.Reason}}: {{.Message}}"
      }
    }
  }
}

```

### Success script configuration
The success script configuration accepts Go templates with the following variables:
- `{{.PrevRelease.Version}}`: The previous release version.
- `{{.NewRelease.Version}}`: The new release version.
- `{{.Changelog}}`: The release changelog.
- `{{.Commits}}`: The commits between the previous and the new release.
- `{{.RepoInfo}}`: The repository information.

A more detailed documentation about the available variables can be found [here](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/hooks#SuccessHookConfig).

### No release script configuration
The no release script configuration also accepts Go templates with the following variables:
- `{{.Reason}}`: The reason why no release was created. One of 
- `{{.Message}}`: The message why no release was created

A more detailed documentation about the available variables can be found [here](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/hooks#NoReleaseConfig).

## Licence

The [MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright © 2023 [Christoph Witzko](https://twitter.com/christophwitzko)
