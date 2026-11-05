# Scope, provenance and licensing

This package contains selected project-generated source files, board-specific configuration and documentation. It does not contain the proprietary stock firmware or Android donor binaries, game content, full kernel sources, extracted disassembly listings, third-party toolchains or copied upstream source trees.

The pointer helper and backlight implementation were written during this project with AI assistance. The backlight source reproduces small Android HAL ABI structure declarations; the interface reference is AOSP Android 5.0.2_r1, whose source carries its own license notices:

- https://android.googlesource.com/platform/hardware/libhardware/+/android-5.0.2_r1/include/hardware/hardware.h
- https://android.googlesource.com/platform/hardware/libhardware/+/android-5.0.2_r1/include/hardware/lights.h

The keylayout/IDC files are GA36-specific mappings and settings. Power documentation describes observations and a change specific to the exact local vendor kernel, without bundling that kernel.

No general project LICENSE has been selected by the owner yet. Do not describe this as an MIT-, GPL- or Apache-licensed complete ROM. The absence of a project license also means readers should not assume an unrestricted reuse grant. Public source availability and permission to redistribute someone else's firmware are separate matters.

This file records the publication scope and provenance; it is not a guarantee of worldwide legal clearance. A complete flashable image is excluded until its individual components and redistribution terms have been reviewed.
