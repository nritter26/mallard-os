# Progress

Running log for MallardOS (TinyGo on ESP32-S3 N16R8). Roadmap and design live in `MallardOS.md`.

## Current: v0.1 (boot + serial shell)

- [x] Design doc drafted (TinyGo-first, C only for gaps)
- [ ] Toolchain set up (TinyGo 0.44+, flash works on the N16R8)
- [ ] Hello world over USB serial
- [ ] Shell loop (read line, parse, dispatch)
- [ ] `help`
- [ ] `sysinfo`
- [ ] `uptime`
- [ ] `meminfo`

## Experiments (test order for the unknowns)

Update the status as each one is tested on the board.

| # | Experiment | Status | Date | Notes |
|---|---|---|---|---|
| 1 | USB serial console | Untested | | |
| 2 | Dual-core test: two goroutines pinned to different cores, GPIO toggle, shared counter with mutex | Untested | | |
| 3 | LittleFS over flash `BlockDevice` | Untested | | |
| 4 | PSRAM allocation | Untested | | |
| 5 | TLS client connection | Untested | | |
| 6 | `recover` restart under the `cores` scheduler | Untested | | |

## Decisions

- Language split: Go (TinyGo) primary, C only behind HAL interfaces, Python for host-side tooling.

## Blockers

- None yet.

## Next

- Run experiment 1 (USB serial console), since the whole v0.1 shell depends on it.
- Then experiment 2 (dual-core test).

## Log

### Template

```
### YYYY-MM-DD
- What I did:
- What worked:
- What broke / open questions:
```
