Adapted from `pbakaus/impeccable` 4.5.2, `reference/polish.md` (Apache-2.0), for Devin project use.

Polish is refinement, not a hidden redesign. Preserve the incumbent visual direction, content, behavior, and everything outside the user's scope. If the concept itself is wrong, explain that and ask before replacing the visual world.

## Establish the system

Read `DESIGN.md`, representative components, shared tokens, and neighboring flows. Classify drift as a missing token, a one-off implementation, a conceptual mismatch, or a local defect; fix it at the narrowest correct level. Ask when a binding design principle cannot be inferred.

## Gather evidence

Use the screen at desktop and mobile sizes. Exercise the real interaction path and check functional completion, user context, unfinished states, content lengths, and input methods. A detector result is evidence of a defect, not proof of quality.

## Triage

Fix in this order:
1. Broken or blocked tasks, data loss, misleading state, and inaccessible paths.
2. Missing loading, empty, error, success, disabled, and permission states.
3. Flow, hierarchy, responsive behavior, and design-system drift.
4. Visual and motion inconsistencies.
5. Dead code, redundant styles, and asset cleanup.

Do not perfect one corner while the rest of the flow remains below the same quality bar.

## Polish the whole path

- Keep the primary task and current state obvious; keep terminology and content factual.
- Align layout and spacing to the project system. Test long text, localization, zoom, and mobile widths.
- Use semantic color consistently and verify contrast for text, controls, and focus states.
- Use coherent icons and real images with stable aspect ratios and useful alt text.
- Ensure every control has appropriate hover, focus, active, disabled, loading, error, and success behavior.
- Keep keyboard order logical, focus visible, touch areas practical, and motion interruptible and performant.
- Validate empty, offline, slow, missing, and long-content states the product can encounter.
- Remove debug output, unused imports, dead code, obsolete styles, and polish-created duplication.

## Verify and finish

Walk the complete path with mouse, keyboard, and touch where practical. Check desktop, intermediate, and mobile layouts; empty/loading/error/success/disabled states; zoom, contrast, focus, semantics, console errors, and agreement with `DESIGN.md`. Run the detector and project build. Make one bounded correction pass and one final confirmation pass; stop after that rather than polishing indefinitely.
