# AC-005 Development Knowledge Layer

## Objective Fit

1. Any coding agent working in a governed CODE repo gets durable engineering guidance without depending on agent-specific memory
2. This fills the gap between governance (AGENTS.md), process (development-cycle.md), and validation (build-release.md) — none of those carry reusable coding knowledge
3. repokit is agent-agnostic; knowledge must live in repo artifacts, not in any single agent's memory system
4. Direct roadmap work — strengthens the CODE overlay template

## Summary

Add two new template artifacts to the CODE overlay: `docs/development-guidelines.md` for durable coding guidance any agent should follow, and a `docs/knowledge/` directory pattern for expandable engineering notes that elaborate on development-guidelines topics. Seed the guidelines template with genuinely reusable engineering notes (not process rules). Seed the knowledge directory with an index and example structure.

## In Scope

- Create `overlays/code/files/docs/development-guidelines.md.tmpl` with reusable engineering guidance sections
- Create `overlays/code/files/docs/knowledge/README.md.tmpl` as an index for accumulated lessons
- Add repokit's own `docs/development-guidelines.md` with real entries relevant to this project
- Create repokit's own `docs/knowledge/README.md` as index
- Seed repokit's `docs/knowledge/` with at least one real knowledge entry (e.g. propagation patterns for template repos)
- Regenerate `examples/code/` to include both new artifacts
- Update `overlays/code/README.md` to list the new files

## Out Of Scope

- Changes to AGENTS.md, development-cycle.md, or build-release.md
- DOC overlay changes (stretch only, same as AC-004)
- Migrating existing Claude memory entries into repo knowledge — memory stays as optional acceleration
- Defining a formal schema or tooling for knowledge files — keep it simple markdown

## Implementation Notes

- The doc layer hierarchy after this AC:
  - `AGENTS.md` — governance and boundaries
  - `docs/development-cycle.md` — workflow and process
  - `docs/build-release.md` — validation and release flow
  - `docs/development-guidelines.md` — durable coding guidance for any agent
  - `docs/knowledge/` — deeper supporting notes, examples, and reusable lessons that expand development-guidelines topics
- `development-guidelines.md` should contain guidance classes like:
  - identifier strategy and key design
  - schema or data migration pitfalls
  - external integration reconciliation patterns
  - generated artifact propagation rules
  - error handling and validation boundaries
  - testing expectations beyond what build-release.md covers
- The template version should be generic enough to seed any CODE repo; repokit's own version will have project-specific entries
- `docs/knowledge/` entries should be one file per topic, with a README index
- Agent-specific memory remains optional acceleration on top of these repo artifacts

## Acceptance Tests

- [Manual] `overlays/code/files/docs/development-guidelines.md.tmpl` exists with reusable guidance sections
- [Manual] `overlays/code/files/docs/knowledge/README.md.tmpl` exists as a knowledge index template
- [Manual] repokit's `docs/development-guidelines.md` has real entries relevant to repokit and representative of reusable CODE-repo guidance
- [Manual] repokit's `docs/knowledge/README.md` exists as an index
- [Manual] repokit's `docs/knowledge/` has at least one real knowledge entry
- [Manual] `examples/code/docs/development-guidelines.md` matches the template shape
- [Manual] `examples/code/docs/knowledge/README.md` matches the template shape
- [Manual] `overlays/code/README.md` lists both new files
- [Manual] No duplication with AGENTS.md, development-cycle.md, or build-release.md

## Documentation Updates

- `overlays/code/README.md` — add new file entries
- `docs/README.md` — add entries for development-guidelines.md and knowledge/

## Status

COMPLETE
