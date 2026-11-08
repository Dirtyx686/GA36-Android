# Source and validation update — 8 November 2026

I am developing this experimental Android 7.0 / API 24 port with AI assistance and testing it on one GA36-MB V1.2-20260416 handheld. Runtime identification reports Allwinner A33 and 512 MB RAM. This update publishes documentation only.

## Device-confirmed results

I confirmed both sticks, centre positions, full travel and reliable analog input after cold start on the version currently on my card. Previous device tests also confirmed Android/apps, audio, gamepad navigation, right-stick mouse with R1 click, display wake and Power operation. These confirmations must not be read as a complete regression acceptance of every V65 change or as a test of the separate source-built kernel.

PC USB, OTG hub/flash-drive support and secondary-SD operation remain unconfirmed. Android 7 remains the permanent boot target on the V38 direct-start/vendor-kernel baseline.

## Local development, not installed

- A bounded UART1 analog receive candidate has been implemented in the separate source-built kernel. Local build and parser/error-path tests pass; hardware behavior is still untested.
- A GPL `gpio-sunxi` source candidate builds and passes static interface checks. Exact correspondence to the old vendor binary is not proved; GPIO behavior still needs device tests. Licensed modern LCD and joystick references need adaptation before they can replace the old modules.
- A static BusyBox 1.36.1 candidate builds with legacy `timeout -t` compatibility. Seventeen ARM-emulation checks pass, covering basic shell/file/archive/compression operations, timeout success/failure behavior, and syntax checks of eight stored helper scripts. The saved build recipe also builds in a new directory. Sources, configuration, changes and license notices are prepared locally. It has not been installed or device-tested, and some old applets are absent.
- The early BusyBox 1.18.3 identified in an older initramfs is not present in the saved V65 direct-start root. It is not automatically a dependency of the current release plan.

Kernel and BusyBox source/build packages remain local for now. No binary or installable package accompanies this update. TWRP remains optional and paused indefinitely.

## Provenance and licensing review

The old Android `busybox-smp` matches the [official BusyBox 1.20.0 ARMv6 binary](https://busybox.net/downloads/binaries/1.20.0/busybox-armv6l) byte-for-byte. Its [build manifest](https://busybox.net/downloads/binaries/1.20.0/MANIFEST), BusyBox source, uClibc source and Aboriginal Linux 1.2.0 build scripts have been saved locally. Historical build correspondence is distinct from the complete source/configuration package prepared for the new replacement.

I am retaining the present Mali, audio and boot components technically. That decision does not complete their redistribution evidence:

- The exact Mali r4p0 file still lacks an established accompanying redistribution grant. [Arm's package-specific licensing information](https://developer.arm.com/Downloads/-/Mali%20GPUs%20User-Space%20Binary%20Drivers) does not automatically license an unrelated binary.
- Matching Astar audio sources and a matching lowpass archive were found in another vendor tree. Missing licensing evidence for the actual lowpass/codec/RIL chain remains unresolved. Licensed older audio sources are possible replacement references, not a tested A33 replacement.
- [Allwinner's published A33 Boot0 sources](https://github.com/allwinner-zh/bootloader/tree/e5ceeca211883d9f6b25f84e0e3c5fe2afaaf350) include GPL notices. The saved boot chain also identifies a modified U-Boot 2017.09 build; its complete corresponding source/configuration is not established by the older public U-Boot 2011.09 tree.

Targeted source/license requests have been drafted locally, not sent. MIT is selected for my own original code only; imported material keeps its existing terms. A warranty disclaimer or experimental label does not settle missing third-party rights.

## Publication boundary

This is still an experimental project, not a ready-to-flash release. No stock firmware, Android/SD images, kernel files, BusyBox binaries, games, game assets, private data or raw device logs are added by this update. Full image release remains pending the source/license and hardware-test gates in the [release roadmap](../../release-readiness.md).

Earlier dated reports remain historical snapshots. This note supersedes their open-analog-path and BusyBox-provenance statements only to the extent described above; it does not turn offline checks into device acceptance.
