# GA36 Android — my R36S clone development

Experimental Android development on my **GA36-MB V1.2-20260416** board, a
**R36S-style clone** identified at runtime as **Allwinner A33 with 512 MiB RAM**.
This is a record of my own device and tests. Compatibility with other R36S
boards, clones or revisions is not established.

I develop the project with AI assistance and test it on my handheld. Device
feedback, saved logs and host/build checks are distinguished below.

## Current status — 11 November 2026

| Component | Current position |
| --- | --- |
| GA7 | **Android 7 Beta V1**, P39 baseline; the existing device setup runs and sound is owner-confirmed. A separate clean-image candidate is prepared locally. |
| GA11 | **Android 11 Alpha V1**, P39 baseline; framework and launcher run on the handheld. The latest owner feedback confirms sound, GTA trailer playback and that the previously reported GTA exit crash is no longer observed. Exact tested build/card was not specified. |
| Kernel | Source-built **K11, Linux 3.4.39 #82**. Source, configuration and matching source hashes have been collected locally. Retained vendor modules still need separate source/provenance work. |
| Launcher | Device P39 uses versionCode101. Clean candidates use the versionCode102 DE/EN derivative, defaulting to English. Themes, console navigation and 2/4/6 game layouts are implemented; fresh-image device acceptance remains open. |
| Memory | **256 MiB zRAM + 4 GiB SD swap**. Activation was logged on earlier running setups. Physical RAM is still 512 MiB; no repeated performance benchmark establishes a general speed-up. |
| Clean images | Separate GA7 and GA11 images prepared locally for **32 GB or larger cards**, with first-boot growth. Decompression, hashes and filesystem checks pass; another-card first boot, expansion and app/data import still need physical acceptance. |
| Desktop Flasher | New native **4.0 Windows/Linux DE/EN** program prepared locally, supporting GA7/GA11 and card-reader APK/OBB import. Thirty file-based tests pass; the owner reports Linux operation. The initial release will be Linux only; a Windows binary release is deferred. |
| Deferred | USB-to-PC/OTG and the second SD slot. TWRP remains optional and paused. |

**No ready-to-flash Android image is published in this repository.** Current
images are release candidates, not a compatibility or stability guarantee.
Intermittent wake behavior, stability under memory pressure and complete
per-game control regression are not signed off by the latest general feedback.

[Current P39/K11 report](docs/updates/2026-11-11/p39-k11-status.md) ·
[Current status JSON](docs/updates/2026-11-11/latest-status.json) ·
[Release readiness](docs/release-readiness.md)

## Personal development here; public releases separately

This repository remains my development journal for the exact board above.
A separate **R36S-Clone-GA36-Android-Releases** project is being prepared for
community source work, build instructions and eventual GA7/GA11 releases.
It has not yet been created or published. Source and release packages are
being collected locally; no ROM upload is part of this status update.

The local clean images contain no user accounts or game APK/OBB packages.
Their current launcher still embeds some game logos; these resources and the
retained vendor/boot components need separate review before a public image.

## Native Desktop Flasher 4.0

The new flasher uses **its own Win32/X11 desktop window**, without a browser,
WebView or local web server. It supports separate GA7/GA11 images, SHA-256
checks, card writing with complete readback, APK/OBB queueing and INSTALL.TXT
reading. File transfer uses a microSD card reader because console USB-to-PC
connectivity remains unresolved. No ROM is included in the flasher package.

The native package remains local while release preparation continues. The first
planned binary release is **Linux only**. Windows is deferred. The owner reports
that the Linux flasher works; this does not certify every physical-card action.

[Earlier browser-based GUI tool](tools/flasher/README.md) — historical version,
with simulated-card screenshots and no Android image. Its earlier Windows
feedback does not validate the new native 4.0 write path.

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
sources are being assembled for the separate source/release project.

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

[Related Reddit post](https://www.reddit.com/r/R36S/comments/1xomg76/comment/19s6oe904jeyt/)
