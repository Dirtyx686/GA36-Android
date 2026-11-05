# GA36 Android research

Experimental Android 5.0 bring-up on an R36S-style handheld with the **GA36-MB V1.2-20260416** board.

This is a **source and documentation repository**, not a flashable firmware release. No stock firmware, Android system image, game ROM, SD-card image, private key or raw device log is included.

## Device and current status

Runtime analysis identifies an Allwinner A33 platform, 512 MB RAM and a vendor Linux 3.4.39 kernel. The photographed chips have Rockchip markings; Android donor model strings alone are not used as hardware identification. The donor system identifies as Android 5.0 / API 21.

Confirmed by the device owner:

- Android starts from a separate test SD card.
- Installed apps launch and audio works.
- Gamepad menu navigation works.
- Right stick moves a mouse pointer; R1 clicks.
- Display off and wake work.
- Power-button operation works after the V19 kernel correction.

Still pending: PC USB file transfer, USB hub/flash-drive operation and visual confirmation of launcher label layout. This has been tested on one board; compatibility with other revisions or A33 devices is not established.

## Repository contents

- `src/ga36-pointer.c`: ARM32 evdev-to-uinput mouse helper, with self-test.
- `src/lights-ga36.c`: experimental Android 5 ARM32 backlight HAL with delayed resume brightness restoration.
- `src/test-kernel-branches.c`: ARM execution test for the Power-event suppression branch change; it does not access a real device.
- `config/`: GA36-specific input mappings and device policy.
- `docs/power-button.md`: exact-build Power-button findings.
- `docs/usb-status.md`: confirmed USB observations and remaining questions.
- `docs/discussion.md`: prepared GitHub Discussion text.
- `tools/check-publication.py`: checks the exact publication file set and content hashes against the reviewed manifest.

These files are components of the experimental port, not a complete Android build tree or installer. The local boot helper, disk-writing tools and generated firmware are deliberately excluded.

## Sources and attribution

The work was developed with AI assistance and tested on the owner's device. Android interface definitions follow AOSP; runtime observations and board-specific mappings come from analysis of this device. See `NOTICE.md` for provenance and the current licensing status. No license is granted for third-party firmware through this repository.

## Related post

[Owner's Reddit post](https://www.reddit.com/r/R36S/comments/1xomg76/comment/19s6oe904jeyt/)
