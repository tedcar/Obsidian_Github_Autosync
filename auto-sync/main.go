package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/kardianos/service"
    "github.com/zalando/go-keyring"
    "gopkg.in/natefinch/lumberjack.v2"
)

const keyringService = "obsidian-auto-sync"
const keyringUser = "github-pat"

func main() {
    // setup rotating log file early
    if dir, err := defaultConfigDir(); err == nil {
        log.SetOutput(&lumberjack.Logger{
            Filename: filepath.Join(dir, "sync.log"),
            MaxSize: 1,   // megabytes
            MaxBackups: 3,
            MaxAge: 7,    // days
            Compress: true,
        })
    }

    initFlag := flag.Bool("init", false, "run interactive setup and install service")
    runOnce := flag.Bool("once", false, "run one sync cycle and exit (for testing)")
    flag.Parse()

    if *initFlag {
        if err := interactiveSetup(); err != nil {
            log.Fatal(err)
        }
        fmt.Println("Setup completed. Service installed and started.")
        return
    }

    cfg, err := loadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    if *runOnce {
        if err := syncVault(cfg); err != nil {
            log.Fatalf("sync failed: %v", err)
        }
        return
    }

    // Service config
    svcConfig := &service.Config{
        Name:        "ObsidianAutoSync",
        DisplayName: "Obsidian Vault Auto-Sync",
        Description: "Background service that syncs an Obsidian vault to GitHub.",
        Option: service.KeyValue{
            "UserService": "true",
        },
    }

    prg := &program{cfg: cfg}
    s, err := service.New(prg, svcConfig)
    if err != nil {
        log.Fatal(err)
    }

    if err := s.Run(); err != nil {
        log.Fatal(err)
    }
}

func interactiveSetup() error {
    // Prompt user (basic stdin prompts)
    var vaultPath string
    fmt.Print("Vault absolute path (leave blank for current directory): ")
    fmt.Scanln(&vaultPath)
    if vaultPath == "" {
        cwd, _ := os.Getwd()
        vaultPath = cwd
    }
    vp, err := filepath.Abs(vaultPath)
    if err != nil {
        return err
    }

    // Interval
    var interval int
    fmt.Print("Sync interval minutes (default 60): ")
    fmt.Scanln(&interval)
    if interval <= 0 {
        interval = 60
    }

    // Repo URL
    var repo string
    fmt.Print("GitHub repository HTTPS URL (e.g., https://github.com/user/repo.git): ")
    fmt.Scanln(&repo)

    // Username
    var user string
    if repo != "" {
        // try infer user from URL if pattern matches https://github.com/<user>/<repo>.git
        parts := strings.Split(repo, "/")
        if len(parts) >= 5 {
            user = parts[3]
        }
    }
    fmt.Printf("GitHub username [%s]: ", user)
    var tmp string
    fmt.Scanln(&tmp)
    if tmp != "" {
        user = tmp
    }

    // PAT
    fmt.Print("Personal Access Token (input hidden after enter): ")
    var pat string
    fmt.Scanln(&pat)

    if pat == "" {
        return fmt.Errorf("PAT cannot be empty")
    }

    if err := keyring.Set(keyringService, keyringUser, pat); err != nil {
        return fmt.Errorf("keyring: %w", err)
    }

    cfg := &Config{
        VaultPath:      vp,
        IntervalMinutes: interval,
        RepoURL:        repo,
        Username:       user,
        RemoteName:     "origin",
    }
    if err := saveConfig(cfg); err != nil {
        return err
    }

    // Configure remote URL with embedded PAT for non-interactive auth
    authRepoURL := repo
    if !strings.Contains(repo, "@") {
        authSegment := fmt.Sprintf("%s:%s@", user, pat)
        authRepoURL = strings.Replace(repo, "https://", "https://"+authSegment, 1)
    }
    // add or set remote
    if err := runGit(vp, "remote", "add", cfg.RemoteName, authRepoURL); err != nil {
        // if remote exists, try set-url
        _ = runGit(vp, "remote", "set-url", cfg.RemoteName, authRepoURL)
    }

    // install service using kardianos
    svcConfig := &service.Config{
        Name:        "ObsidianAutoSync",
        DisplayName: "Obsidian Vault Auto-Sync",
        Description: "Background service that syncs an Obsidian vault to GitHub.",
        Option: service.KeyValue{
            "UserService": "true",
        },
    }
    prg := &program{cfg: cfg}
    s, err := service.New(prg, svcConfig)
    if err != nil {
        return err
    }
    if err := s.Install(); err != nil {
        // best-effort detection of an "already installed" error message
        if !strings.Contains(err.Error(), "already") {
            return err
        }
    }
    if err := s.Start(); err != nil {
        return err
    }
    return nil
} 