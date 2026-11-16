# Changelog

## Unreleased

### Features

- `railway_service_instance` gains `static_egress_enabled` and the computed `static_egress_ips`, so a service's static outbound addresses can be switched on and fed straight into a firewall allowlist such as an AWS security group. The default is `false`, so applying a configuration that omits it releases static egress that was enabled in the dashboard.

### Fixes

- `railway_project` and `railway_environment` creation wait out Railway's limits of one new project, and one new environment, per workspace every 30 seconds, retrying for up to 90 seconds, instead of failing when one apply creates several.

## 0.1.0

First release of this fork. Versioning restarts here, so this entry describes the state of the provider rather than replaying the upstream history that produced it.

### Resources

- **Project and environment** — `railway_project`, `railway_environment`
- **Services** — `railway_service` (project-wide identity), `railway_service_instance` (per-environment deploy configuration: image or repo source, registry credentials, regions and replica counts, `vcpus`, `memory_gb`, healthcheck, start command, cron schedule)
- **Storage** — `railway_volume` (project-scoped), `railway_volume_instance` (per-environment mount path, service attachment, and a provisioned `size_mb` that can be raised)
- **Buckets** — `railway_bucket`, `railway_bucket_instance`, `railway_bucket_cors_configuration`, plus the `railway_bucket_credentials` ephemeral resource
- **Networking** — `railway_custom_domain`, `railway_service_domain`, `railway_tcp_proxy`
- **Variables** — `railway_variable`, `railway_variable_collection`, `railway_shared_variable`

### Notable differences from upstream 0.6.2

- Volumes are their own resources. Upstream modelled a volume as a nested attribute of `railway_service` with a computed, read-only size, which left size unsettable and gave a plan no way to distinguish moving a mount from discarding the data behind it. Size is now settable and growable, and only replacing the `railway_volume` destroys data.
- Deploy configuration is per-environment via `railway_service_instance`, which also exposes resource sizing and healthcheck settings that had no equivalent before.
- `regions` differs in shape between the two: a list of `{region, num_replicas}` on `railway_service`, a map keyed by region name on `railway_service_instance`.
- Bucket resources are new.

### Deprecations

Every deploy attribute on `railway_service` — `source_image`, `source_repo`, `source_repo_branch`, `root_directory`, `config_path`, `cron_schedule`, `regions`, `volume`, and the `source_image_registry_*` credentials — is deprecated in favour of `railway_service_instance` and the volume resources. They still work, so configurations written against upstream continue to plan and apply while they are migrated.
