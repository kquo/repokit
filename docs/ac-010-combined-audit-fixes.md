# AC-010 Combined Audit Fixes

## Objective Fit

1. Closes four open audit items from plan.md in a single pass
2. Fixes doc drift, build validation gaps, and release tool UX issues found during combined review
3. Improves safety and correctness without changing the workflow model
4. Direct roadmap work

## Summary

Combined audit across adopt/enhance safety, build validation, release git behavior, and docs drift produced 8 findings. Adopt/enhance safety is clean. The remaining findings are: build tool treats `go fmt`/`go fix` failures as non-blocking, release tool has no recovery guidance after partial git failure, docs drift in three locations, and test coverage gaps in reltool/buildtool.

## Findings

### Critical / High

**F1. `go fmt` and `go fix` failures are non-blocking** (`internal/buildtool/buildtool.go:88,95`)
- Both use `runCapturedSoft` which treats failures as informational
- A formatting or API compatibility issue could pass the build silently

**F2. Release git steps have no recovery path** (`internal/reltool/reltool.go:96-110`)
- If `git tag` succeeds but `git push` fails, the local tag is orphaned
- Re-running the release fails with "tag already exists"
- No rollback, no recovery suggestion in error output

**F3. Release message doc says "ideally" but code enforces hard limit** (`docs/build-release.md:88`)
- Doc: "ideally 60 characters or fewer"
- Code: hard rejection at 60 chars
- Overlay template already correct ("60 characters or fewer" without "ideally")

### Medium

**F4. `--apply` flag missing from bootstrap-model recommended arguments** (`docs/bootstrap-model.md:174-189`)
- Flag is implemented and mentioned in enhance mode description but omitted from the argument reference

**F5. Adopt mode conditional metadata requirements undocumented** (`docs/bootstrap-model.md:196-198`)
- When type is specified for adopt, stack (CODE) or publishing-platform+style (DOC) become required
- Not mentioned in docs

### Low / Informational

**F6. voice.md/calendar.md alternatives listed as current** (`docs/bootstrap-model.md:156-157`)
- Ownership model lists "style.md or voice.md" and "content-plan.md or calendar.md"
- Only the first option in each pair exists in overlays

**F7. Adopt/enhance safety: no findings**
- Proposals never overwrite live files
- Manifest parsing handles corruption gracefully
- Section patching preserves all existing content
- Three-way comparison conservative on missing data

**F8. Missing test coverage in reltool and buildtool**
- `Run()` untested in both packages (subprocess-dependent, known coverage ceiling)
- `ensureGitRepo()` and `runGit()` untested in reltool
- `runCapturedSoft()` behavior untested in buildtool

## In Scope

- F1: Make `go fmt` build-breaking; keep `go fix` as soft (it's advisory by design) but document the choice
- F2: Add recovery guidance to release error output when a git step fails after prior steps succeeded
- F3: Change "ideally 60 characters or fewer" to "must be 60 characters or fewer" in `docs/build-release.md`
- F4: Add `--apply` to the recommended arguments list in `docs/bootstrap-model.md`
- F5: Document conditional metadata requirements for adopt mode in `docs/bootstrap-model.md`
- F6: Clarify that voice.md and calendar.md are planned variants, not current overlay content
- F8: Add test for the new `go fmt` non-empty-output failure branch in buildtool, and test `ensureGitRepo` returns error when run outside a git work tree

## Out Of Scope

- Atomic git transactions or automatic rollback in the release tool (too complex for the benefit; recovery guidance is sufficient)
- Full `Run()` integration tests for buildtool and reltool (subprocess-dependent ceiling documented in AC-003)
- Adding voice.md or calendar.md overlays (separate plan.md item)

## Implementation Notes

- F1: `go fmt` exits 0 even when it rewrites files, reporting changed filenames on stdout. The fix is not a helper-function swap — it requires new logic: run `go fmt` via `runCapturedSoft` (which captures stdout), then check whether the output is non-empty. If non-empty, treat it as a build failure (print the changed files and return an error). `go fix` stays soft with a code comment explaining the advisory-only design choice.
- F2: When a git step fails in reltool `Run()`, include context about what succeeded and what the user should do to recover. For example: "git push tag failed; the tag exists locally — retry with `git push origin <tag>` or delete with `git tag -d <tag>`"
- F3: One-word change in build-release.md
- F4: Add `-a, --apply` line to the recommended arguments block
- F5: Add a note after the adopt requirements that overlay-specific metadata is required when type is specified
- F6: Change "style.md or voice.md" to "style.md (voice.md planned)" and similarly for calendar.md
- F8: Test the `go fmt` non-empty-output failure branch by calling the format-check logic with content that needs formatting and verifying it returns an error. Test `ensureGitRepo` by running it in a temp dir that is not a git repo.

## Acceptance Tests

- [Automated] `go fmt` failure is build-breaking (non-empty output from format check returns error)
- [Automated] `ensureGitRepo` returns error when run outside a git work tree
- [Automated] Release message limit doc says "must be" not "ideally"
- [Manual] Release tool error output includes recovery guidance after partial failure
- [Manual] `docs/bootstrap-model.md` lists `--apply` in recommended arguments
- [Manual] `docs/bootstrap-model.md` documents conditional adopt metadata requirements
- [Manual] `docs/bootstrap-model.md` ownership model clarifies planned vs current file variants

## Documentation Updates

- `docs/build-release.md` — fix "ideally" to "must be"
- `docs/bootstrap-model.md` — add `--apply` to arguments, document adopt conditionals, clarify planned variants

## Status

COMPLETE
