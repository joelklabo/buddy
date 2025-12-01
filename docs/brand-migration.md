# Brand Migration Plan – Issue 3oa.41

Scope: update visible text/assets from "nostr-codex-runner" to "buddy" once repo/module rename is executed.

Checklist
- README: project name, badges (CI/release/coverage/staticcheck), quickstart commands (`buddy`, `buddy wizard`, `buddy run <preset>`), release links, install URLs.
- Docs index & pages: replace legacy name; ensure examples use buddy binary and new config path (`~/.config/buddy/config.yaml`).
- Wizard copy: update references to buddy; include alias note.
- Preset/CLI docs: update command samples.
- Systemd/launchd recipes: service names and ExecStart paths to `buddy`/`nostr-buddy`.
- Scripts: `scripts/install.sh`, `scripts/run.sh` text.
- Assets: social preview image if it contains text; any SVG/PNG badges.
- Sample configs: paths/comments referencing old name.
- References to `nostr-codex-runner` in code comments/logs (search and replace carefully after module rename).
- FAQ/Backcompat: keep a section about legacy binary/env for one release.

Ordering
1) After repo + module rename (issues 3oa.38/3oa.39).
2) Update goreleaser archive names and install URLs (issue 3oa.28/3oa.40 handles the pipeline; this doc focuses on wording).
3) Run `rg "nostr-codex-runner"` to catch stragglers; update docs accordingly.

Success criteria
- README first screen shows buddy branding and correct commands.
- All docs examples use buddy names; only backcompat section mentions old name.
