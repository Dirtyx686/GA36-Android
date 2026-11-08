# Scope, provenance and licensing

This package contains selected project-generated source files, board-specific configuration and documentation. It does not contain the proprietary stock firmware or Android donor binaries, game content, full kernel sources, extracted disassembly listings, third-party toolchains or copied upstream source trees.

The pointer helper and backlight implementation were written during this project with AI assistance. The backlight source reproduces small Android HAL ABI structure declarations; the interface reference is AOSP Android 5.0.2_r1, whose source carries its own license notices:

- https://android.googlesource.com/platform/hardware/libhardware/+/android-5.0.2_r1/include/hardware/hardware.h
- https://android.googlesource.com/platform/hardware/libhardware/+/android-5.0.2_r1/include/hardware/lights.h

The keylayout/IDC files are GA36-specific mappings and settings. Power documentation describes observations and a change specific to the exact local vendor kernel, without bundling that kernel.

The owner has selected MIT for original project-owned code. The scoped license package remains separate from this documentation update; existing file-level notices still apply, and this announcement is not a blanket relicensing of mixed-provenance files. Do not describe this as an MIT-, GPL- or Apache-licensed complete ROM. Readers should not infer an unrestricted reuse grant for every file or third-party component. Public source availability and permission to redistribute someone else's firmware are separate matters.

This file records the publication scope and provenance; it is not a guarantee of worldwide legal clearance. A complete flashable image is excluded until its individual components and redistribution terms have been reviewed.
