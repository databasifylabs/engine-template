---
name: building-an-engine
description: "Turn this engine template into a working Databasify engine: fill in the Go package (required Nodes and Cluster, optional groups and their Unsupported stand-ins, capabilities, host profile, settings catalog), adapt the Dockerfile and versions.json to the real database, and check it locally with the SDK's build, smoke and conformance tools. Use when building or extending an engine from this template. Registration in Databasify is a separate plugin repository."
---

# Building an engine

This repository is a Databasify engine on the public [engine SDK](https://github.com/databasifylabs/engine-sdk). `cmd/databasify/main.go` runs `sdk.Run(engine.New())`, and that binary is PID 1 in every container of the engine. `README.md` covers the whole path from the first commit to a published image; this skill is the working checklist.

## The Go package (`engine/`)

`contract.Engine` has `Name()`, `Capabilities()`, `HostProfile()` and one accessor per group.

- **`Name()`** is the engine id. It must equal `DATABASIFY_ENGINE` and the `io.databasify.engine` label in the Dockerfile, the `engine` field in `versions.json`, and the id the plugin registers.
- **`Nodes` and `Cluster` are required.** The app's create, add-node, remove-node, restart, upgrade and recovery flows call them on every engine.
  - The stubs implement `contract.Partial`: `Implements()` lists the commands they offer, and every other command answers `unsupported`.
  - Implement a command, add it to `Implements()`, and add a test.
  - Delete `Implements()` once the whole group works.
  - `Ready`, `Cluster.Enable` and `Cluster.Rejoin` are long commands: they report phases through `report`, and the app polls `operation status`.
- **Optional groups** (`Accounts`, `Archive`, `Migration`, `Pools`, `Config`) return their `contract.Unsupported*` stand-in until the engine fully supports them, and the app hides those features. Replace the stand-in with your own type when you add one.
- **`Capabilities()`** lists:
  - the ports: `client`, `pool` (or none) and `peers`, which open in each node's firewall;
  - reserved users and databases;
  - password encryptions and the default one.
- **`HostProfile()`** holds the sysctl settings, floors and the transparent hugepages choice that `host tune` applies on the server. The zero value means no tuning.
- **Settings:** a `Config` group serves `config fields`, `show` and `tune` from a catalog built with the SDK's `settings` package (`settings.Integer`, `Choice`, `Size`, `.Restarting()`, `.TunedFromServer()`). Offer the settings DigitalOcean Managed Databases offers for that database, no more and no fewer. Values tuned from the server's CPU and memory stay automatic.
- **Processes:** write `/srv/<engine>/conf/processes.json` for the SDK's supervisor (`manifest` package). Secrets arrive only on stdin; write them as 600 files owned by the engine user, never into argv or the manifest.
- **Every command is safe to repeat.** The app retries a lost operation under a new key.

## The image (`images/<engine>/`, `versions.json`)

- Rename `images/example/` and the `example-9` entry after your engine. A key is `<engine>-<line>`.
- Point `base` at the official image of your database: its tag and digest.
- Set `registry` to `ghcr.io/<owner of this repository>`.
- Keep the packages the tool contract needs: `bash ps ss ip jq openssl curl`.
- Remap the database user to UID/GID 64999, as the Dockerfile does for `valkey`.
- Keep the labels `io.databasify.{engine,release,line,sdk}`, `ENV DATABASIFY_ENGINE/RELEASE/LINE` and `ENTRYPOINT ["/usr/local/bin/databasify", "serve"]`.
- `smoke/processes.json` starts the database. `smoke/check.sh` waits for it, checks `$DATABASIFY_RELEASE` and every bundled tool, and runs inside the container. Keep the empty `conformance` file so the smoke test runs the SDK's conformance suite.
- Add a line by copying an entry and running `$IMAGES refresh <key>`, which pins the packages and moves everything to the newest patch of the line.

## Check it locally

`$IMAGES` stands for `tools/images.sh` in a clone of [engine-sdk](https://github.com/databasifylabs/engine-sdk) at the SDK version `go.mod` requires. Run it from this repository's root:

```bash
go vet ./... && go test -race -count=1 ./...
$IMAGES check
$IMAGES build <engine>-<line> --load --tag smoke:<engine>-<line>
$IMAGES smoke smoke:<engine>-<line> <engine>
```

The conformance suite always tests `engine capabilities`, `engine host-profile`, `engine info` and `node status`. It tests every other command once `engine capabilities` lists it, and checks that every unlisted command answers `unsupported`. Passwords must never appear in a reply.

## Registering it

The engine reaches Databasify through a plugin made from [plugin-example-database](https://github.com/databasifylabs/plugin-example-database):
- declare the engine id in its `composer.json` under `extra.databasify.engines`;
- point `image()` at `<registry>/databasify-<engine>`;
- describe the connection with `connection()`.

Databasify shows it as a community engine, and the admin confirms a warning before installing it.
