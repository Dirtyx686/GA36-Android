# Clean dual-boot replacement — 12 November 2026

Exact board: **GA36-MB V1.2-20260416**, R36S-style clone, runtime Allwinner A33,
512 MiB physical RAM. The working P39 device setup shares K11 Linux 3.4.39 #82.

The owner reports that the earlier separate GA7/GA11 clean images do not boot,
while the existing dual-boot setup does. Host hashes and filesystem checks did
not establish physical boot compatibility. Those standalone packages are now
marked as failed device candidates, and the updated local flasher rejects them.
The precise boot failure cause has not been proven.

The replacement keeps the original Allwinner bootloader prefix and raw offsets
rather than substituting a DOS MBR. Its native pre-Android menu selects GA7 or
GA11 with the same kernel. D-pad selects, A/START confirms; the five-second
automatic GA7 start pauses when a direction is pressed.

The new clean English candidate targets **64 GB or larger cards**, at least
63,909,658,624 bytes. Both userdata filesystems are fresh, initially 24 GiB each.
The retained 256 MiB zRAM/4 GiB shared SD swap configuration remains. Larger
cards grow the raw outer pool before mounting and split the userdata budget
equally after system/vendor/swap space and a reserve. This growth needs a
physical larger-card test.

The GA36 launcher defaults to English. The new resource variant omits bundled
game artwork; the private original launcher remains unchanged. The GA7 system
SkyRunner package was also removed. No user accounts, installed games, OBBs,
saves or old startup logs are copied into the fresh data filesystems. The GA11
userdata root label is system_data_root_file, matching its platform policy.

The native Linux flasher 4.1-dual candidate validates the original boot prefix
and leaves it unchanged after complete write/readback. A separate private
Windows CMD candidate uses a native console and ordinary Windows elevation,
with no browser or PowerShell. Windows execution and real writing are untested.
The planned initial public binary remains Linux only.

The old APK/OBB queue requires a GA36_APPS FAT area. The preserved dual-boot
layout has no such area; importing into this new image through that menu is
not supported yet. USB-to-PC/OTG and the second SD slot remain deferred.

Full image decompression, raw SHA-256, component-region hashes and filesystem
checks pass on the host. The compressed artifact is about 950 MiB. Physical
acceptance remains pending. No new image, executable, private key or user log is uploaded by this
status update. The separate release repository still awaits owner approval.

[Machine-readable status](latest-status.json)
