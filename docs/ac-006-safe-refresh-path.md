# AC-006 Safe Refresh Path Improvements

## Objective Fit

1. Enhance mode becomes more precise at identifying real improvements vs noise, reducing false positives and missed deltas for template maintainers
2. This is the primary R3 priority — the enhance workflow is the main template-maintenance surface and its comparison logic is the weakest link
3. Must preserve the review-first, AC-driven workflow — no auto-apply, no bypassing human review
4. Direct roadmap work

## Summary

Improve enhance mode's comparison precision in two dimensions: governed-section comparison (currently keyword signal matching) and overlay/file comparison (currently whole-file). This is an umbrella AC with phased delivery. Phase 1 targets the two comparison refinements. Later phases cover version-awareness, classifier extensibility, and assisted-apply.

## Phases

### Phase 1 — Comparison Refinement (implementation scope)

**Governed-section comparison:**
- Current: `governanceSectionCovered()` checks for keyword signals ("exploratory", "authorization", "targeted", etc.). Catches broad presence/absence but misses meaningful differences within sections that share the same keywords.
- Goal: Compare at the constraint level, not just signal presence. Two sections that both mention "authorization" but differ on what requires authorization should produce a candidate.

**Overlay/file comparison:**
- Current: `reviewMappedFile()` normalizes and compares entire file content. Any difference flags the whole file.
- Goal: Section-level diffing for structured markdown files (those with `##` headings). When only one section differs, the candidate should identify which section changed, not just "file differs."

### Phase 2 — Version Awareness (implementation scope)

**Problem:** Enhance mode does two-way comparison (template vs reference). It cannot distinguish whether a difference exists because the user intentionally customized a file or because the template evolved after bootstrap. Every difference gets equal treatment.

**Solution:** Write a `.repokit-manifest` file at bootstrap time that records checksums of all generated files and their template sources. During enhance, use this manifest for three-way comparison to classify each difference by origin.

**Bootstrap manifest (`.repokit-manifest`):**
- Written during `new` and `adopt` bootstrap into the generated repo
- Records template version, per-file rendered content checksum, source template path, and source template checksum
- Line-oriented text format (`repokit-manifest-v1`) parseable with stdlib only
- For adopt mode, records canonical checksums (what the template intended) regardless of propose/skip transforms

**Three-way classification during enhance:**

| User changed file? | Template source changed? | ChangeOrigin | Effect |
|---|---|---|---|
| no | no | skip | No candidate produced |
| yes | no | `"user"` | Evaluate normally — these are potential improvements to adopt into template |
| no | yes | `"template"` | Defer — reference is stale, template already evolved past bootstrap baseline |
| yes | yes | `"both"` | Flag for careful review — both sides modified |
| no manifest | — | `""` | Fall back to current two-way comparison (backward compatible) |

Key insight: enhance extracts improvements FROM reference INTO template. User changes are the interesting ones (potential improvements). Template-only changes mean the reference is stale, not that it has improvements.

**Dual checksums resolve the placeholder problem:**
- Rendered checksum (after `{{REPO_NAME}}` substitution) → compare against current reference file to detect user changes
- Source checksum (raw `.tmpl` file) → compare against current template source to detect template changes
- No need to re-render templates with original placeholders

**`planRender` refactoring:**
- Extract `planCanonical()` (pre-adopt-transform operations) and `applyAdoptTransforms()`
- Build manifest from canonical operations (correct baseline regardless of propose/skip)
- Append manifest write to executable operations

**Governed-section handling:**
- Use whole-file AGENTS.md checksum from manifest as pre-filter
- Template-only change on AGENTS.md → skip all section candidates (reference is stale)
- User-only change → evaluate section candidates normally with ChangeOrigin `"user"`
- Both changed → proceed with section-level comparison with ChangeOrigin `"both"`

### Phase 3 — Classifier Extensibility (implementation scope)

**Problem:** Classification logic is spread across three hardcoded locations:
- `classifyEnhancement` — four `if`/`switch` branches that map content traits to portability/disposition/reason tuples
- `projectSpecificMarkers` — two hardcoded heuristics (repo name mention, absolute user paths)
- `sectionSignals` — a six-case switch mapping section names to keyword signal sets

Adding a new overlay type, a new project-specific marker, or a new governance signal requires editing Go code and recompiling. As overlays grow (e.g. API, CLI, library variants), this becomes a maintenance bottleneck.

**Solution:** Replace the three hardcoded locations with a data-driven rule table loaded at startup. Each rule is a declarative struct expressing: match condition → classification outcome. The rule table is defined in Go code (not external config) to preserve compile-time safety, but structured so adding a rule is a one-line table entry rather than a new code branch.

**Classification rule table:**
```go
type classificationRule struct {
    Name        string
    Match       func(ctx classificationContext) bool
    Portability string
    Disposition string
    Reason      string
    Priority    int // lower wins; first matching rule applies
}
```

Where `classificationContext` carries all the inputs currently scattered across the three functions:
```go
type classificationContext struct {
    Content        string
    ReferenceRoot  string
    TemplateTarget string
    Governance     bool
    Area           string
}
```

**Default rules (preserving current behavior):**
1. Project-specific: content mentions reference repo name or contains absolute user paths → `("project-specific", "defer", ...)`
2. Governance: `governance == true` → `("portable", "accept", ...)`
3. Workflow helpers: template target ends in `.go.tmpl`, `.sh.tmpl`, or equals `TEMPLATE_VERSION` → `("portable", "accept", ...)`
4. Default fallback → `("needs-review", "adapt", ...)`

**Project-specific marker extensibility:**
- `projectSpecificMarkers` becomes a list of `markerRule` structs instead of hardcoded `if` branches
- Each marker rule has a name and a match function
- New markers (e.g. "contains CI-specific secrets path", "mentions cloud project ID") are table entries

**Governance signal extensibility:**
- `sectionSignals` switch cases become a map of section name → signal definitions
- Each signal definition has a name and a match function operating on normalized text
- Adding signals for new governed sections is a map entry, not a new `case` branch

**What this does NOT change:**
- The enhance workflow, ranking, AC generation, or three-way comparison logic
- The `enhancementMapping` list (file-to-template mappings remain hardcoded — they're structural, not classification)
- External config files — rules remain in Go code for compile-time safety

### Phase 4 — Assisted Apply (implementation scope)

**Problem:** After enhance produces an AC doc, the template maintainer must manually locate the reference file, compare it to the template target, and make the edit by hand. For `accept`+`portable` candidates — especially whole-file replacements — this is tedious and error-prone. The review step is valuable, but the apply step has no tooling support.

**Solution:** Add an `--apply` flag to enhance mode that, after the normal review-and-AC-creation step, writes `.template-proposed` files for each actionable candidate. This reuses the existing proposal pattern from adopt mode (`proposeIfExists`). The maintainer reviews the proposals, accepts/rejects them manually, then deletes the proposal files. No automatic overwrite, no bypassing human review.

**Workflow with `--apply`:**
1. Normal enhance review runs (compare, classify, rank, create AC doc)
2. For the selected candidate (and optionally deferred actionable candidates):
   - Always write the reference content to `proposalPath(templateTarget)` (e.g. `README.template-proposed.md.tmpl`)
   - This applies whether or not the template target already exists — proposals are always side-by-side artifacts, never live files
3. Print a summary of proposed files
4. Maintainer reviews proposals, applies or discards them

**What `--apply` does NOT do:**
- It does not overwrite existing template files
- It does not skip the AC doc creation
- It does not auto-accept — proposals are side-by-side files requiring manual merge
- It does not modify the reference repo

**Collision detection:**
- Reuses `proposalPath()` from adopt mode to generate `.template-proposed` paths
- If a proposal file already exists from a prior run, overwrite it (proposals are ephemeral)
- `CollisionImpact` field on the candidate already indicates risk level

**Config and flag changes:**
- Add `Apply bool` to `Config` struct
- Add `--apply` / `-a` flags to the flag parser
- `validateConfig`: `--apply` is only valid with `--mode enhance`
- `--apply` and `--dry-run` can coexist (dry-run previews what would be proposed)

**Scope boundary:** This phase writes proposals for file-level candidates only. Governance section-level candidates (where the template target is `base/AGENTS.md` and only one section differs) get a proposal for the entire AGENTS.md file — section-level merge is out of scope and noted in the AC doc for manual attention.

## In Scope (Phase 1 — COMPLETE)

- Refine `governanceSectionCovered()` to compare at constraint level, not just keyword presence
- Add section-level diffing for structured markdown in `reviewMappedFile()`
- Update `EnhancementCandidate` to carry section-level delta info when available
- Update `printEnhancementSummary` to show section-level detail
- Update `renderACDoc` to include section-level detail in the generated AC
- Tests for the new comparison logic
- Tests for section-level file diffing

## In Scope (Phase 2)

- New file `internal/bootstrap/manifest.go` — `Manifest`/`ManifestEntry` structs, `computeChecksum`, `buildManifest`, `formatManifest`, `parseManifest`, `readManifest`
- Add `source` field to `operation` struct to track which template file produced each output
- Refactor `planRender` into `planCanonical` + `applyAdoptTransforms` to expose pre-transform operations for manifest building
- Modify `runNewOrAdopt` to build manifest from canonical ops and write `.repokit-manifest` into generated repos
- Add `ChangeOrigin` field to `EnhancementCandidate`
- Modify `ReviewEnhancement` to read manifest from reference repo
- Modify `reviewMappedFile` to perform three-way comparison when manifest is available
- Modify `reviewGovernedSections` to use whole-file AGENTS.md checksum as pre-filter
- Modify `classifyEnhancement` to defer template-only changes
- Update `printEnhancementSummary` and `renderACDoc` to show change origin
- Tests for all manifest functions and three-way comparison scenarios

## Out Of Scope

- Changes to the AC workflow itself (ranking, selection, rendering structure)
- Per-section checksums in manifest (whole-file checksums are sufficient for Phase 2)
- External config files for classification rules (rules stay in Go code for compile-time safety)
- Changes to `enhancementMapping` list (file-to-template mappings are structural, not classification)
- Section-level merge for governance candidates (Phase 4 proposes the whole file; manual section merge)
- Automatic creation or overwrite of live template files (Phase 4 always writes `.template-proposed` side-by-side proposals)

## In Scope (Phase 3)

- Define `classificationRule` struct with match function, portability, disposition, reason, and priority
- Define `classificationContext` struct carrying all inputs (`Content`, `ReferenceRoot`, `TemplateTarget`, `Governance`, `Area`)
- Replace `classifyEnhancement` body with a rule table scan (first matching rule by priority wins)
- Populate default rule table with the four existing classification branches (project-specific, governance, workflow helpers, default)
- Replace `projectSpecificMarkers` hardcoded `if` branches with a `markerRule` slice (name + match function)
- Replace `sectionSignals` switch/case with a `map[string][]signalDef` (section name → signal definitions)
- Each `signalDef` has a name and a match function operating on normalized text
- Tests verifying default rules reproduce identical behavior to the current hardcoded logic
- Tests verifying a custom rule inserted into the table takes precedence at the correct priority

## In Scope (Phase 4)

- Add `Apply bool` field to `Config` struct
- Add `--apply` / `-a` flags to the flag parser
- Update `validateConfig` to restrict `--apply` to enhance mode only
- Add `applyProposals` function that writes `.template-proposed` files for actionable candidates
- Modify `RunEnhance` to call `applyProposals` after AC doc creation when `--apply` is set
- For all candidates: always write reference content to `proposalPath(templateTarget)` — proposals are always side-by-side artifacts, never live files
- For governance candidates: propose the entire reference AGENTS.md file (not section-level merge)
- Print summary of proposed files after writing
- Dry-run support: `--apply --dry-run` previews proposals without writing
- Tests for all apply scenarios

## Implementation Notes (Phase 1)

- Governed-section refinement: normalizes each section into constraints (one per bullet), compares constraint sets after signal pre-filter
- File diffing: parses `##`-delimited sections, compares per-section, reports which sections differ; falls back to whole-file for unstructured files
- `extractConstraints` handles bullet lists (`- `/`* `) and numbered lists (`1. `/`2. `)
- `DeltaSections []string` field on `EnhancementCandidate` carries section-level info
- `diffMarkdownSections` walks both template and reference sections to catch additions and removals

## Implementation Notes (Phase 2)

- Manifest format is `repokit-manifest-v1`, line-oriented, sorted by path — no JSON/YAML, trivially parseable with `strings.Split`/`strings.Cut`
- Dual checksums: rendered content (`sha256`) for detecting user changes, source template (`source-sha256`) for detecting template changes
- `planCanonical` returns operations with canonical paths and content; adopt transforms applied separately
- Manifest records canonical operations regardless of adopt propose/skip — establishes correct baseline
- For AGENTS.md: whole-file checksum pre-filters before section-level comparison
- Missing or corrupt manifest → graceful fallback to current two-way comparison, no crash
- `.repokit-manifest` is NOT gitignored — useful metadata that should travel with the repo
- New imports: `crypto/sha256`, `encoding/hex` — all stdlib

## Implementation Notes (Phase 3)

- `classificationRule` table is a package-level `[]classificationRule` variable, sorted by priority
- `classifyEnhancement` iterates the table and returns the first match; the last entry is the catch-all default
- `projectSpecificMarkers` becomes `var defaultMarkerRules []markerRule` — each entry is `{Name string, Match func(content, refRoot string) bool}`
- `sectionSignals` becomes `var defaultSignalDefs map[string][]signalDef` — each `signalDef` is `{Name string, Match func(text string) bool}`
- The existing `containsAll`/`containsAny` helpers remain unchanged; signal match functions use them
- No external config parsing, no YAML/JSON — rules are Go structs initialized in `var` blocks
- Existing tests should pass without modification (behavior-preserving refactor); new tests verify extensibility

## Implementation Notes (Phase 4)

- `applyProposals` takes the selected candidate, deferred actionable candidates, the template root, and the reference root
- For each candidate, reads reference file content and writes it to `proposalPath(templateTarget)` — always a `.template-proposed` artifact, never a live file
- Governance candidates use the reference AGENTS.md content; the proposal target is `proposalPath("base/AGENTS.md")`
- File candidates use `candidate.Path` for content and `proposalPath(candidate.TemplateTarget)` for the proposal path
- `--apply` without `--dry-run` writes AC doc AND proposal files in the same pass
- `--apply` with `--dry-run` prints what would be proposed without writing anything
- Proposal files are ephemeral — they should be reviewed and deleted after merge, not committed
- No changes to `ReviewEnhancement`, classification, or ranking logic

## Acceptance Tests (Phase 1 — COMPLETE)

- [Automated] Two governance sections with same keywords but different constraints produce a candidate
- [Automated] Two governance sections with equivalent constraints produce no candidate
- [Automated] Overlay file with one changed section reports that specific section, not whole-file
- [Automated] Overlay file with no `##` structure falls back to whole-file comparison
- [Automated] Enhancement summary includes section-level detail when available
- [Automated] Generated AC doc includes section-level detail when available

## Acceptance Tests (Phase 2)

- [Automated] `formatManifest`/`parseManifest` round-trip produces identical manifest
- [Automated] `parseManifest` rejects unrecognized format version
- [Automated] `readManifest` returns false for missing file without error
- [Automated] `readManifest` returns false for corrupt file without error
- [Automated] Bootstrap `new` mode writes `.repokit-manifest` with correct entries
- [Automated] Bootstrap `adopt` mode writes manifest with canonical checksums
- [Automated] Enhance with manifest: user-only change evaluates normally (not deferred)
- [Automated] Enhance with manifest: template-only change deferred as stale
- [Automated] Enhance with manifest: both changed produces candidate with ChangeOrigin `"both"`
- [Automated] Enhance with manifest: neither changed produces no candidate
- [Automated] Enhance without manifest falls back to current two-way behavior
- [Automated] Enhancement summary includes change-origin when available
- [Automated] Generated AC doc includes change-origin when available
- [Manual] Bootstrap a repo, modify a file, update template source, run enhance, verify three-way classification

## Acceptance Tests (Phase 3)

- [Automated] Default classification rules produce identical output to pre-refactor hardcoded logic for all existing test scenarios
- [Automated] Project-specific marker rules detect repo name mention and absolute user paths
- [Automated] A custom classification rule inserted at higher priority overrides the default for matching content
- [Automated] A custom marker rule added to the marker list is evaluated alongside defaults
- [Automated] A custom signal definition added for an existing section is recognized by `sectionSignals`
- [Automated] Signal definitions for unknown section names return empty (no panic, no false match)
- [Automated] Rule table with no matching rules falls through to the default catch-all

## Acceptance Tests (Phase 4)

- [Automated] `--apply` with actionable candidate writes `.template-proposed` file at correct path
- [Automated] `--apply` still creates the AC doc (apply does not skip the AC)
- [Automated] `--apply` with candidate where template target does not exist still writes `.template-proposed` file (not the live target)
- [Automated] `--apply` with governance candidate proposes entire AGENTS.md file
- [Automated] `--apply --dry-run` prints proposal paths but writes no files
- [Automated] `--apply` without `--mode enhance` fails validation
- [Automated] Enhance without `--apply` produces no proposal files (existing behavior preserved)
- [Manual] Run enhance with `--apply` against a real governed repo, review proposals, verify they match the reference content

## Documentation Updates

- `docs/bootstrap-model.md` — document manifest and three-way comparison in Enhancement Review section
- `docs/bootstrap-model.md` — document `.repokit-manifest` as a template-owned artifact
- `docs/bootstrap-model.md` — update Enhancement Review section to note data-driven classification
- `docs/bootstrap-model.md` — document `--apply` flag behavior in enhance mode section

## Status

PHASE 1 COMPLETE — comparison refinement shipped
PHASE 2 COMPLETE — version awareness with bootstrap manifest and three-way comparison shipped
PHASE 3 COMPLETE — classifier extensibility shipped
PHASE 4 COMPLETE — assisted apply with `.template-proposed` files shipped
