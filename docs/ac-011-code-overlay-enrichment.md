# AC-011 CODE Overlay Enrichment

## Objective Fit

1. Generated CODE repos get more actionable release/upgrade guidance and a concrete AC example out of the box
2. Reduces the gap between the rich docs in repokit itself and the thinner docs generated repos start with
3. These are the last two CODE overlay priorities in plan.md
4. Direct roadmap work

## Summary

Enrich the CODE overlay with two improvements: (1) add upgrade and refresh guidance to `build-release.md.tmpl` so generated repos know how to check for template updates, and (2) add a concrete example AC doc alongside the existing `ac-template.md.tmpl` so generated repos have a worked example of the AC pattern, not just a blank template.

## In Scope

- Add a "Template Upgrade" section to `overlays/code/files/docs/build-release.md.tmpl` covering: how to check the current template version via `TEMPLATE_VERSION`, how to compare the generated repo against the source template manually, and what `.repokit-manifest` records for future tooling-assisted comparison. This section describes the operator's review workflow, not a generated-repo command — enhance mode is a template-maintainer tool, not a generated-repo refresh command.
- Add a "Release Artifacts" section to `overlays/code/files/docs/build-release.md.tmpl` documenting `TEMPLATE_VERSION` as a template-owned artifact, and `CHANGELOG.md` conditionally ("if the repo uses a changelog"). The overlay does not generate `CHANGELOG.md` — repos opt into it.
- Create `overlays/code/files/docs/ac-example.md.tmpl` — a concrete filled-in AC example showing the pattern in use (not a real AC, just a teaching artifact)
- Regenerate rendered examples under `examples/code/docs/`
- Update `overlays/code/README.md` to list the new template file

## Out Of Scope

- Changes to repokit's own `docs/build-release.md` (already has this content)
- Changes to the AC template itself (`ac-template.md.tmpl`)
- DOC overlay enrichment (separate plan.md items)
- Changes to bootstrap behavior or Go code

## Implementation Notes

- The "Template Upgrade" section should be generic (no repokit-specific paths) and describe the operator review workflow: check `TEMPLATE_VERSION` against the source template, diff files manually, and note that `.repokit-manifest` records bootstrap-time checksums for future tooling-assisted comparison. Do not describe enhance mode as a generated-repo command — it runs from the template repo, not the generated repo.
- The example AC should use a realistic but fictional scenario (e.g. "AC-001 Add user authentication") to demonstrate each section's purpose. Include a note at the top that it is an example, not an active AC.
- The "Release Artifacts" section lists `TEMPLATE_VERSION` as always present (written by bootstrap) and `CHANGELOG.md` as conditional ("if the repo maintains a changelog"). Does not generate `CHANGELOG.md`.

## Acceptance Tests

- [Automated] Bootstrap `new` mode for CODE produces `docs/ac-example.md`
- [Automated] Bootstrap `new` mode for CODE produces `docs/build-release.md` containing "Template Upgrade" section
- [Automated] Bootstrap `new` mode for CODE produces `docs/build-release.md` containing "Release Artifacts" section
- [Automated] Bootstrap `adopt` mode proposes `docs/ac-example.md` when it already exists
- [Automated] Bootstrap `adopt` mode proposes `docs/build-release.md` with updated content when it already exists
- [Manual] Generated `ac-example.md` is clearly marked as an example and demonstrates all AC sections
- [Manual] Overlay README lists the new file

## Documentation Updates

- `overlays/code/README.md` — list `docs/ac-example.md`

## Status

COMPLETE
