# GA36 Android

I'm working on an Android port for my R36S-style handheld with a **GA36-MB V1.2-20260416** board. I started with Android 5.0, moved through Android 6, and now have **Android 7.0 Nougat / API 24** running on the device.

I develop this project **with AI assistance**. I use AI tools to help analyse the firmware, write and adapt code, investigate problems, and prepare documentation. I test the changes on my own handheld. A successful build or an AI-generated explanation is not enough for me to call a feature working — I distinguish local checks from actual device tests below.

## Photos of my device

These are photos of my test handheld and board, not rendered mockups. The Android photos were taken on 6 November 2026 and do not identify or validate every later build.

<img src="assets/photos/android7-device.jpg" alt="GA36 handheld showing Android 7.0 and vendor kernel 3.4.39 #522" width="360">

*Android 7.0 running on the GA36 handheld, using the vendor kernel shown on screen.*

<details>
<summary>Memory view and board photos</summary>

<img src="assets/photos/android7-memory.jpg" alt="Android memory settings on the GA36 test handheld" width="320">

*Memory settings during that test: about 496 MB reported total, 304 MB average use and 191 MB free. This is a snapshot, not a performance benchmark.*

<img src="assets/photos/board-component-side.jpg" alt="GA36-MB V1.2-20260416 component side with chips and display connector" width="360">
<img src="assets/photos/board-button-side.jpg" alt="Opposite side of the GA36 board with button contacts, card slots and USB ports" width="360">

*Both sides of my GA36-MB V1.2-20260416 board. Printed chip markings differ from the runtime platform identification described below.*

</details>

## My hardware

Runtime analysis identifies my board as **Allwinner A33 with 512 MB RAM**. The photographed chips have Rockchip markings, which differ from that runtime identification.

So far, I've tested this on **one GA36-MB V1.2-20260416 board**. I can't confirm compatibility with other revisions or other A33 handhelds.

## Where the project stands

**Updated 8 November 2026: V69 is the latest installed Android 7 development build confirmed running. It now uses my source-built D03 Linux 3.4.39 kernel (#70). V70 wake/display work is in preparation; it is not an accepted release.**

V69 was installed and checked by direct readback. I reported that it boots and runs; fresh device logs separately confirm the D03 kernel identity and the native system-selection menu before Android starts. This supersedes the earlier V65/vendor-kernel development status below in dated reports.

The native menu offers Android 7 and a **locked Android 11 placeholder**. Real input events and an Android 7 selection were recorded before Android init. Android 7 is the configured default after five seconds; navigation pauses the countdown. **Android 11 has not been built as a complete system or booted.** This is not a working Android 7/11 dual-boot release.

<img src="assets/illustrations/system-selection.jpg" alt="Illustration of the GA36 system-selection menu with Android 7 enabled and Android 11 locked" width="360">

*User-supplied menu illustration. Runtime confirmation comes from device logs; this image does not establish that Android 11 works. The menu currently uses German labels; this README is maintained in English for international readers.*

Earlier testing confirmed Android/apps, audio, gamepad navigation, both sticks, right-stick mouse with R1 click and Power-button operation on the earlier development baseline. Those results are historical; I still distinguish them from complete per-feature acceptance of the new kernel and startup path.

**Known current issue:** waking after Power can intermittently leave the display black and require a reset. V70 preparation addresses the wake/display path and the repeated logo/dark-handoff sequence. Earlier successful wake tests are not a claim that this regression is fixed.

Device startup logs show **256 MiB zRAM and a 4 GiB SD swap file active**, with priorities 100 and 10. Both were unused in the captured snapshots; runtime load/performance testing is still needed. Physical RAM remains 512 MiB.

**PC USB, OTG hub/flash-drive support and secondary-SD operation remain unconfirmed.** USB investigation is deferred. Android 5 and 6 remain earlier development stages, with backups kept locally.

See the [latest V69/D03 progress report](docs/updates/2026-11-08/v69-d03-status.md) and [current status summary](docs/updates/2026-11-08/latest-status.json). Full cold-start/restart, audio, controls and suspend/wake regression acceptance remains open.

## Own kernel — D03 running, sources kept local for now

D03 (#70) is installed and its identity is confirmed by device logs. It incorporates the diagnostic work on the vendor partition layout, audio machine-driver integration and the analog UART receive path. Static checks cover 377 versioned imports across five board modules, without CRC overrides; some vendor modules are still retained for local testing.

Boot success does not complete graphics, audio, analog-control or suspend/wake validation. Kernel source archives and build packages remain outside this GitHub publication for now; source/rights work for retained components is still open.

## BusyBox preparation — local only

I have built a static ARM BusyBox 1.36.1 replacement, including compatibility with the older `timeout -t` calls used by my helpers. **17 ARM-emulation checks passed**, and the saved build recipe also produced a second working local build. Sources, configuration, changes and license notices are prepared locally. Device integration is still pending; this is not a promise of compatibility with every command in the older BusyBox.

## Recovery experiments

I've also worked on an experimental TWRP build and a GA36-specific recovery configuration. Local ARM emulator self-tests passed, but I haven't confirmed a working recovery UI on the device.

I've **paused TWRP development indefinitely** for now. After the recovery tests, I restored the V38 Android startup path and verified it by direct card readback. The V38 direct-start layout is historical context; the current V69 test uses the D03 kernel and native pre-Android menu.

## What you'll find here

The implementation files I've published so far mainly come from the **earlier Android 5 bring-up**:

- `src/ga36-pointer.c`: ARM32 evdev-to-uinput mouse helper.
- `src/lights-ga36.c`: experimental Android 5 ARM32 backlight HAL.
- `src/test-kernel-branches.c`: isolated ARM test for the Power-event suppression change.
- `config/`: board-specific input mappings and device policy.
- `docs/power-button.md` and `docs/usb-status.md`: earlier technical findings.

I've added Android 7 progress reports and selected validation metadata. The full Android 7, custom-kernel and recovery implementations are **not published here yet**.

**This remains experimental. This repository is not a ready-to-flash release or a complete installation package. I haven't included stock firmware, Android images, game ROMs, full card images or private keys.**

## Development reports

These are dated snapshots of the work. They distinguish device-confirmed results from offline checks and pending tests. Older reports describe historical builds and must not be read as current device acceptance.

- [Latest V69/D03 startup, memory and Android 11 preparation update](docs/updates/2026-11-08/v69-d03-status.md)
- [Current status summary](docs/updates/2026-11-08/latest-status.json)
- [Earlier source, validation and release-status update](docs/updates/2026-11-08/source-and-validation-update.md)
- [Earlier V65 progress and pending tests](docs/updates/2026-11-08/progress-update.md)
- [Earlier V65 curated status](docs/updates/2026-11-08/status.json)
- [Image release readiness and third-party software](docs/release-readiness.md)
- [Android 7 progress update](docs/updates/2026-11-06/progress-update.md)
- [Android 7 status and evidence categories](docs/updates/2026-11-06/android7-status.json)
- [Custom kernel V39 checks](docs/updates/2026-11-06/kernel-v39-checks.json)
- [TWRP status](docs/updates/2026-11-06/twrp-status.json)
- [Evidence index](docs/updates/2026-11-06/evidence-index.json)

## Planned public test release

I plan to use my **own source-built kernel** for the public Android 7 test image. D03 has now booted on the handheld, but full hardware acceptance and the retained-component source/license work remain open. **TWRP will be optional**, remains paused, and will not block the release. Installation and recovery must work by writing the SD card without TWRP.

The [release roadmap](docs/release-readiness.md) covers kernel validation, Android integration, source/license preparation, a clean image, testing on another card and an experimental GitHub pre-release. No public image is available yet.

## What I want to work on next

- Resolve the intermittent black-screen wake issue and validate controls, audio, display and Power on the D03/V69 baseline.
- Investigate PC USB and OTG hub/flash-drive operation.
- Finish full device regression tests and source/license preparation for the custom-kernel integration.
- Put together reproducible build instructions and a board-specific installation and recovery process.

If you have matching hardware, board documentation or relevant driver knowledge, I'd welcome your findings. Please don't assume that sharing the same SoC makes another device compatible with this port.

## Publication checks and attribution

I keep the publication separate from my full local working directory. `PUBLICATION-MANIFEST.json` records the reviewed files, and `tools/check-publication.py` checks the file list, UTF-8/JSON text, explicitly reviewed JPEG photos/illustrations, sizes and SHA256 hashes. It rejects extra files and symlinks while excluding local root `.git` metadata.

Run it with:

```sh
python3 tools/check-publication.py
```

New or changed files need review and a corresponding manifest update. The manifest isn't signed, and the checker isn't a guarantee of legal clearance or a complete secret scan.

See [NOTICE.md](NOTICE.md) for provenance and licensing status. I have selected MIT for my own original code; the scoped license package is separate from this documentation update. This does not relicense imported files or grant permission to redistribute third-party firmware.

## Related post

[My Reddit post](https://www.reddit.com/r/R36S/comments/1xomg76/comment/19s6oe904jeyt/)
