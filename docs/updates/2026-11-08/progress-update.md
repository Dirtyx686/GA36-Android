# GA36 Android 7 — 8 November 2026

I develop this port with AI assistance and test it on my own GA36-MB V1.2-20260416 handheld. Runtime identification is Allwinner A33 / 512 MB RAM despite the photographed Rockchip markings. Compatibility with other boards is not established.

## Latest local build: V65

V65 has been written to the test SD card and verified by direct readback. It retains the V38 direct-start baseline and patched vendor Linux 3.4.39 #522. The boot region was not rewritten for V65. **Physical V65 acceptance is pending.** V56 cold start was confirmed previously; earlier tests established apps, sound, gamepad navigation, right-stick mouse/R1 click, display wake and Power operation.

The local launcher is version 1.9.0-v65, with eight themes, two/four/six-tile layouts and fullscreen mode. Holding X+Y together for 1.5 seconds requests background-app cleanup. Host tests cover hold/release/repeat behavior, one-shot dispatch, foreground and input-method protection, permission failures and cancellation. These tests use Android test doubles; they do not prove physical controller/service delivery. The cleanup does not delete app data.

V65 removes overlapping automatic installer/compilation work and obsolete maintenance components while preserving runtime caches. First-boot package maintenance remains pending on the device. No new boot-time or performance improvement has been measured. Private application inventories, application resources and raw logs are excluded from this publication.

## Still open

- Test V65 cold start, first-boot maintenance and launcher behavior on hardware.
- Verify global X+Y, short-key/FN behavior, sound, display wake and Power menu.
- Recheck package-installer permission-dialog responsiveness and measure later boots.
- PC USB and OTG hub/flash-drive support remain unresolved; investigation is deferred.
- Secondary SD has no confirmed physical card response.
- V39 custom kernel: built locally; five modules pass static compatibility checks. Analog receive path and hardware tests remain open. It has not been installed on the SD card.
- TWRP development remains paused indefinitely.

## Publication

The published implementation remains mainly the earlier Android 5 components. This update adds documentation and curated status, not the full Android 7 implementation or a firmware release. Older dated reports remain historical evidence.

No stock ROM, Android/card image, game content, foreign binary or private log is included. A public image still requires component redistribution/source compliance review, a newly prepared clean user-data image, and testing on another card. See [release readiness](../../release-readiness.md).
