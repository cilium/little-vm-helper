# Solution for Issue #716

## 🛠️ Proposed Solution (by Aditya Waghamare)

### Analysis
The `little-vm-helper` tool uses Qemu to run VMs, currently hardcoding the `-append` kernel command-line flags in the `run` command. To support custom kernel arguments or replace the default `-append` string (e.g., to enable specific LSMs like `bpf`), we need to add a new command-line flag (e.g., `--append` or `--kernel-args`) to the `run` command structure, parse it, and use it when constructing the Qemu command invocation.

### Fix
Add a `--kernel-args` / `--append` flag to the `run` command configuration and update the Qemu runner to override or append to the default kernel arguments.

### Implementation
```go
// Add to command/run struct or options:
var kernelArgs string

// In flag definitions:
cmd.Flags().StringVar(&kernelArgs, "kernel-args", "", "Additional or replacement kernel command line arguments (passed to Qemu -append)")

// In Qemu execution logic:
appendArgs := "root=/dev/vda console=ttyS0 earlyprintk=ttyS0 panic=-1 lsm=lockdown,capability,landlock,yama,apparmor"
if kernelArgs != "" {
    appendArgs = kernelArgs
}
qemuArgs = append(qemuArgs, "-append", appendArgs)
```

### Testing
Verify by running `little-vm-helper run --kernel-args "root=/dev/vda console=ttyS0 earlyprintk=ttyS0 panic=-1 lsm=lockdown,capability,landlock,yama,apparmor,bpf"` and confirming Qemu receives the custom `-append` argument.

Signed-off-by: Aditya Waghamare <adityawaghamare7620@gmail.com>

---
*Submitted by Aditya Waghamare*
💰 **Payout Address (Base L2 / EVM):** `0xb61dBcdBc3407F71EaCb64D4CBFAcf9FFfe2415C`