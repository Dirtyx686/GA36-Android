# Experimental image release readiness

Status: 10 November 2026. **No public Android SD image is provided.** Android 11 P23 with K09 #80 is installed/readback-verified and reported booting; earlier P22 logs confirm framework/GUI startup. Settings navigation is now user-confirmed. The mouse remains unresponsive, fresh P23 detailed logs are pending, and complete runtime acceptance is open. Android 7 D04 #72 is preserved. See the [current report](updates/2026-11-10/p23-k09-status.md).

## Release roadmap

My planned public test image uses a **source-built GA36 kernel**. Android 11 now uses K09 #80; Android 7 retains D04 #72. Retained modules and complete hardware/source acceptance remain release work. **TWRP is optional and is not a release requirement.**

1. **Complete hardware/input acceptance.** Diagnose the P23 mouse failure and test digital/analog controls, sound, Power/wake, repeated starts/restarts and graphics. Settings navigation is confirmed, while complete Settings functionality is not.
2. **Stabilize the Android 11 integration while preserving Android 7.** Capture fresh P23 logs, investigate DocumentsUI/external-storage hangs and measure memory/I/O/boot performance. P24 is planned, not implemented. The last detailed logs are P22; no stable-release claim follows from basic boot success.
3. **Prepare reproducible sources and licenses.** Publish the kernel source, configuration, board changes and build instructions. Establish redistribution terms for every included third-party component and satisfy applicable source, notice and license obligations. A source-built kernel does not resolve rights for the remaining binaries.
4. **Create a clean SD image.** Use newly prepared empty user data. Exclude private data, games, unreviewed game artwork and installation remnants from every partition and staging area.
5. **Test installation on another card.** Document the exact supported board, card capacity requirements, writing and verification steps, first boot and recovery by rewriting the SD card. Verify cold start, sound, controls, wake and the Power menu.
6. **Publish an experimental GitHub pre-release.** Include the verified image only after the above gates pass, SHA256 checksums, source/build information, installation/recovery instructions, known limitations and all required licensing notices. State the experimental status and applicable warranty limitations without claiming blanket immunity.

### Optional TWRP package

TWRP can be developed and tested separately and offered later as an optional package. Its development is still paused indefinitely; this roadmap does not claim that work has resumed or that recovery works on hardware. The SD installation and recovery instructions must remain usable without TWRP.

## Component review findings

- The Android 7 donor base came from an Android SDK system-image package. SDK package terms and the separate open-source component licenses need individual review; the presence of AOSP code does not clear the complete archive. [Android SDK terms](https://developer.android.com/studio/terms).
- Older Mali/EGL/gralloc binaries and donor audio/ION libraries have identifiable upstream repository matches, but matching redistribution terms and complete source correspondence have not been established. Repository availability alone is insufficient. [Arm driver downloads](https://developer.arm.com/Downloads/-/Mali%20GPUs%20User-Space%20Binary%20Drivers).
- The vendor kernel/modules and boot components still require matching source/license evidence. The old system BusyBox was matched byte-for-byte to the official 1.20.0 ARMv6 binary; a separate source-built BusyBox 1.36.1 candidate has sources/configuration/notices prepared locally and 17 passing emulation checks, but is not installed or device-accepted. The earlier V39/D03 stages are historical. K09 #80 now boots Android 11 locally and D04 #72 is retained for Android 7; full corresponding-source/rights preparation for retained components remains incomplete. [BusyBox licensing](https://busybox.net/license.html).
- Amaze has an identifiable official release and GPL licensing; inclusion still requires the applicable source and notice obligations. It does not clear unrelated firmware components. [Amaze source and releases](https://github.com/TeamAmaze/AmazeFileManager).

Missing evidence is an unresolved release condition, not permission to redistribute. No blanket third-party redistribution permission is claimed. MIT has been selected for my own original code; imported material retains its own terms. The scoped license package is separate from this documentation update; see [NOTICE](../NOTICE.md).

The project remains experimental, without a promise of fitness, compatibility or reliability. An experimental label or warranty disclaimer does not replace licensing obligations or remove liability that cannot legally be excluded.

## Latest review outcome

The current graphics, audio and boot components are being retained technically. Matching redistribution evidence for the existing Mali and audio files remains incomplete; GPL A33 Boot0 source was found, but complete corresponding sources for the modified U-Boot remain unestablished. Keeping these components does not mark the image release gates complete. See the [earlier rights review](updates/2026-11-08/source-and-validation-update.md) and [current P23 status](updates/2026-11-10/p23-k09-status.md). No image, kernel or BusyBox binary is included in this documentation update. Android 11 has since booted, but P23 still has unresolved input/stability issues and is not a finished public release.
