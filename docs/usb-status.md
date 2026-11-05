# USB status

Goals: PC file transfer through the OTG port, plus USB hub/flash-drive use.

V18 observations:

- Android init skipped `/init.ga36.usb.rc` because it was group-writable (0664).
- The gadget remained disabled and did not receive its expected functions.
- A host root hub existed, but no attached hub/flash drive enumerated.
- vold listed the external USB volume as NoMedia.

V19 secures init configuration files to0644. Physical USB operation after that correction remains unconfirmed. The GPIO/VBUS role detection has not been forced to host mode.

Selected board configuration observations:

```ini
usb_port_type = 2
usb_detect_type = 1
usb_id_gpio = port:PH08<0><1><default><default>
usb_det_vbus_gpio = "axp_ctrl"
usb_drv_vbus_gpio = port:power4<1><0><default><0>
usb_host_init_state = 0
```

These settings were read from this device's board configuration; they do not establish that a particular adapter supplies power or selects the correct role. Matching schematics, sources or known-working GA36 OTG configurations would help.

Only the external USB host controller is intended for removable storage management. The internal card is not exported as a USB mass-storage device. No real USB stick has been formatted as part of this publication preparation.
