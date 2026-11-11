# Databasify Engine Template

- [Introduction](#introduction)
- [Creating Your Engine](#creating-your-engine)
- [Writing the Engine](#writing-the-engine)
- [Choosing the Image](#choosing-the-image)
- [Testing](#testing)
- [Publishing](#publishing)
- [Registering the Engine](#registering-the-engine)

## Introduction

A template for a new [Databasify](https://github.com/databasifylabs/databasify) database engine, built on the [engine SDK](https://github.com/databasifylabs/engine-sdk). It passes the SDK's conformance suite as it is. It runs the official `valkey` image as a stand-in database, answers `node status`, and says `unsupported` to everything else.

An engine needs two repositories: this one, which builds and publishes the images, and a plugin made from [plugin-example-database](https://github.com/databasifylabs/plugin-example-database), which registers the engine in Databasify. Databasify shows a plugin's engine as a community engine, built and supported by you.

## Creating Your Engine

Create a repository from this template, then rename `example` to your engine id everywhere it appears:
- `Name()` in `engine/engine.go`;
- `images/example/`;
- the keys and `engine` fields in `versions.json`;
- `DATABASIFY_ENGINE` and the labels in the Dockerfile;
- `engine: example` in both workflows.

Also set the module path in `go.mod`:

```bash
go mod edit -module github.com/<owner>/engine-<id>
```

## Writing the Engine

`engine/` implements the SDK's `contract.Engine`.
- `Nodes` and `Cluster` are required. Their stubs implement `contract.Partial`: add each command to `Implements()` as you implement it, and delete `Implements()` once the whole group works.
- Every optional group (`Accounts`, `Archive`, `Migration`, `Pools`, `Config`) returns its `Unsupported` stand-in until you replace it.
- Fill in `Capabilities()` (ports, reserved names, password encryptions) and `HostProfile()` (host tuning).

```go
func (Nodes) Implements() []string {
	return []string{"node status", "node configure", "node start"}
}
```

## Choosing the Image

`versions.json` pins the official image of your database, one entry per version line.
- A key is `<id>-<line>`.
- Set `registry` to `ghcr.io/<owner>`.
- Keep the packages the tool contract needs: `bash`, `ps`, `ss`, `ip`, `jq`, `openssl` and `curl`.

The refresh pins every package and moves each pin to the newest patch of its line:

```bash
$IMAGES refresh <id>-<line>
```

`$IMAGES` is `tools/images.sh` in a clone of [engine-sdk](https://github.com/databasifylabs/engine-sdk) at the version `go.mod` requires, run from this repository's root.

## Testing

```bash
go vet ./... && go test -race -count=1 ./...
$IMAGES check
$IMAGES build example-9 --load --tag smoke:example-9
$IMAGES smoke smoke:example-9 example
```

The smoke test starts the image, runs `images/example/smoke/check.sh` inside it, and then runs the conformance suite.
- The suite always tests `engine capabilities`, `engine host-profile`, `engine info` and `node status`.
- It tests every other command once `engine capabilities` lists it.
- It checks that every unlisted command answers `unsupported`.

## Publishing

Push to GitHub. On every push, `.github/workflows/images.yml` does three things:
1. It runs the tests.
2. It builds, smoke-tests and scans each line for `linux/amd64` and `linux/arm64`.
3. It publishes `ghcr.io/<owner>/databasify-<id>:<release>-<sha7>`. The default branch also moves `<release>`, `<line>` and `latest`.

```bash
git push origin main
```

A new package starts private. Make it public, because Databasify pulls images without credentials.

## Registering the Engine

In a plugin made from [plugin-example-database](https://github.com/databasifylabs/plugin-example-database), declare the engine id in `composer.json` and point `image()` at your images:

```json
"extra": { "databasify": { "engines": ["<id>"] } }
```

Databasify runs only images built on an SDK inside its supported range. Rebuild on a new SDK when the range moves, or databases on your engine become read-only.
