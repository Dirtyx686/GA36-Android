# GA36 Android

I'm working on an Android port for my R36S-style handheld with a **GA36-MB V1.2-20260416** board. I started with Android 5.0, moved through Android 6, and now have **Android 7.0 Nougat / API 24** running on the device.

I develop this project **with AI assistance**. I use AI tools to help analyse the firmware, write and adapt code, investigate problems, and prepare documentation. I test the changes on my own handheld. A successful build or an AI-generated explanation is not enough for me to call a feature working — I distinguish local checks from actual device tests below.

## My hardware

Runtime analysis identifies my board as **Allwinner A33 with 512 MB RAM**. The photographed chips have Rockchip markings, which differ from that runtime identification.

So far, I've tested this on **one GA36-MB V1.2-20260416 board**. I can't confirm compatibility with other revisions or other A33 handhelds.

## Where the project stands

**As of 8 November 2026, V65 is my latest local Android 7 build.** I have written it to my test card and verified the written contents by direct readback. **The physical V65 device test is still pending.** My V56 cold start was confirmed on the handheld; that is a separate, earlier result.

Earlier device testing established:

- Android startup and apps
- Audio
- Gamepad menu navigation
- Right-stick mouse control with R1 click
- Display off and wake
- Power-button operation

V65 retains the previously working native audio, wake, Power and graphics files. This does not replace regression testing of the new build. **PC USB connectivity, OTG hub/flash-drive support and secondary-SD operation remain unconfirmed.** USB investigation is currently deferred.

Android 7 is my permanent boot target, using the **V38 direct-start baseline and vendor Linux 3.4.39 #522**, with the existing patches. Android 5 and 6 are no longer bootable systems on the test card; I have kept earlier states in local backups.

V65 adds launcher layout/theme options and an X+Y background-app cleanup action. Local behavior tests pass, but actual controller delivery and first-boot maintenance still need device testing. I removed overlapping automatic installation/compilation work; I have not measured a V65 performance improvement or boot-time reduction.

See the [latest progress report](docs/updates/2026-11-08/progress-update.md) for the test boundaries.

## My custom kernel work

I've compiled a separate **V39 development kernel** locally. Static compatibility checks pass for five modules: `disp`, `mali`, `lcd`, `gpio-sunxi` and `udt_joystick`. All **377 versioned imports** match, without CRC overrides.

That doesn't mean the kernel works on the handheld yet. The analog input receive path is still open, and I haven't tested boot, graphics, controls, audio, wake or Power behavior with this kernel.

**V39 has not been installed on the SD card.** I also haven't measured a performance improvement.

## Recovery experiments

I've also worked on an experimental TWRP build and a GA36-specific recovery configuration. Local ARM emulator self-tests passed, but I haven't confirmed a working recovery UI on the device.

I've **paused TWRP development indefinitely** for now. After the recovery tests, I restored the V38 Android startup path and verified it by direct card readback. My current Android build runs on that boot baseline.

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

- [Latest V65 progress and pending tests](docs/updates/2026-11-08/progress-update.md)
- [Current curated status](docs/updates/2026-11-08/status.json)
- [Image release readiness and third-party software](docs/release-readiness.md)
- [Android 7 progress update](docs/updates/2026-11-06/progress-update.md)
- [Android 7 status and evidence categories](docs/updates/2026-11-06/android7-status.json)
- [Custom kernel V39 checks](docs/updates/2026-11-06/kernel-v39-checks.json)
- [TWRP status](docs/updates/2026-11-06/twrp-status.json)
- [Evidence index](docs/updates/2026-11-06/evidence-index.json)

## Planned public test release

I plan to use my **own source-built kernel** for the public Android 7 test image. The custom kernel still needs analog-input work and hardware validation; the current vendor-kernel build is a development baseline. **TWRP will be optional**, remains paused, and will not block the release. Installation and recovery must work by writing the SD card without TWRP.

The [release roadmap](docs/release-readiness.md) covers kernel validation, Android integration, source/license preparation, a clean image, testing on another card and an experimental GitHub pre-release. No public image is available yet.

## What I want to work on next

- Keep testing controls, audio, display wake and Power behavior on Android 7.
- Investigate PC USB and OTG hub/flash-drive operation.
- Complete the custom kernel's analog input path before moving on to hardware tests.
- Put together reproducible build instructions and a board-specific installation and recovery process.

If you have matching hardware, board documentation or relevant driver knowledge, I'd welcome your findings. Please don't assume that sharing the same SoC makes another device compatible with this port.

## Publication checks and attribution

I keep the publication separate from my full local working directory. `PUBLICATION-MANIFEST.json` records the reviewed files, and `tools/check-publication.py` checks the file list, UTF-8/JSON content, sizes and SHA256 hashes. It rejects extra files and symlinks while excluding local root `.git` metadata.

Run it with:

```sh
python3 tools/check-publication.py
```

New or changed files need review and a corresponding manifest update. The manifest isn't signed, and the checker isn't a guarantee of legal clearance or a complete secret scan.

See [NOTICE.md](NOTICE.md) for provenance and licensing status. I haven't selected a general project license yet. Nothing in this repository grants permission to redistribute third-party firmware.

## Related post

[My Reddit post](https://www.reddit.com/r/R36S/comments/1xomg76/comment/19s6oe904jeyt/)
