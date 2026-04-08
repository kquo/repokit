# AC-012 DOC Overlay Enrichment

## Objective Fit

1. Generated DOC repos get platform-specific publishing guidance and editorial variant options out of the box
2. Closes the last two priorities in plan.md
3. Brings the DOC overlay closer to parity with the enriched CODE overlay
4. Direct roadmap work

## Summary

Enrich the DOC overlay in two areas: (1) add platform-specific publishing examples to `publishing-workflow.md.tmpl` so generated repos have concrete guidance beyond the generic five-step workflow, and (2) add `voice.md.tmpl` and `calendar.md.tmpl` as optional alternate templates alongside the existing `style.md.tmpl` and `content-plan.md.tmpl`. Also add `docs/agent-roles/` to the DOC overlay (DEV, QA, maintainer) to match the CODE overlay pattern established in AC-008/AC-009.

## In Scope

- Enrich `overlays/doc/files/publishing-workflow.md.tmpl` with a "Platform-Specific Notes" section containing starter guidance for common platforms (Hugo, Jekyll, Substack, Ghost, WordPress, Notion)
- Create `overlays/doc/files/voice.md.tmpl` as an alternate to `style.md.tmpl` — focused on voice, persona, and audience rather than formatting rules
- Create `overlays/doc/files/calendar.md.tmpl` as an alternate to `content-plan.md.tmpl` — date-driven editorial calendar rather than priority-ordered backlog
- Create DOC-specific `docs/agent-roles/` templates: `README.md.tmpl`, `dev.md.tmpl`, `qa.md.tmpl`, `maintainer.md.tmpl`. These use DOC-appropriate language (editorial review, publishing workflow adherence, source/fact checking) instead of CODE build/release concepts.
- Update `overlays/doc/files/README.md.tmpl` to mention both variant pairs (style/voice, content-plan/calendar)
- Regenerate rendered examples under `examples/doc/`
- Update `overlays/doc/README.md` to list all new template files

## Out Of Scope

- Changing the bootstrap wizard to let users choose between style/voice or content-plan/calendar at generation time (both variants are generated; users delete whichever they don't need)
- Changes to CODE overlay
- Changes to bootstrap Go code

## Implementation Notes

- Platform-specific notes should be brief (2–3 lines per platform) and clearly labeled as starter guidance to customize. Use the `{{PUBLISHING_PLATFORM}}` placeholder where relevant.
- `voice.md.tmpl` covers: voice and persona definition, audience description, tone guidelines, and examples of on-brand vs off-brand language. Uses `{{DOC_STYLE}}` placeholder.
- `calendar.md.tmpl` covers: publishing cadence, upcoming scheduled content with target dates, and a simple date/title/status table. Uses `{{PROJECT_PURPOSE}}` and `{{PUBLISHING_PLATFORM}}` placeholders.
- Both variant files are generated alongside the originals. Update `publishing-workflow.md.tmpl` to note the variants: "This repo includes `style.md` and `voice.md` — keep whichever fits your editorial model and delete the other. Similarly for `content-plan.md` and `calendar.md`." This makes the swap path explicit rather than implied.
- DOC agent-roles use DOC-specific language:
  - `dev.md`: "Follow the publishing workflow", "Verify content against style.md or voice.md", "Never publish without explicit approval"
  - `qa.md`: "Verify content accuracy and source claims", "Check consistency against style guide", "Run editorial checks from publishing-workflow.md"
  - `maintainer.md`: combines both — implements content, performs self-review against editorial standards before presenting as complete
  - No references to build commands, test coverage, or pre-release checklists (those are CODE concepts)

## Acceptance Tests

- [Automated] Bootstrap `new` mode for DOC produces `publishing-workflow.md` containing "Platform-Specific Notes"
- [Automated] Bootstrap `new` mode for DOC produces `voice.md` and `calendar.md`
- [Automated] Bootstrap `new` mode for DOC produces `docs/agent-roles/dev.md`, `qa.md`, and `maintainer.md`
- [Automated] Bootstrap `adopt` mode for DOC proposes `voice.md`, `calendar.md`, `docs/agent-roles/dev.md`, and updated `publishing-workflow.md` when they already exist
- [Manual] Platform-specific notes contain guidance for at least three platforms
- [Manual] `voice.md` and `calendar.md` are usable standalone alternatives
- [Manual] `publishing-workflow.md` explains the variant swap path (which to keep, which to delete)
- [Manual] DOC agent-role docs use editorial language, not build/release language
- [Manual] Generated DOC README.md mentions both variant pairs
- [Manual] Overlay README lists all new files

## Documentation Updates

- `overlays/doc/README.md` — list all new template files
- `docs/bootstrap-model.md` — update ownership model to remove "(planned)" from voice.md and calendar.md

## Status

COMPLETE
