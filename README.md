# AKANE

This fork is under early development.

# Modifications

- Now `do` will not apply namespace for `runsc-sandbox`
- `gofer` will not `chroot()` into target directory
- Guess `/proc/sys/vm/mmap_min_addr` as `32768` if failed to retrieve
- On `aarch64` platforms, will guess ASID bits as `8` if failed to retrieve
- Support kernel without `PTRACE_SYSEMU` (Linux < 5.3 for `aarch64` platforms)
- Support hosts with 39-bits virtual memory
- Disable `mseal` of `systrap` platfrom, for some malfunctioning syscall filters (especially on Android) will kill the process if unsupported syscall is used
- Experimental `fakehostfs` that emulates file UID/GID and more on Android. Not fully tested, `mmap()` is not implemented, expect bugs
