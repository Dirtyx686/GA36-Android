# Experimental image release readiness

Status: 8 November 2026. **No public Android 7 SD image is provided.** V65 is a private test-card build awaiting device acceptance. Cleanup of that card does not make it a sanitized redistribution image.

## Release roadmap

My planned public test image must use my **own source-built GA36 kernel**. The currently running vendor kernel remains a development baseline; it is not the kernel target for this release. **TWRP is optional and is not a release requirement.**

1. **Finish and validate the custom kernel.** Use the V39 development work as a starting point, complete the analog input receive path, and test boot, display/graphics, audio, buttons, sticks, wake and Power on the handheld. V39 is currently built and statically checked only; it has not been installed or device-tested.
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
- The vendor kernel/modules, boot components and BusyBox require applicable license/source compliance checks. The separate V39 development kernel is not established as corresponding source for the running vendor kernel. [BusyBox licensing](https://busybox.net/license.html).
- Amaze has an identifiable official release and GPL licensing; inclusion still requires the applicable source and notice obligations. It does not clear unrelated firmware components. [Amaze source and releases](https://github.com/TeamAmaze/AmazeFileManager).

Missing evidence is an unresolved release condition, not permission to redistribute. No blanket third-party redistribution permission is claimed. No general project license has been selected; see [NOTICE](../NOTICE.md).

The project remains experimental, without a promise of fitness, compatibility or reliability. An experimental label or warranty disclaimer does not replace licensing obligations or remove liability that cannot legally be excluded.
