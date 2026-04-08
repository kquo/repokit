# repokit

Template repo for generating governed repositories from:

- a common base contract in `base/`
- a repo-type overlay in `overlays/code/` or `overlays/doc/`
- a deterministic Go bootstrap command that renders concrete files into a target repo

## Quick Start

Bootstrap a new **CODE** repo:

```bash
go run <template-root>/cmd/bootstrap \
  -m new -y CODE \
  -n my-service \
  -p "API gateway for internal services" \
  -s "Go CLI"
```

Bootstrap a new **DOC** repo:

```bash
go run <template-root>/cmd/bootstrap \
  -m new -y DOC \
  -n my-docs \
  -p "Public developer documentation" \
  -u "Static site generator" \
  -v "Clear, factual, concise"
```

Adopt an existing repo:

```bash
go run <template-root>/cmd/bootstrap \
  -m adopt \
  -n existing-repo \
  -p "Short project purpose" \
  -s "Go service" \
  -d
```

Review another governed repo for template improvements:

```bash
go run <template-root>/cmd/bootstrap \
  -m enhance \
  -r <reference-root> \
  -d
```

Run `--help` for all flags:

```bash
go run <template-root>/cmd/bootstrap --help
```

## Intended Use

This repo is a reference frame used from inside a target working directory.

Three modes are supported:

- `new`: bootstrap an empty or near-empty folder into a governed `CODE` or `DOC` repo
- `adopt`: apply the methodology to an existing repo with conservative proposal behavior, fit assessment, and section-level patching for `AGENTS.md`
- `enhance`: inspect another governed repo for portable methodology improvements and create an AC doc for the highest-priority actionable candidate

## Operating Model

The target repo stays self-contained.
The template repo is read-only at bootstrap time and is not imported as a submodule, package, or runtime dependency.

The bootstrap tool is Go-based so the template works across macOS, Linux, and Windows without requiring a specific shell.

The flow is:

1. user opens a coding agent in the target directory
2. user gives the agent the absolute path to this template repo
3. agent runs the bootstrap command from this template repo, targeting the current repo
4. agent inspects the target repo state
5. agent chooses bootstrap mode: `new` or `adopt`
6. agent gathers required inputs
7. agent writes concrete files into the target repo
8. generated repo records its template marker and becomes independently managed

## Operator Guide

Use `new` when the target directory is empty or nearly empty and you want a full rendered baseline.

Use `adopt` when the target repo already exists and you want conservative behavior: fit assessment, proposal files instead of overwrites, and section-level patching for `AGENTS.md` that adds only missing governed sections.

Use `enhance` only from inside this template repo to inspect another governed repo for portable improvements. Enhance compares at the constraint level for governance sections and per-section for structured markdown files. When a `.repokit-manifest` exists in the reference repo, enhance uses three-way comparison to distinguish user customizations from stale template content. With `--apply`, it writes `.template-proposed` files for assisted merge. No template files are overwritten automatically.

## Self-Hosting Status

This repo carries the core `CODE`-repo artifacts at the root:

- [`AGENTS.md`](AGENTS.md)
- [`arch.md`](arch.md)
- [`plan.md`](plan.md)
- [`CHANGELOG.md`](CHANGELOG.md)
- [`docs/README.md`](docs/README.md)
- [`docs/agent-roles/`](docs/agent-roles/)

The template repo is governed as a `CODE` repo and uses `enhance` for self-maintenance.

## Rendered Examples

Generated examples:

- [`examples/code/`](examples/code/)
- [`examples/doc/`](examples/doc/)

See [`docs/bootstrap-model.md`](docs/bootstrap-model.md).
