# GA36 Android — Android 7.0 progress update

The project has progressed from Android 5.0 through Android 6 to **Android 7.0 Nougat / API 24** on the **GA36-MB V1.2-20260416** handheld. Runtime analysis identifies the platform as **Allwinner A33 with 512 MB RAM**.

## Current Android build

**V38 is the latest device-confirmed Android 7 direct-boot baseline**, using the existing vendor kernel. The owner confirmed that Android 7 starts successfully.

- Android 7 is the permanent boot target.
- Android 5 and 6 are no longer bootable systems on the test card; earlier states remain in local backups.
- Audio was confirmed working after the V33 correction.
- Earlier development established gamepad navigation, right-stick mouse control with R1 click, display wake and Power-button operation. These confirmations do not replace regression tests on later builds.

A newer **V40 userspace optimization update** has been installed on the V38 boot baseline and verified by direct card readback. It reduces unused-package activity and diagnostic overhead without installing the custom kernel. **Physical verification and performance measurements are pending; no speed improvement is claimed.**

## Custom kernel development

A separate **V39 kernel development build** has been compiled locally. Static checks pass for five modules: `disp`, `mali`, `lcd`, `gpio-sunxi` and `udt_joystick`. Their 377 versioned imports match the build, with no CRC overrides.

The analog input receive path and hardware validation remain open. **This kernel has not been installed on the SD card.** Static compatibility does not establish that boot, graphics, controls, audio or wake work with the new kernel.

## Recovery development

An experimental TWRP build and GA36-specific recovery configuration were developed. Local ARM emulator self-tests passed, but a working recovery UI on the handheld has not been confirmed.

**TWRP work is currently paused indefinitely at the owner's request.** The V38 Android startup path was restored and verified by direct readback after recovery testing. The latest installed V40 update still needs a device test.

## Remaining work

- Verify the latest installed Android 7 update and retest controls, audio, wake and Power behavior.
- Continue investigating PC USB and OTG hub/flash-drive support.
- Complete and test the custom kernel's analog input path when hardware testing is authorized.
- Prepare reproducible build instructions and a board-specific installation and recovery process.

## Repository scope

The original published source files document the earlier Android 5 bring-up. This update adds **Android 7 development status and selected validation metadata**, not the full Android 7, kernel or recovery implementation.

This remains an experimental project for one tested board revision, not a ready-to-flash release. **No stock firmware, Android images, vendor binaries, private keys or game ROMs are included in this update.**

See the adjacent JSON reports for the evidence categories and outstanding tests.
