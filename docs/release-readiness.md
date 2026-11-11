# GA7 / GA11 release preparation

Status: **11 November 2026**. Exact target: **GA36-MB V1.2-20260416**, R36S-style
clone, runtime Allwinner A33 and 512 MiB RAM. Other R36S hardware is not validated.
No Android image is public here. The separate source/release repository is
being prepared locally and has not yet been created or published.

## Prepared

- GA7 Beta V1 and GA11 Alpha V1 clean P39 image candidates; K11 #82.
- Fresh userdata without accounts or game APK/OBB packages, separate OS images,
  minimum 31,138,512,896 bytes, partition/filesystem expansion for larger cards.
- Full image decompression/hashes, filesystem checks, swap allocation checks,
  ARM growth-helper checks and DE/EN launcher checks passed on the host.
- Native Desktop Flasher 4.0; 30 file-based tests and Linux/Windows build/static
  checks passed. Linux operation is owner-reported, with operation unspecified.
- Complete local K11 source/configuration collection and a pinned Android 11
  manifest: 863 upstream projects, 13 with collected source edits.

## Remaining release work

1. Test each clean image on another card: first boot, full capacity, launcher,
   sound, gamepad, repeated suspend/wake, Power menu and APK/OBB import.
2. Prepare the Linux-only flasher first. A Windows binary release remains
   deferred and is not part of that initial package.
3. Complete matching source/provenance and notice work for retained boot,
   Mali/EGL/gralloc, audio and kernel module components. Do not infer permission
   from another repository merely containing those binaries.
4. Review/remove unestablished embedded game-logo resources from any proposed
   public launcher/image. Existing test images are not silently changed by
   preparation of a source-only public tree.
5. Verify the source/build packages and document absent inputs. A collected
   source snapshot is not a demonstrated clean-room or bit-identical ROM build.
6. Present the candidate repository and release asset list to the owner before
   making the new release project and ROM downloads public.

The latest owner feedback confirms Android 11 sound, GTA trailer playback and
that the exit crash is no longer observed; exact tested build was unspecified.
This supersedes an entirely pending sound test, but does not complete all
regression or fresh-image tests. USB-to-PC/OTG and slot 2 remain deferred.

TWRP remains optional and paused. Recovery must work by restoring a known-good
SD image without TWRP. MIT applies to original code only; see [NOTICE](../NOTICE.md).
