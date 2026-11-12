# GA36 Android — my R36S clone development

Experimental Android development on my **GA36-MB V1.2-20260416** board, a
**R36S-style clone** identified at runtime as **Allwinner A33 with 512 MiB RAM**.
This is a record of my own device and tests. Compatibility with other R36S
boards, clones or revisions is not established.

I develop the project with AI assistance and test it on my handheld. Device
feedback, saved logs and host/build checks are distinguished below.

## Download P39 — EXPERIMENTAL

**64 GB or larger · Linux flasher · no games or user data**

| Download | File |
| --- | --- |
| **[Dual-boot ROM](https://github.com/Dirtyx686/R36S-Clone-GA36-Android-Releases/releases/download/dualboot-p39-experimental/GA36-DualBoot-GA7-GA11-Compact-EN-P39.ga36.gz)** | GA7 + GA11 compact package |
| **[IMAGE-INFOS.json](https://github.com/Dirtyx686/R36S-Clone-GA36-Android-Releases/releases/download/dualboot-p39-experimental/IMAGE-INFOS.json)** | Required beside the ROM |
| **[Linux desktop flasher](https://github.com/Dirtyx686/R36S-Clone-GA36-Android-Releases/releases/download/dualboot-p39-experimental/GA36-Flasher-Linux-4.2-Compact.tar.gz)** | Extract and run on Linux x86-64 |
| [Checksums](https://github.com/Dirtyx686/R36S-Clone-GA36-Android-Releases/releases/download/dualboot-p39-experimental/SHA256SUMS.txt) | SHA-256 download verification |
| [All release files and sources](https://github.com/Dirtyx686/R36S-Clone-GA36-Android-Releases/releases/tag/dualboot-p39-experimental) | Expand **Assets** on the release page |

Keep the ROM and `IMAGE-INFOS.json` in the same folder. Use GA36 Flasher 4.2:
the compact `.ga36.gz` file is not a raw image for generic flashers. Writing
overwrites the selected card. **GA11 audio is intermittent; sleep/wake has bugs.**

## Current status — 12 November 2026

| Component | Current position |
| --- | --- |
| GA7 | Android 7 Beta V1, P39. The owner reports the clean dual image boots; GA7 sound works. |
| GA11 | Android 11 Alpha V1, P39. The owner reports it boots; audio is intermittent. Earlier successful audio/trailer tests do not establish stable sound in the released setup. |
| Dual boot | Public experimental P39 release, English defaults, clean data without games, accounts or saves. Preserves the working boot prefix and offsets; 64-GB+ card required. |
| Kernel | Source-built K11, Linux 3.4.39 #82, shared by both systems. Complete source snapshot is attached to the release. |
| Launcher | GA36 launcher, English default, DE/EN strings, themes and 2/4/6-tile layouts. |
| Memory | 512 MiB physical RAM; 256 MiB zRAM and one shared 4-GiB SD swap file. This does not add physical RAM or prove a speed improvement. |
| Storage | Initially 24 GiB fresh userdata per system. Larger-card expansion still needs a physical test. |
| Linux flasher | Native 4.2 Compact, DE/EN. Host checks pass and the offline source rebuild is bit-identical to the delivered Linux binary. Compact restore and filesystem checks passed on a regular host test file. |
| Windows | Private test candidate; the initial public binary flasher is Linux only. |
| Open work | Sleep/wake bugs, intermittent GA11 audio, larger-card growth, sustained memory pressure and APK/OBB import for this dual layout. USB/OTG and SD slot 2 are deferred. |

The compact ROM is 754,064,053 bytes and requires at least 63,909,658,624 bytes
on the destination card. It writes and verifies about 8.67 GB each while retaining
the 64-GB layout. It requires a compatible GA36 flasher, not a generic raw writer.
All 14 uploaded release assets have matching local and GitHub SHA-256 digests.

## Personal development here; public releases separately

This repository remains my development journal. Current downloadable ROMs,
Linux flasher, K11 sources, GA11 integration, native boot and launcher source,
new device photographs and contribution notes are now collected in the
[R36S-Clone-GA36-Android-Releases project](https://github.com/Dirtyx686/R36S-Clone-GA36-Android-Releases).

This is my first project of this kind, developed with AI assistance and tests
on my own handheld. Everyone is welcome to contribute, fork and independently
further develop the project under the included component licenses. Original
project code is MIT-licensed; third-party terms remain in effect.

The release is an **experiment**, not a stable daily-use firmware. Historical
reports below remain dated records and do not override this current status.
The earlier browser-based flasher in this development repository is historical;
use the native Linux 4.2 Compact flasher linked above for the current ROM.

## My handheld and board

<img src="assets/photos/android11-settings.jpg" alt="GA36 handheld displaying Android 11 Settings" width="360">

*Owner-supplied device-test photograph from the earlier Android 11 work. It does
not identify or validate the current P39 clean image by itself.*

<details>
<summary>Device, board and earlier startup photographs</summary>

<img src="assets/photos/board-component-side.jpg" alt="GA36-MB V1.2-20260416 component side" width="360">
<img src="assets/photos/board-button-side.jpg" alt="Button side with card slots and USB ports" width="360">

Printed chip markings differ from the Allwinner runtime identification.

<img src="assets/photos/android7-device.jpg" alt="Earlier Android 7 device test" width="320">
<img src="assets/photos/android7-memory.jpg" alt="Earlier Android memory snapshot" width="320">
<img src="assets/photos/android11-boot-logo.jpg" alt="Earlier Android 11 boot branding" width="280">
<img src="assets/photos/android11-starting.jpg" alt="Earlier Android startup stage" width="280">
<img src="assets/photos/android11-launcher.jpg" alt="Earlier Android 11 launcher on the device" width="320">

The launcher photo includes an installed-game tile; no game package or data is
provided. Memory photos are snapshots, not current performance benchmarks.

</details>

## Sources and history

Published `src/` and `config/` files still include early Android 5 bring-up
code. They are not a complete source representation of the current P39/K11
integration. Current kernel, Android integration, launcher and native boot
sources are available in the separate source/release project linked above.

- [Clean dual-boot replacement — 12 November](docs/updates/2026-11-12/clean-dual-boot-status.md)
- [P39 / K11 — 11 November](docs/updates/2026-11-11/p39-k11-status.md)
- [P23 / K09 — 10 November](docs/updates/2026-11-10/p23-k09-status.md)
- [V69 / D03 preparation](docs/updates/2026-11-08/v69-d03-status.md)
- [Earlier source and validation review](docs/updates/2026-11-08/source-and-validation-update.md)
- [Earlier Android 7 bring-up](docs/updates/2026-11-06/progress-update.md)

Historical reports are retained as dated evidence; their pending flags do not
override the current report. The old menu illustration is also historical:
[earlier system-selection illustration](assets/illustrations/system-selection.jpg).

## Publication scope and attribution

[NOTICE.md](NOTICE.md) records provenance and licensing boundaries. MIT is
selected for original project-owned code only; third-party source, drivers,
firmware and photographs retain their own terms. No signing private key,
account data, raw device log or game package is added by this update.

`PUBLICATION-MANIFEST.json` pins the reviewed repository files. To check them:

```sh
python3 tools/check-publication.py
```

This checks the file set, hashes and selected content; it is not a signed
release or a legal/security clearance.

[Release post with device photos](https://www.reddit.com/r/R36S/comments/1xuisew/experiment_android_7_11_dualboot_release_for_the/)
