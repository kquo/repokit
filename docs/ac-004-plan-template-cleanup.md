# AC-004 Plan Template Cleanup

## Objective Fit

1. Maintainers and generated repos get a plan.md that separates forward-looking work from process guidance, reducing drift and redundancy
2. This is the natural next template-quality step after AC-001/002/003 stabilized the workflow and docs
3. development-cycle.md and build-release.md already exist as the correct homes for process text
4. Direct roadmap work — improves a core template artifact

## Summary

Simplify `plan.md.tmpl` to six sections (Product Direction, Current Platform, Objective-Fit Rubric, Priorities, Deferred, Constraints). Remove process guidance that belongs in `docs/development-cycle.md`. Clean up repokit's own `plan.md` to match: drop completed items, drop redundant Notes and Improvement Intake sections. Optionally add Deferred/Constraints to the DOC overlay's `content-plan.md.tmpl`.

## In Scope

- Update `overlays/code/files/plan.md.tmpl` to new six-section shape
- Clean up repokit's own `plan.md`: remove completed R1/R2, remove Notes, remove Improvement Intake process text, add Deferred table, add Constraints
- Regenerate `examples/code/plan.md`

## Stretch

- Add Deferred/Constraints to `overlays/doc/files/content-plan.md.tmpl` only if there is a clear editorial reason; otherwise leave it simpler

## Out Of Scope

- Changes to AGENTS.md governance sections
- Changes to docs/development-cycle.md or docs/build-release.md (they already have the right content)
- Fixing plan.md in other repos (e.g. skout) — separate effort
- Adding new overlay types or template files
- Product-specific sections like Interaction Principles or Signal Hierarchy — those belong in individual repos, not the template

## Implementation Notes

- Completed priorities belong in CHANGELOG.md and git history, not plan.md
- Process guidance ("capture follow-on improvements here", "use docs/ for AC drafts") already lives in docs/development-cycle.md — do not duplicate
- Deferred table gives ideas a parking lot without cluttering Priorities
- Constraints captures project-specific anti-patterns distinct from AGENTS.md governance rules
- Items in each section should be under 200 characters

## Acceptance Tests

- [Manual] `plan.md.tmpl` has exactly six sections: Product Direction, Current Platform, Objective-Fit Rubric, Priorities, Deferred, Constraints
- [Manual] repokit's `plan.md` has no struck-through items, no Notes section, no Improvement Intake process text
- [Manual] repokit's `plan.md` has a Deferred table and Constraints section with real entries
- [Manual] `examples/code/plan.md` matches the updated CODE template shape
- [Manual] No process/workflow guidance remains in any plan.md — that lives only in docs/development-cycle.md

## Documentation Updates

- overlays/code/README.md if plan.md.tmpl changes are noted there

## Status

COMPLETE
