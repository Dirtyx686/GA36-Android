# Clean dual-boot device test — 12 November 2026

Target: **GA36-MB V1.2-20260416** R36S-style clone, runtime Allwinner A33,
512 MiB RAM. Candidate: P39 clean English GA7 + GA11, K11 #82 unchanged.
This test plan describes the locally held package; it is not a ROM download.

## Prepare a new card

Use a new 64-GB or larger microSD with at least **63,909,658,624 bytes**.
Keep the working test card intact. Use the local native Linux flasher
**4.1-dual**, select **Dual Boot**, select the intended new card, and confirm
its deletion. Every existing file on that target is erased. Wait for writing
and complete readback verification before safely removing it. Cancel any PC
prompt to format its partitions: the original Allwinner layout is intentional.

Earlier separate GA7/GA11 image packages failed the owner's boot test and are
marked not_for_flashing. Use the new dual image described in the
[reference manifest](dualboot-image.json).

## Boot selection

The D-pad selects GA7 or GA11. **A or START** confirms. With no input, GA7
starts after five seconds; pressing a direction pauses the timer. **B** returns
the selection to GA7. Both boot paths use the shared kernel. First startup on
a larger card can take longer while its storage grows; completion at the
launcher is not yet established for this fresh image.

## Test each system

Test GA7 first, shut down through its menu, then start and select GA11.
Record each row for both systems; the table below is a checklist, not a pass.

| Check | What to observe |
| --- | --- |
| First boot | Boot menu, matching English logo and arrival at GA36 launcher |
| Clean state | English default, no old games, accounts or saves |
| Launcher | D-pad navigation, A/B, access to menu, Settings and file manager |
| Settings and files | Select and open entries; no unexpected return or hang |
| Sound | Built-in speaker plays a test sound and volume can be changed |
| Sleep/wake | Three short-Power screen-off/wake cycles without a stuck black screen |
| Power menu | Long Power shows the custom menu; shutdown/restart works |
| Storage | Both systems initially have separate 24-GiB userdata; inspect capacity |
| Larger card | After each OS has booted, verify equal growth and retained files after reboot |
| Memory | Return the card for logs proving current zRAM and SD-swap activation |

A successful first boot alone does not sign off larger-card growth, repeated
wake or stability under memory pressure. No current-game test is required for
this initial clean-system check; the image intentionally includes no games.

## Report back

Note the selected OS, whether it reached the launcher, any visible error,
and how long it ran before a freeze or reset. After testing, shut down through
the menu and return the system card to the PC reader for the new logs.

The old GA36_APPS APK/OBB queue cannot target this preserved dual layout.
USB-to-PC/OTG and the second SD slot remain deferred.
