# Mallard-OS — Development Progress

Running development log for Mallard-OS (TinyGo on Seeed XIAO ESP32-S3 N16R8).

The architecture, design decisions and long-term roadmap live in `MallardOS.md`.

---

## Current Milestone: v0.2 — Interactive Shell and Built-in Commands

### Development checklist

#### v0.1 — Boot and Interactive Serial Terminal

- [x] Design document drafted (TinyGo-first, C only for hardware gaps).
- [x] TinyGo toolchain installed and compiling.
- [x] Flashing ESP32-S3 N16R8 successfully.
- [x] Hello World / boot messages over USB serial.
- [x] HAL initialization.
- [x] CPU frequency configuration (160 MHz).
- [x] USB Serial/JTAG communication.
- [x] Hardware-independent console abstraction.
- [x] Terminal logging.
- [x] Character reception and echo.
- [x] Input line buffering.
- [x] Enter key handling.
- [x] Backspace handling.

**Status: COMPLETE**

#### v0.2 — Shell and Built-in Commands

- [x] Create `internal/shell/`.
- [x] Implement command registration.
- [x] Implement command parsing.
- [x] Implement command dispatch.
- [x] Implement command handler interface.
- [x] Integrate shell with console line editor.
- [x] Implement `help`.
- [x] Implement `info`.
- [x] Implement `uptime`.
- [x] Implement `meminfo`.
- [x] Register built-in commands.
- [x] Verify command execution on physical hardware.

**Status: COMPLETE**

---

## Current Architecture

```text
cmd/
└── mallard-os/
    └── main.go

internal/
├── boot/
│   ├── boot.go
│   └── boot_test.go
│
├── cmd/
│   ├── command.go
│   ├── help.go
│   ├── info.go
│   ├── meminfo.go
│   └── uptime.go
│
├── console/
│   ├── console.go
│   ├── line.go
│   └── log.go
│
├── drivers/
│   └── ws2812/
│       └── hardware-specific driver implementations
│
├── hal/
│   ├── gpio.go
│   ├── gpio_machine.go
│   ├── manager.go
│   ├── serial.go
│   └── serial_machine.go
│
├── logger/
│   └── logger.go
│
├── runtime/
│   └── runtime.go
│
├── shell/
│   ├── parser.go
│   ├── registry.go
│   └── shell.go
│
└── svc/

Progress.md
MallardOS.md
README.md
```

---

## Verified Hardware Configuration

| Component | Current state |
|---|---|
| MCU | ESP32-S3 |
| Board | Seeed XIAO ESP32-S3 N16R8 |
| Flash capacity | 16 MB board specification |
| PSRAM capacity | 8 MB board specification; allocation unverified |
| CPU frequency | 160 MHz |
| Language | TinyGo |
| TinyGo version | 0.44.0 |
| Go version | 1.26.8 |
| LLVM version | 22.1.4 |
| Serial transport | USB Serial/JTAG |
| USB device | `/dev/ttyACM0` |
| TinyGo target | `xiao-esp32s3` |
| Flash baud rate | 460800 |

---

## Shell Verification

Successfully tested on physical hardware.

```text
mallard>
mallard> help
Mallard-OS Commands:
help - Display available commands
info - Display Mallard-OS system information
meminfo - Display runtime memory information
uptime - Display system uptime

mallard> info
Mallard-OS v0.1
Embedded Operating Environment
Architecture: ESP32-S3
Runtime: TinyGo

mallard> uptime
...

mallard> meminfo
...

mallard> unknown
unknown command: unknown
```

*Output is representative; uptime and memory values depend on the running firmware.*

### Verified functionality

- Command registration.
- Command lookup.
- Command execution.
- Command descriptions.
- Shell command parsing and dispatch.
- Unknown-command handling.
- Built-in command handlers.
- Console output through `CommandIO`.
- Runtime uptime tracking.
- TinyGo runtime memory statistics.

### Current built-in commands

| Command | Description | Status |
|---|---|---|
| `help` | Display available commands | PASS |
| `info` | Display Mallard-OS system information | PASS |
| `uptime` | Display system uptime | PASS |
| `meminfo` | Display runtime heap statistics | PASS |

### Memory monitoring

Mallard-OS uses TinyGo's `runtime.ReadMemStats()` API.

Current metrics include:

- Heap allocated.
- Heap in use.
- Heap idle.
- Heap total.
- Heap objects.
- Total allocations.
- Garbage collection cycles.

Automatic garbage collection is provided by TinyGo.

**Important:** These values describe runtime heap statistics. They do not represent total physical SRAM, PSRAM or flash capacity.

---

## Experiments

Test order for investigating unknown hardware and runtime capabilities.

Update status only after testing on the physical board.

| # | Experiment | Status | Date | Notes |
|---|---|---|---|---|
| 1 | USB serial console | PASS | 2026-11-20 | USB Serial/JTAG, input, echo, line editing and Enter verified. |
| 2 | Shell and command execution | PASS | 2026-11-21 | Command registry, parser, dispatch and built-in commands verified. |
| 3 | Dual-core test: two goroutines pinned to different cores, GPIO toggle, shared counter with mutex | Untested | | Verify actual core execution, synchronization and concurrency. |
| 4 | LittleFS over flash `BlockDevice` | Untested | | Investigate persistent filesystem support. |
| 5 | PSRAM allocation | Untested | | Verify external RAM allocation on N16R8. |
| 6 | TLS client connection | Untested | | Verify networking and TLS compatibility. |
| 7 | `recover` restart under the `cores` scheduler | Untested | | Investigate panic recovery and scheduler behavior. |

### Experiment results

#### Experiment 1 — USB Serial Console

**Date:** 2026-11-20

**Result:** PASS

What was verified:

- Successful firmware compilation.
- Successful flashing and MD5 verification.
- USB Serial/JTAG initialization.
- Terminal input reception.
- Character echo.
- Line buffering.
- Enter and Backspace handling.
- Terminal logging.

Conclusion:

The USB serial transport is operational on the ESP32-S3. The console layer now provides the foundation for the Mallard-OS shell.

#### Experiment 2 — Shell and Command Execution

**Date:** 2026-11-21

**Result:** PASS

What was implemented:

- Shell package.
- Command parser.
- Command registry.
- Command lookup and registration.
- Command dispatch.
- Command handler interface.
- Built-in command registration.
- `help` command.
- `info` command.
- `uptime` command.
- `meminfo` command.

What was verified:

- Commands can be registered.
- Commands can be discovered through `help`.
- Commands execute through the registry.
- Command output is routed through `CommandIO`.
- Unknown commands return an appropriate error.
- Uptime is tracked from system initialization.
- TinyGo runtime memory statistics can be displayed.

Conclusion:

The interactive shell is operational. Mallard-OS now has a working command execution foundation that can be extended with filesystem, hardware and system management commands.

---

## Decisions

### Language architecture

- Go (TinyGo) is the primary implementation language.
- C is permitted only as a hardware capability gap-filler behind Go interfaces.
- Python is reserved for host-side development tools, testing and automation.

### Software architecture

- `boot`: System initialization and boot sequence.
- `hal`: Hardware abstraction and device interfaces.
- `drivers`: Hardware-specific implementations.
- `console`: Terminal input/output and line editing.
- `logger`: System-level logging.
- `runtime`: Runtime utilities, uptime and memory statistics.
- `shell`: Command parsing, registration and execution.
- `cmd`: Built-in command handlers and registration.
- `svc`: Planned service supervision.

### Hardware abstraction

Hardware access must remain behind HAL interfaces wherever practical.

The console layer must not depend directly on `machine.Serial`.

### Shell architecture

- Commands are registered through `shell.Registry`.
- Commands use a shared `CommandHandler` function signature.
- Command handlers receive `shell.CommandIO` and command arguments.
- The shell is independent of the underlying serial transport.
- Built-in commands live in `internal/cmd/`.

### Runtime philosophy

- Keep the system lightweight.
- Avoid unnecessary dependencies.
- Prefer native TinyGo capabilities.
- Validate uncertain runtime features experimentally.
- Do not assume standard Go runtime behavior is available until verified on the target.
- Monitor runtime memory usage without assuming heap statistics represent all physical memory.

---

## Blockers

None currently identified.

### Known limitations

- No persistent filesystem.
- Dual-core runtime behavior has not been experimentally verified.
- PSRAM allocation has not been verified.
- No networking or TLS validation yet.
- No task manager or service supervisor.
- No application management system.

---

## Next

### v0.3 — Persistent Filesystem

Before implementing filesystem commands, investigate LittleFS support on the XIAO ESP32-S3 N16R8.

Planned development sequence:

1. Investigate TinyGo flash and `BlockDevice` support.
2. Verify LittleFS compatibility.
3. Design a storage HAL interface.
4. Implement the storage backend.
5. Implement filesystem initialization and mounting.
6. Add filesystem commands:
   - `ls`
   - `mkdir`
   - `cat`
   - `write`
7. Verify persistence across device reboots.

### Next hardware experiment — Dual-Core Runtime

Before committing to multicore-dependent services, verify the selected `cores` scheduler.

Test requirements:

- Two goroutines.
- Different CPU cores.
- GPIO activity.
- Shared counter.
- Mutex synchronization.
- Verify actual multicore execution.
- Check for deadlocks and watchdog resets.
- Keep USB serial responsive.

### Long-term roadmap

- Persistent filesystem (LittleFS).
- PSRAM memory management.
- Task and service supervision.
- Networking and TLS.
- Recovery and restart behavior.
- Hardware monitoring and management.
- Application management.

---

## Development Log

### 2026-11-20

**What I did:**

- Implemented hardware-independent console abstraction.
- Added terminal logging.
- Implemented terminal line editor.
- Integrated console and line editor into `main.go`.
- Built and flashed Mallard-OS onto the ESP32-S3.

**What worked:**

- TinyGo compilation.
- ESP32-S3 firmware flashing.
- Flash verification.
- HAL initialization.
- CPU frequency configuration at 160 MHz.
- USB Serial/JTAG communication.
- Character echo.
- Line buffering and Enter detection.
- Backspace editing.
- Terminal logging.

**What broke / open questions:**

- No known blocking issues.
- Command parsing and dispatch remained to be implemented at this stage.
- Dual-core runtime behavior remained untested.

**Milestone:**

USB serial terminal operational.

### 2026-11-21

**What I did:**

- Implemented shell command registration and dispatch.
- Implemented command parsing.
- Added the command handler interface.
- Integrated the shell with the console.
- Implemented built-in commands:
  - `help`
  - `info`
  - `uptime`
  - `meminfo`
- Added TinyGo runtime memory statistics.
- Tested command execution on physical hardware.

**What worked:**

- Shell command loop.
- Command parsing and dispatch.
- Command registration.
- Help output.
- System information output.
- Uptime reporting.
- Runtime memory statistics.
- Unknown-command handling.
- TinyGo automatic garbage collection remains managed by the runtime.

**What broke / open questions:**

- No known blocking issues.
- Dual-core runtime behavior remains untested.
- Persistent storage remains unimplemented.
- PSRAM allocation remains unverified.

**Current milestone:**

v0.1 and v0.2 complete. Interactive shell operational with four built-in commands.

---

**Mallard-OS — An OS with a duck going one step at a time.**
