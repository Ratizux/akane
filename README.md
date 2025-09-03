# AKANE

This fork is under early development.

# Modifications

- The `gofer` process no longer uses `chroot()` on the target directory
- Fall back to assuming `/proc/sys/vm/mmap_min_addr` = `32768` if retrieval fails
- On `aarch64`, assume ASID bits = 8 if retrieval fails
- Add support for kernels without `PTRACE_SYSEMU` (Linux < 5.3 on aarch64)
- Add support for hosts with 39-bit virtual memory
- Disable `mseal` on the `systrap` platform (some buggy syscall filters, especially on Android, may terminate the process if an unsupported syscall is used)
- Experimental fakehostfs
    - Emulates file UID/GID and other behaviors on Android
    - Not fully tested
    - mmap() not implemented
