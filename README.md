# repokit
Template repo that bootstraps and adopts governed repositories, and maintains itself through enhance mode. Built from:

- a common base contract in `base/`
- a repo-type overlay in `overlays/code/` or `overlays/doc/`
- a deterministic Go bootstrap command that renders concrete files into a target repo

## Why
Most AI-assisted repo work fails not because the model is weak, but because the collaboration contract is implicit, inconsistent, and hard to reproduce. `repokit` exists to make that contract explicit by providing a governance and workflow template for deterministic human-AI project collaboration.

It gives a repo a stable structure for how work is proposed, reviewed, documented, and maintained, so the human and the agent are operating from the same visible rules instead of hidden session context. The goal is not more process for its own sake; it is less coordination drift, less prompt-only memory, and more repeatable project maintenance.

## Modes

### `new` and `adopt`
Consumer modes, run from a target repo or empty directory. Repokit is read-only source.

**`new`** — bootstrap an empty or near-empty directory into a governed `CODE` or `DOC` repo.

```bash
go run <template-root>/cmd/bootstrap \
  -m new -y CODE \
  -n my-service \
  -p "API gateway for internal services" \
  -s "Go CLI"
```

```bash
go run <template-root>/cmd/bootstrap \
  -m new -y DOC \
  -n my-docs \
  -p "Public developer documentation" \
  -u "Static site generator" \
  -v "Clear, factual, concise"
```

**`adopt`** — apply governance to an existing repo with conservative behavior: fit assessment, proposal files instead of overwrites, and section-level `AGENTS.md` patching that adds only missing governed sections.

```bash
go run <template-root>/cmd/bootstrap \
  -m adopt \
  -n existing-repo \
  -p "Short project purpose" \
  -s "Go service" \
  -d
```

### `enhance`
Template-maintenance mode, run from inside this repo. The only mode that runs from repokit itself and the only mode that can propose changes back into the template.

Enhance inspects another governed repo by reference path, comparing at the constraint level for governance sections and per-section for structured markdown files. When a `.repokit-manifest` exists in the reference repo, enhance uses three-way comparison to distinguish user customizations from stale template content. With `--apply`, it writes `.template-proposed` files for assisted merge. No template files are overwritten automatically.

```bash
go run ./cmd/bootstrap \
  -m enhance \
  -r <reference-root> \
  -d
```

Run `bootstrap --help` for all flags.

## Design
The target repo stays self-contained. The template repo is read-only at bootstrap time and is not imported as a submodule, package, or runtime dependency. The bootstrap tool is Go-based so the template works across macOS, Linux, and Windows without requiring a specific shell.

## Self-Hosting Status
This repo is itself governed as a `CODE` repo and carries the core artifacts at the root:

- [`AGENTS.md`](AGENTS.md)
- [`arch.md`](arch.md)
- [`plan.md`](plan.md)
- [`CHANGELOG.md`](CHANGELOG.md)
- [`docs/README.md`](docs/README.md)
- [`docs/agent-roles/`](docs/agent-roles/)

## Rendered Examples
Generated examples:

- [`examples/code/`](examples/code/)
- [`examples/doc/`](examples/doc/)

See [`docs/bootstrap-model.md`](docs/bootstrap-model.md).
