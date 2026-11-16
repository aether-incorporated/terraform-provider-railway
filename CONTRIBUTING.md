## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Generating

To generate or update documentation and GraphQL client, run `go generate`.

## Testing

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

Run them against a dedicated Railway workspace, never one holding real projects. Create one fixture project there with an environment named `staging` and one service in it, then set:

| Variable                      | Value                                             |
| ----------------------------- | ------------------------------------------------- |
| `RAILWAY_TOKEN`               | A token for the test workspace                    |
| `RAILWAY_TEST_WORKSPACE_ID`   | The test workspace's id                           |
| `RAILWAY_TEST_PROJECT_ID`     | The fixture project's id                          |
| `RAILWAY_TEST_ENVIRONMENT_ID` | The id of the fixture project's `staging` environment |
| `RAILWAY_TEST_SERVICE_ID`     | The fixture service's id                          |

In CI, `RAILWAY_TOKEN` is a repository secret and the four ids are repository variables. Without the secret, the acceptance tests are skipped and only the unit tests run.

```shell
make testacc
```
