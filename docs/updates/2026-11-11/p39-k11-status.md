# P39 / K11 development update — 11 November 2026

This update describes the owner's **GA36-MB V1.2-20260416** R36S-style clone.
Runtime platform: **Allwinner A33, 512 MiB physical RAM**. Other board revisions
and ordinary Rockchip-based R36S devices are not established as compatible.

## Device feedback

Android 7 is the Beta V1 development setup; Android 11 is Alpha V1. The current
baseline is P39 with the source-built K11 Linux 3.4.39 #82 kernel. P39 launcher101
was registered on both systems. A subsequent launcher102 language derivative
is included in the separate clean-image candidates, using an English default
locale and DE/EN strings. That derivative still needs its own physical test.

Latest owner feedback confirms Android 11 sound, says the GTA Vice City exit
crash is no longer observed, and confirms that the GTA trailer plays. The exact
card/build was not identified in that feedback, so it does not certify a fresh
clean-image install or a complete regression test. Intermittent wake behavior,
memory-pressure stability and complete per-game controls still require targeted
checks. The former claim that the sound test is entirely pending is superseded.

## Prepared locally

- Separate GA7 and GA11 P39 clean images, without user accounts or game APK/OBB
  packages; minimum target size 31,138,512,896 bytes (marketed 32 GB or larger).
- Flasher partition adjustment and native first-boot filesystem growth for
  larger cards. Full decompression/SHA-256/filesystem checks passed on the host;
  another-card first boot, growth and importer acceptance remain pending.
- 256 MiB zRAM and 4 GiB fully allocated SD swap. Earlier running systems logged
  activation. Physical RAM remains 512 MiB; no repeated speed-up benchmark is
  claimed for the clean images.
- Native Desktop Flasher 4.0, DE/EN, GA7/GA11, full write/readback verification,
  APK/OBB queueing and installation-log reading through the card reader. It uses
  its own Win32/X11 window. Thirty file-based tests pass; the owner reports that
  the Linux flasher works without specifying which operation. Windows builds,
  static checks and PE headers pass; the first binary release is Linux only.
  Windows remains deferred.
- Exact K11 source/configuration and Android 11 pinned manifest/source changes
  collected for a separate public development/release project. The pinned
  manifest contains 863 upstream projects; 13 have collected local source edits.

The earlier browser-based tool remains in this repository as a historical
version. The new native tool and ROM files have not been uploaded here.

## Publication boundaries

The existing local launcher APK includes embedded game-logo resources. Those
resources are excluded from the planned public source selection, and the local
image is not silently changed. The native flasher contains no Android image or
game data. Source/redistribution evidence for retained Mali/audio/module/boot
components remains incomplete. A source-built kernel does not resolve those
separate inputs.

No current ROM release is declared public or fully accepted. The owner requested
that this repository continue as the personal development record and that a
separate **R36S-Clone-GA36-Android-Releases** repository be prepared for sources,
community work and releases. That new repository is not yet created/published.
USB-to-PC/OTG and the second SD slot remain deferred. TWRP remains optional.

See [machine-readable status](latest-status.json) and the
[release readiness record](../../release-readiness.md).
