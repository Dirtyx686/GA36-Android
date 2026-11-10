# Android 11 P23 / K09 update — 10 November 2026

I am developing this GA36-MB V1.2-20260416 port with AI assistance and testing it on one A33/512 MiB handheld. Android 11 has now booted on the device. This supersedes the earlier “Android 11 preparation only / locked” status as a current project description; those dated records remain historical.

## Latest installed state

Android 11 P23 was installed with complete direct readback of the changed ranges and safely ejected. It uses the source-built K09 Linux 3.4.39 kernel, build #80. The installer preserved Android 7 D04 #72 boot/system contents and user data. This report does not perform or authorise another card write.

P21/P22 records confirm K09 and Android 11 framework/GUI startup. My latest P23 feedback reports that Android 11 runs, but the mouse still has no reaction. I subsequently confirmed Settings navigation works. Fresh detailed P23 runtime logs are not yet available; the latest detailed findings are from P22. These evidence levels are deliberately separate.

The four new user-supplied photographs show Android 11 branding, a startup screen, the GA36 launcher and Settings. They are not independently version-labelled P23, do not prove every subsystem, and do not establish gameplay just because a game tile is visible.

## P23 changes and results

- The mouse helper was changed to finish short writes accepted by the legacy Linux 3.4 uinput interface, retry interrupted writes and retain complete motion/synchronisation events. ARM regression and negative-control tests passed locally. The fix is installed, but the latest device report still says the mouse is unresponsive. It is not considered solved.
- Settings no longer forces focus onto the outer scroll container; selectable descendants are allowed. The compiled APK/code was checked, and the latest user feedback now confirms Settings navigation works. Every Settings function has not been tested.
- P22 keyguard/SystemUI changes remain included. The P22 capture had no new SystemUI ANR, but still showed DocumentsUI and external-storage ANRs. P23 is not claimed to have fixed those hangs without fresh logs.
- Audio, gameplay, repeated Power/wake and full button/stick regression on P23 remain to be specifically recorded. General “runs for now” feedback is not acceptance of every function.

## Memory and observed boot timing

The last detailed P22 capture confirms 256 MiB zRAM at priority 100 and a 4 GiB SD swap file at priority 10. Both reported zero usage at that snapshot. Physical RAM remains 512 MiB. These are configured swap resources, not additional physical RAM or a measured performance gain.

One recorded P22 framework-startup time was 130.112 seconds, compared with 160.071 seconds in the referenced P21 capture. This is a comparison of individual recorded runs, not a repeated benchmark, an end-to-end boot guarantee or a P23 timing result. Device I/O, memory pressure and storage-service stalls still need measurement.

## OS selection and next work

Android 7 remains preserved with D04 #72; Android 11 P23 is the active experimental test. OS switching uses the existing restart-based mechanism. USB to PC, the second SD slot and direct OS startup without an extra restart remain deferred. TWRP remains optional and paused.

P24 has been requested for further mouse diagnosis, measured performance work, boot-menu UI/test labels and documented bug fixes. No P24 code change, build or installation has been made in the reviewed state. Development is paused at my request; this GitHub task updates documentation and supplied photos only.

The next device investigation should capture P23 input/service logs and verify pointer lifetime, virtual input registration, axis values/R1 events and motion-write errors before claiming a functional mouse.

## Publication and licensing boundary

Only this curated report, a matching status JSON, README changes and the selected test photographs are published. No raw runtime logs, card images, Android/kernel binaries, private paths or keys, game APKs/data or standalone game artwork are included. The launcher photograph visibly contains an installed-game tile and is included at my explicit request; the game itself is not distributed.

Full kernel/source/build packages remain local for now. The local BusyBox replacement package and earlier emulator tests are not proof it is installed in P23. Rights/corresponding-source questions for retained Mali, audio, modules and boot components remain open; retaining components does not clear a complete image release.

## Local evidence reviewed

- Current-target/continuation record, P23 install/readback report and latest user feedback.
- P23 source/integration and compiled-Settings validation summaries.
- P22 device findings for boot timing, uinput failures, ANRs and active swap.
- User-supplied Android 11 photographs and subsequent Settings confirmation in this chat.

Raw records and their private paths are not published. Earlier “hardware pending” entries are not silently treated as later acceptance: current user feedback and P22 logged findings are identified explicitly above.
