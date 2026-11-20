# MallardOS v0.1

*TinyGo-first embedded operating environment for the Seeed XIAO ESP32-S3 (N16R8). Status: Work in Progress.*

---

## 1. Goals

MallardOS aims to become a lightweight embedded operating environment with:

- Interactive terminal over USB Serial/JTAG.
- Command interpreter (shell) with a command registry.
- Persistent filesystem.
- Task (goroutine) management and a service supervisor.
- Hardware abstraction layer (HAL).
- GPIO, PWM, ADC and I²C management.
- Basic networking (Wi-Fi, HTTP).
- Eventually, an application management system.

The objective is to build an operating environment using Go as the primary development language while retaining access to lower-level hardware capabilities when necessary.

---

## 2. Language Strategy

| Language | Role |
|---|---|
| Go (TinyGo) | Primary: shell, services, HAL, drivers, networking and supervisor |
| C | Gap-filler only, behind HAL interfaces (see section 5) |
| Python | Host side only: serial client, test harness, build and flash scripts |

Everything that touches C should sit behind a small Go interface (`Storage`, `Crypto`, `Power`).

A C implementation should be replaceable with a Go implementation later without requiring changes to the shell or services.

---

## 3. Stack

| Technology | Responsibility | Status |
|---|---|---|
| TinyGo (0.44 or later) | Compiler and embedded runtime | Working |
| `cores` scheduler | Goroutines on both LX7 cores, with core affinity | Untested |
| espradio | Wi-Fi and Bluetooth, radio on core 0 | Planned |
| lneto | Network stack used by espradio | Planned |
| tinyfs (LittleFS) | Filesystem, C library wrapped for Go | Untested |
| `machine` package + TinyGo drivers | GPIO, PWM, ADC, I²C and SPI | Partially implemented |
| USB Serial/JTAG | Terminal communication | Verified |
| espflasher | Firmware flashing | Not currently used; flashing scripts working |

---

## 4. Architecture

MallardOS uses a layered architecture to separate hardware access, system services and user-facing functionality.

```text
                    MallardOS
                        |
        +---------------+---------------+
        |               |               |
      Shell           Services      Applications
        |               |               |
   Command API       System API    Application API
        |               |               |
        +---------------+---------------+
                        |
             Hardware Abstraction (Go)
                        |
                TinyGo Runtime
                        |
                    ESP32-S3
```

### 4.1 Current software architecture

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

### 4.2 Package responsibilities

| Package | Responsibility |
|---|---|
| `boot` | System initialization and boot sequence |
| `hal` | Hardware interfaces and device management |
| `drivers` | Hardware-specific driver implementations |
| `console` | Terminal transport abstraction, output and input handling |
| `logger` | System-level logging |
| `shell` | Command parsing, registration and execution (planned) |

### 4.3 Console architecture

The console is deliberately separated from the hardware abstraction layer.

```text
             Terminal / Shell
                    |
                    v
             Console Package
                    |
             Transport Interface
                    |
                    v
                HAL Console
                    |
                    v
              USB Serial/JTAG
```

Current console components:

- `console.go`: Hardware-independent terminal input/output.
- `log.go`: Human-readable terminal messages.
- `line.go`: Character buffering, echo, Enter and Backspace handling.

The HAL remains responsible for communicating with the underlying USB serial hardware.

This separation allows future console transports to be introduced without rewriting the terminal or shell.

### 4.4 Core split (initial plan)

The initial multicore design proposes:

- Core 0: Radio operations.
- Core 1: Shell and system services.

Goroutine core affinity will be handled through the `cores` scheduler.

This architecture remains experimental until the dual-core hardware test is completed.

---

## 5. C Gaps and Risks

C is reserved for capabilities that cannot reasonably be implemented or accessed through TinyGo.

Everything exposed through C should be wrapped behind Go interfaces.

| Area | Why C may be needed | Status |
|---|---|---|
| LittleFS | tinyfs wraps the C LittleFS library; requires an S3 flash `BlockDevice` | Untested |
| PSRAM (8 MB) | TinyGo initialization and allocation support need verification | Untested |
| USB Serial/JTAG console | Hardware communication successfully verified; native implementation details and alternative USB support remain to be investigated | Verified for current transport |
| TLS / HTTPS | Pure-Go TLS compatibility with lneto needs investigation; mbedTLS is a fallback | Untested |
| Hardware crypto | Go crypto may use software unless hardware acceleration is available | Untested |
| OTA and partition table | ESP-IDF handles these capabilities in C; TinyGo behavior needs verification | Untested |
| Sleep / low power / ULP | May require ESP-IDF C APIs | Untested |
| Other peripherals (I2S, RMT, LCD, camera) | Support may require C implementations | Untested |

### Experiment order

1. USB Serial/JTAG console — completed.
2. Dual-core runtime and synchronization — next.
3. LittleFS over flash `BlockDevice`.
4. PSRAM allocation.
5. TLS client connection.
6. Panic recovery under the `cores` scheduler.

---

## 6. Known Design Constraints

### Scheduling

The scheduler must be evaluated on the actual ESP32-S3.

Services must yield, block on channels or otherwise avoid starving other work.

Tight loops must not monopolize execution.

### TinyGo runtime

TinyGo supports only part of the standard Go library.

Avoid reflection-heavy packages and unnecessary runtime dependencies.

Binary size and memory consumption must be measured on the target hardware.

### Multicore synchronization

Shared state across cores requires appropriate synchronization using mutexes, channels or other verified mechanisms.

Do not assume goroutine core affinity works until experimentally confirmed.

### Panic recovery

Recoverable panics are intended to form the basis of supervised service restarts.

The behavior of `recover` must be verified under the selected multicore scheduler.

### Memory management

The N16R8 board provides external PSRAM, but its availability to TinyGo must be verified.

Do not assume all memory regions are automatically available through normal Go allocation.

---

## 7. Non-Goals (v1.0)

- No MMU.
- No process isolation.
- No dynamic loading.
- No Linux compatibility.
- No on-device second runtime such as MicroPython unless deliberately introduced in v0.7.

MallardOS is an embedded operating environment, not a Linux replacement.

---

## 8. Roadmap

| Version | Deliverable | Status |
|---|---|---|
| v0.1 | Boot sequence + interactive USB terminal | In progress |
| v0.2 | Command registry + built-in commands | Next |
| v0.3 | Filesystem + configuration storage (`ls`, `mkdir`) | Planned |
| v0.4 | Task manager + system monitoring (`ps`, with core shown) | Planned |
| v0.5 | Hardware management services (`gpio list`) | Planned |
| v0.6 | Wi-Fi + HTTP server (`wifi scan`) | Planned |
| v0.7 | Application management | Planned |
| v1.0 | First usable MallardOS release | Planned |

### v0.1 completion criteria

- [x] Boot sequence.
- [x] HAL initialization.
- [x] USB Serial/JTAG communication.
- [x] Terminal input/output.
- [x] Line editing.
- [ ] Shell command loop.
- [ ] Command parsing and dispatch.
- [ ] `help`.
- [ ] `sysinfo`.
- [ ] `uptime`.
- [ ] `meminfo`.

---

## 9. Target Commands by Version

| Command | Version | Purpose |
|---|---|---|
| `help` | v0.1 | Display available commands |
| `sysinfo` | v0.1 | Display system and hardware information |
| `uptime` | v0.1 | Display system uptime |
| `meminfo` | v0.1 | Display memory statistics |
| `ls` | v0.3 | List filesystem contents |
| `mkdir projects` | v0.3 | Create directories |
| `ps` | v0.4 | Display running tasks and core information |
| `gpio list` | v0.5 | List available GPIO resources |
| `wifi scan` | v0.6 | Scan available Wi-Fi networks |

---

## 10. Testing Approach

MallardOS uses three testing levels.

### 10.1 Host-side unit tests

Use Go unit tests where the code is compatible with the host Go runtime.

Priorities:

- Shell parser.
- Command registry.
- Command argument handling.
- Configuration parsing.
- Service lifecycle logic.

### 10.2 Hardware testing

All hardware-dependent features must be verified on the Seeed XIAO ESP32-S3 N16R8.

Examples:

- USB Serial/JTAG.
- GPIO.
- Multicore scheduling.
- PSRAM.
- Flash storage.
- Wi-Fi.
- TLS.

### 10.3 Python serial test harness

Python will provide host-side automation for:

- Connecting to the serial interface.
- Sending shell commands.
- Capturing terminal output.
- Checking expected responses.
- Running repeatable hardware smoke tests.

---

## 11. First Experiment — Multicore Runtime

### Objective

Verify that TinyGo can execute goroutines on both ESP32-S3 LX7 cores using the selected `cores` scheduler.

### Test design

Create two goroutines:

**Worker 1**
- Pinned to core 0.
- Toggles a GPIO.
- Increments a shared counter.

**Worker 2**
- Pinned to core 1.
- Toggles another GPIO.
- Reads the shared counter.

Protect shared state using a mutex.

Both workers should periodically report their activity over USB serial.

### Success criteria

- Both workers execute.
- Core affinity behaves as expected.
- Both GPIOs toggle independently.
- Shared counter remains consistent.
- No deadlocks.
- No unexplained crashes or watchdog resets.
- USB serial remains responsive.

### Result

Pending physical hardware testing.

---

## 12. Development Philosophy

MallardOS follows these principles:

1. Go first, C only when necessary.
2. Hardware access stays behind HAL interfaces.
3. Keep components modular and independently testable.
4. Prefer small dependencies and predictable memory usage.
5. Verify runtime assumptions on physical hardware.
6. Separate system logging from interactive terminal output.
7. Build incrementally, validating each subsystem before moving forward.

The goal is not simply to run Go applications on an ESP32.

The goal is to build a small, modular embedded operating environment with its own terminal, shell, services, hardware management and application architecture.

**MallardOS — An OS with a duck one step at a time**
