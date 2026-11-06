# GA36 Android — Android 7.0 progress

Experimental Android bring-up for the **GA36-MB V1.2-20260416** R36S-style handheld. The project has progressed from Android 5.0 through Android 6 to **Android 7.0 Nougat / API 24**.

Runtime analysis identifies the platform as **Allwinner A33 with 512 MB RAM**. Photographed chip markings differ from the runtime identification; donor model strings alone are not treated as hardware identification. Compatibility with other board revisions or A33 devices is not established.

**This is a development repository, not a ready-to-flash firmware release. No stock firmware, Android images, game ROMs, card images or private keys are included.**

## Current status — 6 November 2026

| Component | Status |
| --- | --- |
| Current Android build | V41 running on the V38 direct-boot baseline; owner reports the currently used functions are working; existing vendor kernel retained |
| Custom kernel | V39 compiled locally; not installed or boot-tested |
| Recovery | Experimental TWRP development paused indefinitely at the owner's request; working device UI not confirmed |
| PC USB / OTG hub and flash drive | Unresolved |

Android 7 is the permanent boot target. Android 5 and 6 are no longer bootable systems on the test card; earlier states remain in local backups.

The owner reports that **V41 is running with the currently used functions working**. Development has established Android startup and apps, audio, gamepad navigation, right-stick mouse control with R1 click, display wake and Power-button operation. This is an initial device report, not an endurance test or a complete compatibility assessment. PC USB and OTG operation have not been separately confirmed.

V41 includes a built-in speaker routing guard and a dark display handoff. It retains the known V38 boot baseline and existing vendor kernel; the custom V39 kernel is not installed.

## Custom kernel V39

The separate local build passes static compatibility checks for five modules: `disp`, `mali`, `lcd`, `gpio-sunxi` and `udt_joystick`. All **377 versioned imports** match, with no CRC overrides.

The analog input receive path is still open. Boot, graphics, controls, sound, wake and Power behavior have not been validated with this kernel. **V39 remains uninstalled.**

## Recovery work

A GA36-specific experimental TWRP build passed local ARM emulator self-tests, but a working recovery UI on the handheld has not been confirmed. Recovery development is paused indefinitely.

The Android V38 startup path was restored and verified by direct readback after recovery testing.

## Development reports

The linked reports are earlier, dated development snapshots. They distinguish device-confirmed results from offline validation and pending hardware tests. They predate the latest V41 owner confirmation; this README records the newer status.

- [Android 7 progress update](docs/updates/2026-11-06/progress-update.md)
- [Android 7 status and evidence categories](docs/updates/2026-11-06/android7-status.json)
- [Custom kernel V39 checks](docs/updates/2026-11-06/kernel-v39-checks.json)
- [TWRP status](docs/updates/2026-11-06/twrp-status.json)
- [Evidence index](docs/updates/2026-11-06/evidence-index.json)

## Published source code

The published implementation files currently cover the **earlier Android 5 bring-up**, not a complete Android 7 build:

- `src/ga36-pointer.c`: ARM32 evdev-to-uinput mouse helper.
- `src/lights-ga36.c`: experimental Android 5 ARM32 backlight HAL.
- `src/test-kernel-branches.c`: isolated ARM execution test for the Power-event suppression change.
- `config/`: board-specific input mappings and device policy.
- `docs/power-button.md` and `docs/usb-status.md`: earlier technical findings.

The full Android 7, kernel and recovery implementations are not included. These components are not a complete installer or a generic port for other A33 devices.

## Remaining work

- Continue regression testing of controls, audio, display wake and Power behavior on Android 7.
- Investigate PC USB and OTG hub/flash-drive operation.
- Complete the custom kernel's analog input path and perform hardware tests when authorized.
- Prepare reproducible build instructions and a board-specific installation and recovery process.

## Publication checks and provenance

The project was developed with AI assistance and tested on the owner's device. See [NOTICE.md](NOTICE.md) for provenance and licensing status. No general project license has been selected, and no third-party firmware redistribution permission is implied.

`PUBLICATION-MANIFEST.json` and `tools/check-publication.py` describe the original Android 5 publication snapshot. They have not yet been updated for this revised README and the added Android 7 reports; the original checker will therefore flag these changes. It is not a complete audit of the current repository.

## Related post

[Owner's Reddit post](https://www.reddit.com/r/R36S/comments/1xomg76/comment/19s6oe904jeyt/)
