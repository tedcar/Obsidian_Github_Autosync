# Obsidian Github Autosync

A single tiny Windows **service** that keeps any Obsidian vault (or arbitrary folder) automatically committed & pushed to a private GitHub repository every hour.

---

## Features
* 💾 **One 4 MB exe – no installer, no Python, no Git config needed**
* 🛡️ Stores your PAT securely in Windows Credential Manager (no plain-text)
* 🔄 Runs as a background *user service* – starts with Windows, no UI
* 🕒 Change-detector & hourly scheduler (interval configurable)
* 📜 Rotating log (`%APPDATA%/obsidian-auto-sync/sync.log` ≤ 1 MB)
* ↩️ Auto-pull with rebase, retry-on-network-error, conflict pause

## Quick-Start
1. **Download** `obsidian-auto-sync.exe` (see Releases) and place it inside your vault.
2. Open *PowerShell* **as Administrator** and run:
   ```powershell
   .\obsidian-auto-sync.exe --init
   ```
3. Follow the wizard prompts:
   • Vault path (press Enter for current dir)  
   • Sync interval in minutes (60 default)  
   • GitHub repo HTTPS URL (`https://github.com/<you>/<repo>.git`)  
   • GitHub username  
   • **PAT** (token with *repo* scope)
4. Approve the UAC dialog.  
   The service installs and starts as **ObsidianAutoSync**.

You're done – edits will sync on the next interval.

## Managing the service
```powershell
# view status
Get-Service ObsidianAutoSync

# restart after changing config.yaml
Stop-Service ObsidianAutoSync; Start-Service ObsidianAutoSync
```

## Uninstall
```powershell
Stop-Service ObsidianAutoSync
sc delete ObsidianAutoSync
Remove-Item "$Env:APPDATA\obsidian-auto-sync" -Recurse
```

## Building from source
```bash
go build -ldflags "-s -w" -o obsidian-auto-sync.exe ./auto-sync
```

---
MIT © 2025 tedcar 
