# Power-button investigation — V19

The exact tested boot image before the change has SHA256:

`9795fc84c711de4e67983c84ff8774cd9db8e372d0c9fae65a12f2b40800f5ef`

Analysis of that kernel found a flag set by a write to `/sys/power/state`. Two conditional branches in the AXP Power-key interrupt path skip native key-down and key-up processing when that flag is1. Physical key activity still appeared in kernel messages, but evdev stopped delivering it after the first screen sleep operation.

For this exact build only:

| Virtual address | Kernel byte offset | Original ARM instruction | Tested replacement |
| --- | --- | --- | --- |
| `0xc03f685c` | `4122716` | `0x0a000019` | `0xe320f000` (NOP) |
| `0xc03f68f4` | `4122868` | `0x0a000019` | `0xe320f000` (NOP) |

The original boot header's SHA1 ID was recalculated after changing the two instructions. The ramdisk and image layout were retained. The patched boot-image SHA256 is:

`c8e0a19a3defdc7a5a74f782ec95f6eb3e58a5a28f68b53b5320f1d44dde4cbb`

These addresses and hashes document the result; they are not a generic patch recipe for other firmware. No boot image or disk-writing script is included.

The vendor `micro_gamepad/power_key` sysfs interface reports1 after a short release (up to500ms), and its timer sets2 after1000ms. An earlier synthetic Power-key workaround misread1 as long and was removed in V19.

Android's intended policy is short press for screen sleep/wake and long press for the global actions menu. The owner confirmed that the Power button works after V19. This is not an endurance or all-power-state validation.

The included ARM branch test checks the original suppression and the NOP replacement under an ARM emulator. It is not a full kernel boot test.
