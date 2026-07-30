# Terraform Railway Provider

A Terraform / OpenTofu provider for [Railway](https://railway.app), published as `aether-incorporated/railway`.

## Documentation

Documentation lives in this repository, not on a registry site:

- [Provider configuration and authentication](./docs/index.md)
- [Resources](./docs/resources)
- [Ephemeral resources](./docs/ephemeral-resources)

## Usage

```hcl
terraform {
  required_providers {
    railway = {
      source  = "aether-incorporated/railway"
      version = "~> 0.1.0"
    }
  }
}

provider "railway" {
  token = var.railway_token
}
```

The token can also come from the `RAILWAY_TOKEN` environment variable. Use an **account/workspace token** — a project-scoped token can read its own project but is rejected by most create and update operations.

OpenTofu resolves an unqualified source against `registry.opentofu.org`. Terraform resolves it against `registry.terraform.io`, so under Terraform the source needs the `registry.opentofu.org/` prefix unless the provider has also been published there.

## About this fork

Forked from [`terraform-community-providers/terraform-provider-railway`](https://github.com/terraform-community-providers/terraform-provider-railway) via [`dpkass/terraform-provider-railway`](https://github.com/dpkass/terraform-provider-railway).

Upstream stopped at 0.6.2 in April 2026 and models a volume as a nested attribute of `railway_service` with a read-only size. That makes a volume's size unsettable and its lifecycle inseparable from the service it is attached to, so there is no way to express — or to read from a plan — whether a change preserves stored data or destroys it. The fork chain replaced that with first-class `railway_volume` and `railway_volume_instance` resources, moved per-environment deploy configuration onto `railway_service_instance`, and added the storage bucket resources.

Versioning restarts at 0.1.0 here rather than continuing the upstream numbering, because the version line is this repository's, not upstream's. See [`CHANGELOG.md`](./CHANGELOG.md) for what 0.1.0 contains.

## Development

```bash
make testacc          # acceptance tests; needs TF_ACC and a workspace-scoped RAILWAY_TOKEN
make download-schema  # refresh schema.graphql from Railway's GraphQL API
```

Acceptance tests create real Railway projects. Railway allows one project per 30 seconds per workspace, so a full run needs either serialisation or retries; unit tests have no such requirement and run clean.

Releases are cut by pushing a `v*` tag, which runs [`.github/workflows/release.yml`](./.github/workflows/release.yml): goreleaser builds every platform, writes `_SHA256SUMS` and a manifest, and GPG-signs the checksums using the `GPG_PRIVATE_KEY` and `PASSPHRASE` repository secrets.
