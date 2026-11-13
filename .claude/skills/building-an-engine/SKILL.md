---
name: building-an-engine
description: "Turn this template into the images of a new Databasify engine: point versions.json and the Dockerfile at the real database and its tools, keep the engine-sdk 1.0 runtime stage, and check it locally with the SDK's check, build, smoke and runtime conformance tools. Use when building or extending an engine's images from this template. The engine's behaviour (its PHP driver and Blade scripts) lives in the Databasify app, not here."
---

# Building an engine's images

This repository builds the images of one Databasify engine. Each image is the official image of the database, its tools, and `databasify`, the [engine SDK](https://github.com/databasifylabs/engine-sdk) 1.0 runtime, as PID 1. There is no engine code here: the engine's PHP driver and Blade scripts live in Databasify (the app's `adding-database-engine` skill). `README.md` covers the whole path; this skill is the working checklist.

## The image (`images/<engine>/`, `versions.json`)

- **Rename** `images/example/`, the `example-9` entry and `engine: example` in both workflows after your engine. A key is `<engine>-<line>`. The engine id must equal `DATABASIFY_ENGINE`, the `io.databasify.engine` label and the id the app registers.
- **Base:** point `base` at the official image of your database, its tag and digest. Set `registry` to `ghcr.io/<owner of this repository>`.
- **Tools:** keep the packages the tool contract needs (`bash ps ss ip jq openssl curl`) and add every tool the driver's scripts and jobs call (backup, restore, migration). A tool built from Go source needs the `go` pin and a stage `FROM --platform=$BUILDPLATFORM ${GO_IMAGE}`.
- **User:** remap the database user to UID/GID 64999, as the Dockerfile does for `valkey`.
- **Runtime:** keep the `FROM ghcr.io/databasifylabs/engine-sdk:<sdk>@sha256:<digest> AS sdk` stage and `COPY --from=sdk /databasify /usr/local/bin/databasify`. The `sdk` pin in `versions.json` must name the same version; `$IMAGES check` refuses a mismatch.
- **Keep:** the labels `io.databasify.{engine,release,line,sdk}`, `ENV DATABASIFY_ENGINE/RELEASE/LINE` (without `DATABASIFY_ENGINE` the binary acts as the host binary and the container never starts) and `ENTRYPOINT ["/usr/local/bin/databasify", "serve"]`.
- **Smoke:** `smoke/processes.json` starts the database under the supervisor. `smoke/check.sh` waits for it, checks `$DATABASIFY_RELEASE` and every bundled tool, and runs inside the container with `bash -euo pipefail -s`. Keep the empty `conformance` file so the smoke test runs the SDK's runtime conformance suite.
- **Lines:** add a line by copying an entry and running `$IMAGES refresh <key>`, which pins the packages and moves everything to the newest patch of the line.

## Check it locally

`$IMAGES` stands for `tools/images.sh` in a clone of [engine-sdk](https://github.com/databasifylabs/engine-sdk) at the tag the `sdk` pin names (`v1.0.0`). Run it from this repository's root:

```bash
$IMAGES check
$IMAGES build <engine>-<line> --load --tag smoke:<engine>-<line>
$IMAGES smoke smoke:<engine>-<line> <engine>
```

Pushing a branch runs the same checks in CI on amd64 and arm64, plus a Trivy scan, and publishes only `<release>-<sha7>` tags.

## The driver

The engine's behaviour is a PHP driver extending `App\Engines\Drivers\AbstractEngineDriver` in Databasify: its facts per line, the groups it offers (`nodes`, `cluster` and `config` required; `accounts`, `archive`, `migration` and `pools` optional), and one Blade script per action. Follow the app's `adding-database-engine` skill. Never add Go code to this repository: the SDK 1.0 runtime is frozen, and every engine feature is a script or a job in the app.
