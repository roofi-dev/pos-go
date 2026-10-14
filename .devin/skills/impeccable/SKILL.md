---
name: impeccable
description: Design and refine frontend interfaces with strong visual craft, clear UX, accessibility, responsive behavior, and intentional motion.
argument-hint: "[polish] [target]"
triggers:
  - user
  - model
---

# Impeccable for Devin

This project-local Devin adaptation is based on **Impeccable 4.5.2** by Paul Bakaus, licensed under Apache-2.0. Upstream source: https://github.com/pbakaus/impeccable. The upstream installer currently does not recognize Devin, so this skill uses Devin's supported `.devin/skills/<name>/SKILL.md` format rather than writing to another agent's configuration directory.

Use this skill for UI design, redesign, refinement, motion, accessibility, responsive layout, typography, and visual-quality work. POS screens are **Operate** surfaces: scanning speed, legibility, reliable controls, and touch/keyboard use take priority over decoration.

## Context and workflow

1. Read `PRD.md`, `DESIGN.md`, and the target component/styles. Preserve the existing product facts, interaction behavior, and established visual direction unless the user explicitly requests a redesign.
2. For a request to refine an existing screen, read [`reference/polish.md`](reference/polish.md). Do not turn refinement into a concealed redesign.
3. Once the direction is settled, read [`reference/craft-floor.md`](reference/craft-floor.md) immediately before editing UI.
4. Inspect the running screen at wide and mobile widths; exercise the main interaction, loading/error/empty/success states, keyboard focus, and touch target sizes. If browser capture is unavailable, state that limitation instead of claiming visual verification.
5. Make one bounded fix pass, then one confirmation pass. Use the official Impeccable detector when available: `npx --yes impeccable@4.1.0 detect <file-or-url>`. Pin the version; do not use an unbounded `latest` install.
6. Run the frontend build and report any remaining intentional limitations. Never imply that an unimplemented checkout or placeholder menu data is production-ready.

The upstream native `impeccable context` launcher is not bundled in this Devin adaptation. Load this project's context directly from `PRD.md`, `DESIGN.md`, and the relevant code; do not claim to have run the unavailable launcher.

## Design principles

- Let real use conditions and the written brief determine the visual language; do not default to generic templates.
- Establish hierarchy through layout, type, contrast, spacing, and restrained depth.
- Give every control clear default, hover, focus, active, disabled, loading, error, and success feedback as applicable.
- Use purposeful, interruptible motion and respect `prefers-reduced-motion`.
- Prefer real, licensed product photography for food menu imagery; do not invent image URLs or present emoji as final food photography.
- Keep text contrast, responsive behavior, keyboard access, and readable touch targets part of the design—not later polish.

## Supported workflow in this project

Invoke `/impeccable polish` for refinement of the current POS screen. The project includes the upstream polish and craft-floor playbooks in `reference/`. Other commands from the upstream Impeccable catalog can be added when their playbooks are vendored and needed.
