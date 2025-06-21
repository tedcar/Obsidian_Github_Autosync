package main

import (
    "bytes"
    "errors"
    "fmt"
    "os/exec"
    "time"
)

func runGit(dir string, args ...string) error {
    cmd := exec.Command("git", args...)
    cmd.Dir = dir
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("git %v failed: %v\n%s", args, err, out.String())
    }
    return nil
}

func syncVault(cfg *Config) error {
    vault := cfg.VaultPath
    // 1. git pull --rebase
    if err := runGit(vault, "pull", "--rebase", cfg.RemoteName, "main"); err != nil {
        // ignore non-fast-forward but log
        fmt.Println(err)
    }

    // 2. git add -A
    if err := runGit(vault, "add", "-A"); err != nil {
        return err
    }

    // 3. check if there is diff
    cmd := exec.Command("git", "diff", "--cached", "--quiet")
    cmd.Dir = vault
    err := cmd.Run()
    if err == nil {
        // exit code 0 => no changes
        return nil
    }
    var exitErr *exec.ExitError
    if !errors.As(err, &exitErr) {
        return fmt.Errorf("git diff error: %w", err)
    }
    // changes detected (exit code 1)

    // 4. commit
    msg := fmt.Sprintf("Auto-sync: %s", time.Now().Format("2006-01-02 15:04:05"))
    if err := runGit(vault, "commit", "-m", msg); err != nil {
        return err
    }

    // 5. push
    if err := runGit(vault, "push", cfg.RemoteName, "HEAD"); err != nil {
        return err
    }
    return nil
} 