# Benchmark & Performance Notes

Current state: no formal benchmark suite yet. Guidance for future work:

- Measure cold-start time for `buddy run mock-echo` and `buddy run claude-dm` on macOS/Linux (arm64/amd64).
- Track latency per request for transports/agents (nostr relay RTT dominates).
- Stress shell action with max_output limits to confirm truncation behavior.
- Profile CPU/memory with `pprof` flags on local runs; identify hot paths in transports/agents.
- Include results in release QA when added; keep assets small and binaries CGO=0 for portability.
