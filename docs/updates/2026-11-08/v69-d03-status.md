# V69 / D03 progress update — 8 November 2026

I develop this GA36-MB V1.2-20260416 Android port with AI assistance and test it on one handheld. Runtime identification reports Allwinner A33 and 512 MiB physical RAM. This is an experimental development update, not a flashable release.

## Current installed baseline

V69 was installed with direct card readback verification. The installation retained the existing D03 #70 kernel and did not write the bootloader or userdata. My recorded feedback confirms general boot/run success. Fresh logs identify Linux 3.4.39, built as `ga36@android7-d03`, build #70. Kernel/source packages remain local.

Earlier build-time reports marked D03 hardware tests pending. Later boot logs now establish that D03 actually runs; they do not by themselves establish every hardware function or long-term stability.

## Native menu before Android

A native framebuffer/evdev menu is launched before Android init. The real device log records a 640x480 menu, navigation to the locked Android 11 entry, navigation back to Android 7 and confirmation of Android 7. It then exits to the D03 Android handoff with Android 11 not started.

The configured default is Android 7 after five seconds. Navigation pauses the countdown; a longer inactivity fallback continues to Android 7. The exact automatic timing paths still need separate acceptance testing. This is early userspace, not a newly flashed bootloader.

The user-supplied menu illustration in the README is labelled as an illustration. It supplements the description; runtime claims are based on logs and recorded user feedback.

## Memory configuration

Device startup logs confirm 256 MiB zRAM at priority 100 and a 4 GiB SD-backed swap file at priority 10. Both reported zero usage at the captured startup snapshots. These values describe configured swap, not additional physical RAM or measured gaming performance. Full runtime load and SD-I/O testing remain open.

## Wake regression and V70

Current recorded feedback reports intermittent failure to restore the display after Power, requiring a reset. V70 is being prepared to investigate/fix the wake/display path and remove the extra short native-logo preview that contributes to the logo/dark-handoff sequence. No V70 installation or successful V70 device acceptance is established by the reviewed preparation record.

Earlier working audio, controls and wake reports refer to earlier test states. A complete D03/V69 regression of sound, buttons/sticks, mouse, Power menu, repeated cold starts/restarts and suspend/wake still needs to be recorded.

## Android 11 / LineageOS 18.1

A local ARM32 device/vendor configuration, interface inventory and separate kernel/Binder preparation exist. Local records describe linked protocol-7 and protocol-8 kernel variants and static module checks. None of this establishes a bootable Android 11 userspace.

The complete Android 11 source tree/system build has not been completed. Graphics/composer, audio/power interfaces, security policy, Binder integration, storage separation and hardware tests remain work items. The Android 11 menu entry stays locked; no working Android 7/11 dual boot is claimed.

## Publication boundary

This update publishes curated English documentation and a labelled menu illustration. It does not publish a kernel, Android image, bootloader, APK, toolchain, game, game asset, raw runtime log or personal data. The earlier BusyBox source-built candidate remains locally prepared and emulator-tested; integration into the current card is not claimed here.

Mali/audio redistribution evidence and complete corresponding sources for retained vendor/boot components remain open. Retaining those components technically does not clear a full-image release. TWRP remains optional and paused; PC USB, OTG and secondary-SD operation remain unconfirmed.

## Evidence used locally

- V69 install/readback report and recorded general boot/run feedback.
- V69 pre-Android menu log and fresh D03 kernel/bootstrap identity records.
- zRAM and SD-swap startup snapshots.
- V70 preparation status recording the intermittent wake issue.
- Android 11 local preparation/status record.

Only this curated summary is published. No raw logs or private paths accompany it. Earlier dated reports remain historical snapshots; the current status is summarised in [latest-status.json](latest-status.json).
