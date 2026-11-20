# MallardOS — Development Progress

Running development log for MallardOS (TinyGo on Seeed XIAO ESP32-S3 N16R8).

The architecture, design decisions, and long-term roadmap live in `MallardOS.md`.

---

## Current: v0.1 — Boot + Interactive Serial Terminal

### Development checklist

- [x] Design document drafted (TinyGo-first, C only for hardware gaps)
- [x] TinyGo toolchain installed and compiling
- [x] Flashing ESP32-S3 N16R8 successfully
- [x] Hello World / boot messages over USB serial
- [x] HAL initialization
- [x] CPU frequency configuration (160 MHz)
- [x] USB Serial/JTAG communication
- [x] Hardware-independent console abstraction
- [x] Terminal logging
- [x] Character reception and echo
- [x] Input line buffering
- [x] Enter key handling
- [x] Backspace handling
- [ ] Shell loop (read line, parse, dispatch)
- [ ] `help`
- [ ] `sysinfo`
- [ ] `uptime`
- [ ] `meminfo`

### Current architecture

```text
cmd/
└── mallard-os/
    └── main.go

internal/
├── boot/
│   ├── boot.go
│   └── boot_test.go
│
├── console/
│   ├── console.go
│   ├── log.go
│   └── line.go
│
├── drivers/
│
├── hal/
│   ├── gpio.go
│   ├── gpio_machine.go
│   ├── manager.go
│   ├── serial.go
│   └── serial_machine.go
│
└── logger/
    └── logger.go
```

### Verified hardware configuration

| Component | Current state |
|---|---|
| MCU | ESP32-S3 |
| Board | Seeed XIAO ESP32-S3 N16R8 |
| Flash | 16 MB detected |
| CPU frequency | 160 MHz |
| Language | TinyGo |
| Serial transport | USB Serial/JTAG |
| USB device | `/dev/ttyACM0` |
| Flash target | `esp32s3-generic` |
| Flash baud rate | 460800 |

### Console verification

Successfully tested on physical hardware:

```text
[INFO] Starting MallardOS...
[INFO] Initializing: Hardware Abstraction Layer
[INFO] Hardware Abstraction Layer initialized
[INFO] Boot sequence completed.
[INFO] Configuring platform clock
[INFO] CPU frequency: 160000000
[INFO] MallardOS hardware initialized

[INFO] MallardOS starting...
[ OK ] USB initialized
[WARN] Running without persistent storage
It's an OS with a duck...
Waiting for input...

> hello mallard
[INFO] Received: hello mallard
>
```

Verified functionality:

- USB serial communication.
- Character reception.
- Character echo.
- Line completion.
- Backspace editing.
- Terminal output.
- Continuous main loop operation.

---

## Experiments

Test order for investigating unknown hardware and runtime capabilities.

Update status only after testing on the physical board.

| # | Experiment | Status | Date | Notes |
|---|---|---|---|---|
| 1 | USB serial console | PASS | 2026-11-20 | USB Serial/JTAG, input, echo, line editing and Enter verified. |
| 2 | Dual-core test: two goroutines pinned to different cores, GPIO toggle, shared counter with mutex | Untested | | Verify actual core execution, synchronization and concurrency. |
| 3 | LittleFS over flash `BlockDevice` | Untested | | Investigate persistent filesystem support. |
| 4 | PSRAM allocation | Untested | | Verify external RAM allocation on N16R8. |
| 5 | TLS client connection | Untested | | Verify networking and TLS compatibility. |
| 6 | `recover` restart under the `cores` scheduler | Untested | | Investigate panic recovery and scheduler behavior. |

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

The USB serial transport is operational on the ESP32-S3. The console layer can now serve as the foundation for the MallardOS shell.

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
- `shell`: Command parsing and execution (planned).

### Hardware abstraction

Hardware access must remain behind HAL interfaces wherever practical.

The console layer must not depend directly on `machine.Serial`.

### Runtime philosophy

- Keep the system lightweight.
- Avoid unnecessary dependencies.
- Prefer native TinyGo capabilities.
- Validate uncertain runtime features experimentally.
- Do not assume standard Go runtime behavior is available until verified on the target.

---

## Blockers

None currently identified.

### Known limitations

- No command interpreter yet.
- No persistent filesystem.
- Dual-core behavior has not been experimentally verified.
- PSRAM allocation has not been verified.
- No networking or TLS validation yet.

---

## Next

### v0.1 — Complete the interactive shell

1. Create `internal/shell/`.
2. Implement command registration.
3. Implement command parsing.
4. Implement command dispatch.
5. Add built-in commands:
   - `help`
   - `sysinfo`
   - `uptime`
   - `meminfo`
6. Integrate the shell with the existing console line editor.
7. Verify command execution on physical hardware.

### v0.2 — Runtime and multicore experiments

Proceed with Experiment 2:

- Two goroutines.
- Different CPU cores.
- GPIO activity.
- Shared counter.
- Mutex synchronization.
- Verify actual multicore execution.

### Long-term roadmap

- Persistent filesystem (LittleFS).
- PSRAM memory management.
- Task and service supervision.
- Networking and TLS.
- Recovery and restart behavior.
- Hardware monitoring and management.

---

## Development Log

### 2026-11-20

**What I did:**
- Implemented hardware-independent console abstraction.
- Added terminal logging.
- Implemented terminal line editor.
- Integrated console and line editor into `main.go`.
- Built and flashed MallardOS onto the ESP32-S3.

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
- Command parsing and dispatch remain unimplemented.
- Dual-core runtime behavior remains untested.

**Current milestone:**
USB serial terminal operational. Preparing the first MallardOS command interpreter.
