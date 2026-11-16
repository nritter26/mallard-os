# Mallard-OS v0.1

*TinyGo-first embedded operating environment for the ESP32-S3 (N16R8). Status: work in progress.*

## 1. Goals

An embedded operating environment with:

- Interactive terminal over USB serial
- Command interpreter (shell) with a command registry
- Persistent filesystem
- Task (goroutine) management and a service supervisor
- Hardware abstraction layer (HAL)
- GPIO, PWM, ADC and I²C management
- Basic networking (Wi-Fi, HTTP)
- Eventually, an application management system

## 2. Language strategy

| Language | Role |
|---|---|
| Go (TinyGo) | Primary: shell, services, HAL, drivers, networking, supervisor |
| C | Gap-filler only, behind HAL interfaces (see section 5) |
| Python | Host side only: serial client, test harness, build and flash scripts |

Everything that touches C sits behind a small Go interface (`Storage`, `Crypto`, `Power`), so a C implementation can be swapped for a Go one later without touching the shell or services.

## 3. Stack

| Technology | Responsibility |
|---|---|
| TinyGo (0.44 or later) | Compiler and bare-metal runtime |
| `cores` scheduler | Goroutines on both LX7 cores, with core affinity |
| espradio | Wi-Fi and Bluetooth, radio on core 0 |
| lneto | Network stack used by espradio |
| tinyfs (LittleFS) | Filesystem, C library wrapped for Go |
| `machine` package + tinygo drivers | GPIO, PWM, ADC, I²C, SPI |
| USB (CDC serial) | Terminal communication |
| espflasher | Flashing without external tools |

## 4. Architecture

```
                 Mallard-OS
                     |
        +------------+------------+
        |            |            |
      Shell       Services     Applications
        |            |            |
   Command API   System API   Application API
        |            |            |
        +------------+------------+
                     |
          Hardware Abstraction (Go)
                     |
          TinyGo runtime (cores scheduler)
                     |
                  ESP32-S3
```

Core split (initial plan): core 0 runs the radio; shell and services run on core 1 via goroutine affinity.

## 5. C gaps and risks

Status column is updated as each item is tested on the board.

| Area | Why C may be needed | Status |
|---|---|---|
| LittleFS | tinyfs wraps the C LittleFS library through CGo; needs the S3 flash `BlockDevice` | Untested |
| PSRAM (8 MB) | No confirmation TinyGo initializes or uses it | Untested |
| USB Serial/JTAG console | USB-OTG device support exists; built-in Serial/JTAG coverage unconfirmed | Untested |
| TLS / HTTPS | Pure-Go TLS story with lneto unclear; mbedTLS is the fallback | Untested |
| Hardware crypto | Go crypto runs in software unless a driver exists | Untested |
| OTA and partition table | ESP-IDF handles this in C; TinyGo behaviour unverified | Untested |
| Sleep / low power / ULP | IDF C APIs | Untested |
| Other peripherals (I2S, RMT, LCD, camera) | May exist only in C | Untested |

Test order for the unknowns: (1) USB serial console, (2) LittleFS over `BlockDevice`, (3) PSRAM allocation, (4) TLS client connection.

## 6. Known design constraints

- Scheduling is cooperative, so services must yield or block on channels, or a tight loop will starve the shell.
- TinyGo supports only part of the standard library; avoid reflection-heavy packages and measure binary size.
- Shared state across cores needs mutexes or channels.
- Recoverable panics (TinyGo 0.42) are the basis for supervised service restarts; verify they behave under the multicore scheduler.

## 7. Non-goals (v1.0)

- No MMU, process isolation or dynamic loading
- No on-device second runtime (such as MicroPython) unless added deliberately in v0.7
- No Linux compatibility

## 8. Roadmap

| Version | Deliverable |
|---|---|
| v0.1 | Boot sequence + serial shell (`help`, `sysinfo`, `uptime`, `meminfo`) |
| v0.2 | Command registry + built-in commands |
| v0.3 | Filesystem + configuration storage (`ls`, `mkdir`) |
| v0.4 | Task manager + system monitoring (`ps`, with core shown) |
| v0.5 | Hardware management services (`gpio list`) |
| v0.6 | Wi-Fi + HTTP server (`wifi scan`) |
| v0.7 | Application management |
| v1.0 | First usable Mallard-OS release |

## 9. Target commands by version

| Command | Version |
|---|---|
| `help`, `sysinfo`, `uptime`, `meminfo` | v0.1 |
| `ls`, `mkdir projects` | v0.3 |
| `ps` | v0.4 |
| `gpio list` | v0.5 |
| `wifi scan` | v0.6 |

## 10. Testing approach

- Host-side Go unit tests for the parser and command registry (no hardware needed)
- Python serial harness that drives the shell and asserts on output
- Per-peripheral smoke tests on the N16R8 before each roadmap milestone

## 11. First experiment

Two goroutines pinned to different cores, each toggling a GPIO and printing over serial, with a shared counter protected by a mutex. If this runs cleanly, the multicore design is viable.
