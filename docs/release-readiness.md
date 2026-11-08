# Experimental image release readiness

Status: 8 November 2026. **No public Android 7 SD image is provided.** V65 is a private test-card build awaiting device acceptance. Cleanup of that card does not make it a sanitized redistribution image.

Before I can publish an image, I need to:

1. Establish redistribution terms for every included third-party binary and satisfy applicable source, notice and license obligations.
2. Prepare new empty user data and remove private data, games and unreviewed application resources from every partition and installation staging area.
3. Document the exact supported board, card requirements, write procedure, verification and recovery path.
4. Test that distributable image on another card, including cold start, sound, controls, wake and Power behavior.

## Component review findings

- The Android 7 donor base came from an Android SDK system-image package. SDK package terms and the separate open-source component licenses need individual review; the presence of AOSP code does not clear the complete archive. [Android SDK terms](https://developer.android.com/studio/terms).
- Older Mali/EGL/gralloc binaries and donor audio/ION libraries have identifiable upstream repository matches, but matching redistribution terms and complete source correspondence have not been established. Repository availability alone is insufficient. [Arm driver downloads](https://developer.arm.com/Downloads/-/Mali%20GPUs%20User-Space%20Binary%20Drivers).
- The vendor kernel/modules, boot components and BusyBox require applicable license/source compliance checks. The separate V39 development kernel is not established as corresponding source for the running vendor kernel. [BusyBox licensing](https://busybox.net/license.html).
- Amaze has an identifiable official release and GPL licensing; inclusion still requires the applicable source and notice obligations. It does not clear unrelated firmware components. [Amaze source and releases](https://github.com/TeamAmaze/AmazeFileManager).

Missing evidence is an unresolved release condition, not permission to redistribute. No blanket third-party redistribution permission is claimed. No general project license has been selected; see [NOTICE](../NOTICE.md).

The project remains experimental, without a promise of fitness, compatibility or reliability. An experimental label or warranty disclaimer does not replace licensing obligations or remove liability that cannot legally be excluded.
