# Buddy Docs, Wizard, and CLI Productization Epic

Goal: ship a binary-first "buddy" CLI with a fast install, guided wizard, presets, and clear docs for users and contributors.

## Issue list (Todo unless noted)

### Docs & Onboarding

1. P0 README information architecture (binary-first; quickstart with `buddy wizard` and presets). ✅
2. P0 Quick install path (curl/Homebrew; 2-minute smoke using a preset; no git clone).
3. P0 Example use cases page (argument-only UX; presets + wizard outputs). ✅ draft
4. P0 Docs landing/index (User vs Contributor entry; links to wizard/presets). ✅ draft
5. P1 Refresh README/landing screenshots to match wizard-first flow (pending).
6. P1 Publish offline local-llm mock recipe link in top-level docs/README (✅).
7. P1 Architecture diagram (include CLI front door and preset loader).
8. P1 Extend/commands guide (register new CLI subcommands/presets).
9. P1 Config reference (search order, preset schema, defaults).
10. P1 Contributing refresh (CLI-first stance; existing standards).
11. P1 Local dev setup (contributors only).
12. P1 Docs style guide/templates.
13. P2 Changelog policy.
14. P2 FAQ/Troubleshooting (install/preset/wizard).
15. P2 Benchmark notes.
16. P2 Security/secrets (wizard + CLI logging).
17. P3 Website/landing polish.
18. P3 Precedent research (great onboarding examples).

### Wizard Track

1. P0 Wizard concept (goals, scope, success metrics). ✅
2. P0 Wizard IA/script (questions → config + preset choice). ✅
3. P1 Wizard UX prototype (survey/promptui/bubbletea vs bufio).
4. P1 Wizard implementation (`buddy wizard`; writes config; optional smoke run). ✅
5. P1 Wizard tests (branching, goldens, stdin sim). ✅ initial
6. P1 Wizard docs (README blurb + page; clip). ✅ initial
7. P2 Wizard telemetry/safety (no secret logging; dry-run).
8. P2 Wizard extensibility (registry for new actions/providers).
9. P3 Wizard polish (presets, color toggle, retries).

### CLI Productization

1. P0 CLI spec/map (`buddy wizard`, `buddy run <preset|config.yaml>`, `buddy presets`, `buddy help`; arguments over flags). ✅ v1
2. P0 Preset library (ship built-ins: `copilot-shell`, `claude-dm`, `local-llm`; assets/presets). ✅
3. P0 Packaging & releases (goreleaser, checksums, Homebrew tap).
4. P1 Install script (curl | sh; checksum; /usr/local/bin or ~/.local/bin).
5. P1 CLI UX copy/errors (friendly, masked secrets, exit codes).
6. P1 Config search precedence (arg path > cwd config.yaml > ~/.config/buddy/config.yaml; env opt-in).
7. P1 Release QA matrix (macOS/Linux/arm64; presets + wizard).
8. P2 Offline bundle (embed default presets/assets; graceful no-network). ✅ presets embedded; mock-echo default in wizard.
9. P2 Help/man page generation. 🚧
10. P3 Windows support decision (scoop/winget or "not yet").
11. P1 Metrics endpoint curl test in CI (✅).

### Name & Migration

1. P0 Name locked: "buddy".
2. P0 Repo rename plan (`nostr-codex-runner` → `buddy` repo; redirects, CI updates). ✅ done
3. P0 Module/binary rename (Go module path, imports, main package; build `buddy`). ✅ done
4. P1 Package manager updates (Homebrew formula rename, release artifacts).
5. P1 Brand migration in docs (README, docs index, wizard copy, badges). ✅ mostly done
6. P1 CLI help/man text with new name. ✅
7. P2 Domain/SEO check (distinct from Buddy.Works; optional microsite).

## Notes

- One issue per commit; run tests after each commit.
- Default collision mitigation: ship `buddy` plus optional alias (`bud` off by default).
- Rename done; legacy references now only for backcompat notes.
