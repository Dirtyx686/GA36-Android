Android 5.0 is now running on my R36S-style handheld with a **GA36-MB V1.2-20260416** board. I am documenting the progress and looking for information about USB OTG on this board.

## Hardware and software

- Board marking: GA36-MB V1.2-20260416
- Runtime platform: Allwinner A33
- RAM: 512 MB
- Vendor kernel: Linux 3.4.39, build #522
- Android: 5.0, ARM32 / API 21
- Android donor model string: TU-7981B

The photographed chip markings say Rockchip RK3326 and RK817, while the running platform identifies as Allwinner A33. The donor's model and processor strings alone should not be treated as hardware identification.

Android boots persistently from a separate test SD card. The original card has been kept separate.

## Confirmed working on the device

- Display and Android launcher
- Launching installed apps
- Audio
- Gamepad navigation
- Right stick as mouse and R1 as left click
- Display off and wake
- Power button operation after the latest kernel fix

## Power-button findings

The actual vendor kernel contains two checks in its Power-key interrupt path that suppress native key-down and key-up events once a flag set by a write to `/sys/power/state` becomes 1. Replacing those two conditional branches with NOP instructions restored Power-button operation in the device test.

For this exact kernel binary, the affected virtual addresses are `0xc03f685c` and `0xc03f68f4`. These addresses are specific to this build and are not instructions for patching another firmware.

A separate vendor interface, `/sys/devices/platform/micro_gamepad/power_key`, reports `1` for a short press and `2` for a long press. An earlier userspace workaround misinterpreted that interface and has been removed; Android now receives the native Power events.

## USB OTG: help wanted

The goals are:

- Connect the handheld to a PC for file transfer.
- Connect a USB hub and flash drive to the handheld.

The latest update corrects an Android init configuration permission problem: `/init.ga36.usb.rc` was group-writable and Android init skipped it. The file permissions are now corrected, and legacy `android_usb` MTP/PTP and ADB configuration handlers are included. Physical USB testing after this correction is still pending.

Earlier logs showed the USB host root hub, but no attached hub or flash drive enumerated. The external USB volume is configured through vold, but this does not establish that hardware detection or VBUS power is working.

The stock board configuration includes:

```ini
[usbc0]
usb_used = 1
usb_port_type = 2
usb_detect_type = 1
usb_id_gpio = port:PH08<0><1><default><default>
usb_det_vbus_gpio = "axp_ctrl"
usb_drv_vbus_gpio = port:power4<1><0><default><0>
usb_host_init_state = 0
```

The kernel's automatic cable-ID/VBUS role detection has been retained. No forced host-power setting is being used.

Does anyone have **GA36-MB-specific USB ID/VBUS wiring information, matching kernel sources, or a known-working OTG configuration**?

## Other remaining issue

Some app names are clipped in the donor launcher. A launcher-layout adjustment still needs visual confirmation.

This is an experimental port, not a ready-to-flash release. No complete installation procedure or release image is provided in this post.

## Related discussion

My Reddit post: https://www.reddit.com/r/R36S/comments/1xomg76/comment/19s6oe904jeyt/
