// Package presence tracks workspace/editor state and pushes it to
// Discord via rich-go. Updates are debounced so Discord's rate limits
// aren't hit.
package presence

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/unstablebuild/rune-go-sdk/api/extensionapi"
	"github.com/unstablebuild/rune-go-sdk/api/workspaceapi"
)

const (
	// debounceDelay coalesces rapid bursts of editor events (typing)
	// into a single Discord update.
	debounceDelay = 500 * time.Millisecond

	// minInterval is the minimum spacing between actual Discord pushes.
	minInterval = 15 * time.Second
)

// Presence owns all state pushed to Discord.
type Presence struct {
	mu sync.Mutex

	// Workspace info (set once at startup).
	workspace workspaceapi.URI

	// Current editor state.
	activePath string
	hasActive bool
	pendTimer *time.Timer

	// lastActivity is the details string of the most recent push, used
	// to skip redundant identical updates. Guarded by mu.
	lastActivity string

	// lastPush is when push() last called SetActivity. Guarded by mu.
	lastPush time.Time

	// startedAt is when the editor session began; used for the Discord
	// elapsed-time timestamp, which persists across file switches.
	startedAt time.Time
}

// New builds a Presence for the given workspace, resolving the workspace
// name and git branch up front.
func New(ctx context.Context, ws *extensionapi.Workspace) (*Presence, error) {
	root, err := wsSchemeRoot(ctx, ws)
	if err != nil {
		return nil, err
	}

	p := &Presence{
		workspace: root,
		startedAt: time.Now(),
	}

	p.push() // initial presence: workspace + branch, no file yet
	return p, nil
}

// wsSchemeRoot resolves the workspace root via the workspace scheme API.
func wsSchemeRoot(ctx context.Context, ws *extensionapi.Workspace) (workspaceapi.URI, error) {
	fs := ws.FileSystem(ctx)
	cwd, err := fs.URI(".")
	if err != nil {
		return workspaceapi.URI{}, err
	}

	return cwd, nil
}

// SetActive marks a file as the one currently focused.
func (p *Presence) SetActive(uri workspaceapi.URI) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Expand to absolute URI before deriving the workspace-relative path.
	absPath, err := workspaceapi.ExpandPathWithURI(uri.Path(), p.workspace)
	if err != nil {
		log.Printf("set activity: error expanding path: %v", err)
		return
	}
	absPathURI := workspaceapi.ParseURI(absPath)
	relPath := workspaceapi.RelPath(p.workspace, absPathURI)
	log.Printf("set activity: abs path: %s", absPathURI.Path())
	log.Printf("set activity: rel path: %s", relPath)
	
	if p.hasActive && p.activePath == relPath {
		return
	}
	p.activePath = relPath
	p.hasActive = true
	p.scheduleUpdate()
}

// Close clears presence when the extension shuts down.
func (p *Presence) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.pendTimer != nil {
		p.pendTimer.Stop()
	}
	client.Logout()
}

// scheduleUpdate debounces pushes to Discord.
func (p *Presence) scheduleUpdate() {
	if p.pendTimer != nil {
		p.pendTimer.Stop()
	}
	p.pendTimer = time.AfterFunc(debounceDelay, p.push)
}

// SetIdle clears the active file so presence shows "Idle" until the
// editor regains focus. Safe to call when already idle.
func (p *Presence) SetIdle() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.hasActive {
		return
	}
	p.activePath = ""
	p.hasActive = false
	p.scheduleUpdate()
}

// push sends the current state to Discord.
func (p *Presence) push() {
	p.mu.Lock()

	if wait := minInterval - time.Since(p.lastPush); wait > 0 {
		// Too soon since the last real push; retry after the cooldown.
		p.pendTimer = time.AfterFunc(wait, p.push)
		p.mu.Unlock()
		return
	}

	activity := client.Activity{
		State:      fmt.Sprintf("Workspace: %s", filepath.Base(p.workspace.Name())),
		SmallImage: runeLogoURL(),
		SmallText:  "Rune Editor",
		Timestamps: &client.Timestamps{
			Start: &p.startedAt,
		},
	}

	if p.hasActive {
		path := p.activePath
		activity.Details = fmt.Sprintf("Editing %s", path)
		activity.LargeImage = fileIcon(path)
	} else {
		activity.Details = "Idle"
		activity.LargeImage = fileIcon("idle")
	}

	// Skip redundant pushes: if the details string is identical to what
	// Discord already shows (e.g. focus returning to the same file),
	// don't burn a rate-limit slot on it.
	if p.lastActivity == activity.Details {
		p.mu.Unlock()
		return
	}

	p.lastActivity = activity.Details
	p.lastPush = time.Now()
	p.mu.Unlock()

	if err := client.SetActivity(activity); err != nil {
		log.Printf("set activity: %v", err)
	}

	log.Printf("set activity success: %+v", activity)
}
