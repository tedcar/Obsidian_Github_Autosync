package main

import (
    "context"
    "log"
    "time"

    "github.com/kardianos/service"
)

// program implements service interface

type program struct {
    cancel context.CancelFunc
    cfg    *Config
}

func (p *program) Start(s service.Service) error {
    // start should not block, so run in goroutine
    ctx, cancel := context.WithCancel(context.Background())
    p.cancel = cancel
    go p.run(ctx)
    return nil
}

func (p *program) run(ctx context.Context) {
    ticker := time.NewTicker(time.Duration(p.cfg.IntervalMinutes) * time.Minute)
    defer ticker.Stop()

    // initial sync right away
    if err := syncVault(p.cfg); err != nil {
        log.Printf("sync error: %v", err)
    }

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := syncVault(p.cfg); err != nil {
                log.Printf("sync error: %v", err)
            }
        }
    }
}

func (p *program) Stop(s service.Service) error {
    if p.cancel != nil {
        p.cancel()
    }
    return nil
} 