# Databasify Engine Template

- [Introduction](#introduction)
- [Creating Your Engine](#creating-your-engine)
- [Choosing the Image](#choosing-the-image)
- [Testing](#testing)
- [Publishing](#publishing)
- [Swapping In Your Database](#swapping-in-your-database)
- [Writing the Driver](#writing-the-driver)

## Introduction

A template for the images of a new [Databasify](https://github.com/databasifylabs/databasify) database engine. Each image is the official image of your database, its tools, and `databasify`, the [engine SDK](https://github.com/databasifylabs/engine-sdk) 1.0 runtime that runs as PID 1. As it is, it runs SQLite on the official `debian` slim image as a stand-in database and passes the SDK's smoke test and runtime conformance suite. The [example database plugin](https://github.com/databasifylabs/plugin-example-database) is the matching driver: a small SQLite engine with its own scripts, built for you to copy.

There is no engine code here. Engine behaviour is a PHP driver with Blade scripts in Databasify, which runs those scripts inside the container through the runtime. This repository changes only when your database or one of its tools moves to a new version.

## Creating Your Engine

Create a repository from this template, then rename `example` to your engine id everywhere it appears:
- `images/example/`;
- the keys and `engine` fields in `versions.json`;
- `DATABASIFY_ENGINE` and the `io.databasify.engine` label in the Dockerfile;
- `engine: example` in both workflows.

## Choosing the Image

`versions.json` pins the official image of your database, one entry per version line.
- A key is `<id>-<line>`.
- Set `registry` to `ghcr.io/<owner>`.
- Keep the packages the tool contract needs: `bash`, `ps`, `ss`, `ip`, `jq`, `openssl` and `curl`. Add the backup and migration tools your driver's scripts call.
- Run the database as a user with UID/GID 64999. The template creates `example`; an official database image usually has its own user, which you remap with `groupmod` and `usermod`.
- Keep the `AS sdk` stage pinned to `ghcr.io/databasifylabs/engine-sdk:<sdk>@sha256:<digest>`, the `ENV DATABASIFY_ENGINE/RELEASE/LINE` block and `ENTRYPOINT ["/usr/local/bin/databasify", "serve"]`. The binary runs as the container runtime only when `DATABASIFY_ENGINE` is set.
- Add a `go` pin only if the image builds a tool from Go source (`GO_IMAGE`).

The refresh pins every package and moves each pin to the newest patch of its line:

```bash
$IMAGES refresh <id>-<line>
```

The line follows the base image's version: `$IMAGES check` requires the base tag to start with `release` and stay inside `line`. The template's line is `13` because its base is Debian 13; with your database's official image it is that database's line.
`$IMAGES` is `tools/images.sh` in a clone of [engine-sdk](https://github.com/databasifylabs/engine-sdk) at the tag the `sdk` pin names (`v1.0.0`), run from this repository's root.

## Testing

```bash
$IMAGES check
$IMAGES build example-13 --load --tag smoke:example-13
$IMAGES smoke smoke:example-13 example
```

The smoke test checks `engine info` and the tool contract, starts the image with `images/example/smoke/processes.json`, runs `images/example/smoke/check.sh` inside it, and then runs the SDK's runtime conformance suite.

## Publishing

Push to GitHub. On every push, `.github/workflows/images.yml`:
1. checks the pins;
2. builds, smoke-tests and scans each line for `linux/amd64` and `linux/arm64`;
3. publishes `ghcr.io/<owner>/databasify-<id>:<release>-<sha7>`. The default branch also moves `<release>`, `<line>` and `latest`.

```bash
git push origin main
```

A new package starts private. Make it public, because Databasify pulls images without credentials.

## Swapping In Your Database

Point `base` at your database's official image, remap its user to 64999, list its tools under `packages`, and change `smoke/processes.json` and `smoke/check.sh` to start and check it. The engine's behaviour stays out of this repository: its driver and scripts live in your plugin, as they do in [`databasifylabs/plugin-example-database`](https://github.com/databasifylabs/plugin-example-database).

## Writing the Driver

The engine reaches Databasify through its driver: a class extending `App\Engines\Drivers\AbstractEngineDriver` that declares the engine's facts (ports, reserved names, password encryptions, host profile) per version line and hands out the groups it offers. Each group method runs one Blade script inside the container. Follow the app's [`adding-database-engine`](https://github.com/databasifylabs/databasify/tree/0.x/.claude/skills/adding-database-engine) skill, which points to the closest existing engine to copy.

Databasify runs only images built on an SDK inside its supported range (`^0.2.2 || ^1.0`).
