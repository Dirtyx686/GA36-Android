# GA36 Android

I'm working on an Android port for my R36S-style handheld with a **GA36-MB V1.2-20260416** board. I started with Android 5.0, moved through Android 6, and now have **Android 7.0 Nougat / API 24** running on the device.

I develop this project **with AI assistance**. I use AI tools to help analyse the firmware, write and adapt code, investigate problems, and prepare documentation. I test the changes on my own handheld. A successful build or an AI-generated explanation is not enough for me to call a feature working — I distinguish local checks from actual device tests below.

## My hardware

Runtime analysis identifies my board as **Allwinner A33 with 512 MB RAM**. The photographed chips have Rockchip markings, which differ from that runtime identification.

So far, I've tested this on **one GA36-MB V1.2-20260416 board**. I can't confirm compatibility with other revisions or other A33 handhelds.

## Where the project stands

**V41 is currently running on my device**, using the V38 direct-boot baseline and the existing vendor kernel. The functions I currently use are working in my initial testing:

- Android startup and apps
- Audio
- Gamepad menu navigation
- Right-stick mouse control with R1 click
- Display off and wake
- Power-button operation

I still need longer testing and further checks after changes. **PC USB connectivity and OTG hub/flash-drive support remain unresolved.**

Android 7 is my permanent boot target. Android 5 and 6 are no longer bootable systems on the test card, but I have kept earlier states in local backups.

V41 includes a built-in speaker routing guard and a dark display handoff. It retains the existing boot baseline and vendor kernel. I haven't installed my separate custom kernel yet.

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

**This repository is not a ready-to-flash release or a complete installation package. I haven't included stock firmware, Android images, game ROMs, full card images or private keys.**

## Development reports

These are dated snapshots of the work. They distinguish device-confirmed results from offline checks and pending tests. They predate my latest V41 confirmation, so this README records the newer overall status.

- [Android 7 progress update](docs/updates/2026-11-06/progress-update.md)
- [Android 7 status and evidence categories](docs/updates/2026-11-06/android7-status.json)
- [Custom kernel V39 checks](docs/updates/2026-11-06/kernel-v39-checks.json)
- [TWRP status](docs/updates/2026-11-06/twrp-status.json)
- [Evidence index](docs/updates/2026-11-06/evidence-index.json)

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
