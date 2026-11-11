# GA7 / GA11 release preparation

Status: **12 November 2026**. Exact board: **GA36-MB V1.2-20260416**,
R36S-style clone, runtime Allwinner A33 and 512 MiB RAM. No ROM is public.

The owner reports that the earlier separate clean images do not boot. The
working dual-boot setup is the basis for the new English clean 64-GB+ candidate.
It keeps K11 #82, the original bootloader prefix and fixed raw storage offsets,
with fresh userdata and without installed games, accounts or saves. A Linux
4.1-dual flasher candidate handles this layout without rewriting its header.

Remaining:

1. Full image hash/readback checks passed. Now test GA7 and GA11 on a new card:
   launcher, sound, buttons, Settings/file manager, power menu and repeated wake.
2. Test equal userdata expansion on a larger card. zRAM is 256 MiB, shared SD
   swap is 4 GiB, physical RAM remains 512 MiB; no general speed-up is claimed.
3. Add APK/OBB transfer for the preserved dual layout. The old GA36_APPS queue
   cannot target the new layout. USB and slot 2 remain deferred.
4. Keep the initial public binary Linux only; the private Windows CMD candidate
   still needs actual Windows execution and physical-card writing tests.
5. Complete retained boot/Mali/audio/module source/provenance and notice work.
   A source-built kernel does not establish redistribution rights for blobs.
6. Present the source repository and asset list for owner approval before a
   new repository or public release. No ROM upload is authorized by this update.

Collected locally: complete K11 source/configuration, 863 pinned Android 11
upstream projects and 13 projects with local source edits, and launcher source
without bundled game artwork. New clean-image physical acceptance is pending.

TWRP remains optional and paused. MIT applies to original project code only.
See [NOTICE](../NOTICE.md) and the [current report](updates/2026-11-12/clean-dual-boot-status.md).
