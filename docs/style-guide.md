# Docs Style Guide & Templates – Issue 3oa.10

Tone & voice
- Direct, concise, imperative. Prefer verbs over adjectives.
- Keep first screen of each page scannable: short paragraphs, bullets, and code blocks.
- Avoid hype; state capability and how to do it.

Formatting
- Headings: sentence case except product names. Avoid deep nesting (>3 levels).
- Bullets: `-` only; keep to one line when possible.
- Commands/paths/env vars in backticks; multi-line commands in fenced `bash` blocks.
- Tables for config fields and presets; keep under 80–90 chars per cell.
- No trailing prose after code fences; end with punctuation.

Content patterns
- Start guides with: What it is → When to use → Steps → Verify → Next steps.
- Put prerequisites up front (OS/arch, deps, keys).
- Include a “Verify” step with an expected output snippet.
- Link to reference docs instead of repeating schema.

Templates
- **How-to**
  1) Goal sentence.
  2) Prereqs.
  3) Steps (numbered). Each step has a command and expected result.
  4) Troubleshooting tip if common failure.

- **Concept page**
  1) Definition (short).
  2) Why it matters / when to choose.
  3) Components diagram/link.
  4) Related tasks/recipes.

- **Reference (config/preset)**
  1) Scope sentence.
  2) Table of fields with defaults and examples.
  3) Notes on search order/precedence.

Accessibility/UX
- Avoid color-dependent meaning; keep text cues.
- Prefer plain ASCII punctuation; avoid smart quotes.
- Spell out abbreviations on first use.

Linking & cross-ref
- Relative links within repo; avoid raw URLs in text (use markdown link text).
- Keep README lean; link to docs index for depth.

Badges and names
- Use `buddy` as canonical name; mention alias `nostr-buddy` only where collisions matter.

Review checklist (for PRs touching docs)
- [ ] Quickstart still works as written.
- [ ] Commands copy-paste friendly.
- [ ] Secrets not logged or echoed in examples.
- [ ] Links resolve after repo rename.
