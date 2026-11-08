# Experimental image release readiness

Status: 8 November 2026. **No public Android 7 SD image is provided.** V69 now runs with the installed source-built D03 kernel (#70), confirmed by device logs and general boot/run feedback. Full regression acceptance remains open. V70 wake/display work is in preparation; Android 11 is still locked and has not booted. See the [latest update](updates/2026-11-08/v69-d03-status.md).

## Release roadmap

My planned public test image must use my **own source-built GA36 kernel**. The current D03 kernel has booted on the device, but retained modules and complete hardware/source acceptance remain release work. **TWRP is optional and is not a release requirement.**

1. **Complete custom-kernel hardware acceptance.** D03 is now installed and runtime-confirmed. Resolve the intermittent black-screen wake issue and complete graphics, audio, digital/analog input, Power, cold-start/restart and memory-load tests. Static module-interface checks are not full acceptance.
2. **Stabilize Android 7 with that kernel.** Check repeated cold starts and longer use. Document unresolved USB/OTG or secondary-SD limitations explicitly; do not advertise them as working.
3. **Prepare reproducible sources and licenses.** Publish the kernel source, configuration, board changes and build instructions. Establish redistribution terms for every included third-party component and satisfy applicable source, notice and license obligations. A source-built kernel does not resolve rights for the remaining binaries.
4. **Create a clean SD image.** Use newly prepared empty user data. Exclude private data, games, unreviewed game artwork and installation remnants from every partition and staging area.
5. **Test installation on another card.** Document the exact supported board, card capacity requirements, writing and verification steps, first boot and recovery by rewriting the SD card. Verify cold start, sound, controls, wake and the Power menu.
6. **Publish an experimental GitHub pre-release.** Include the verified image only after the above gates pass, SHA256 checksums, source/build information, installation/recovery instructions, known limitations and all required licensing notices. State the experimental status and applicable warranty limitations without claiming blanket immunity.

### Optional TWRP package

TWRP can be developed and tested separately and offered later as an optional package. Its development is still paused indefinitely; this roadmap does not claim that work has resumed or that recovery works on hardware. The SD installation and recovery instructions must remain usable without TWRP.

## Component review findings

- The Android 7 donor base came from an Android SDK system-image package. SDK package terms and the separate open-source component licenses need individual review; the presence of AOSP code does not clear the complete archive. [Android SDK terms](https://developer.android.com/studio/terms).
- Older Mali/EGL/gralloc binaries and donor audio/ION libraries have identifiable upstream repository matches, but matching redistribution terms and complete source correspondence have not been established. Repository availability alone is insufficient. [Arm driver downloads](https://developer.arm.com/Downloads/-/Mali%20GPUs%20User-Space%20Binary%20Drivers).
- The vendor kernel/modules and boot components still require matching source/license evidence. The old system BusyBox was matched byte-for-byte to the official 1.20.0 ARMv6 binary; a separate source-built BusyBox 1.36.1 candidate has sources/configuration/notices prepared locally and 17 passing emulation checks, but is not installed or device-accepted. The earlier V39 work is historical; D03 now runs locally, while full corresponding-source/rights preparation for retained components remains incomplete. [BusyBox licensing](https://busybox.net/license.html).
- Amaze has an identifiable official release and GPL licensing; inclusion still requires the applicable source and notice obligations. It does not clear unrelated firmware components. [Amaze source and releases](https://github.com/TeamAmaze/AmazeFileManager).

Missing evidence is an unresolved release condition, not permission to redistribute. No blanket third-party redistribution permission is claimed. MIT has been selected for my own original code; imported material retains its own terms. The scoped license package is separate from this documentation update; see [NOTICE](../NOTICE.md).

The project remains experimental, without a promise of fitness, compatibility or reliability. An experimental label or warranty disclaimer does not replace licensing obligations or remove liability that cannot legally be excluded.

## Latest review outcome

The current graphics, audio and boot components are being retained technically. Matching redistribution evidence for the existing Mali and audio files remains incomplete; GPL A33 Boot0 source was found, but complete corresponding sources for the modified U-Boot remain unestablished. Keeping these components does not mark the image release gates complete. See the [new review update](updates/2026-11-08/source-and-validation-update.md). No image, kernel or BusyBox binary is included in this documentation update. Android 11 preparation and a locked menu entry are not a working Android 11 release.
