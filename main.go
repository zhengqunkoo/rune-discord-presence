// Command rune-discord-presence is a Rune editor extension that publishes
// the current workspace, active file, cursor line and git branch as
// Discord Rich Presence.
package main

import (
	"context"
	"log"

	"github.com/hugolgst/rich-go/client"
	"github.com/unstablebuild/rune-go-sdk/api/config"
	"github.com/unstablebuild/rune-go-sdk/api/extensionapi"
	"github.com/unstablebuild/rune-go-sdk/api/textapi"

	"github.com/ravener/rune-discord-presence/presence"
)

// Replace with the Application ID (Client ID) of YOUR Discord application:
// https://discord.com/developers/applications -> New Application -> OAuth2
const discordAppID = "1548395071664291910"

func main() {
	meta := extensionapi.Metadata{
		DeveloperID:      "ravener",
		DeveloperEmail:   "ravener.anime@gmail.com",
		DeveloperKey:     "1234",
		ExtensionID:      "rune-discord-presence",
		ExtensionName:    "Discord Rich Presence",
		ExtensionVersion: "0.1.0",
		Permissions: extensionapi.NewPermissions(
			extensionapi.PermissionEditor,
			extensionapi.PermissionFileSystem,
		),
	}

	err := extensionapi.ServeWorkspaceExtension(
		extensionapi.FuncWorkspaceExtension(extend), meta,
	)
	if err != nil {
		log.Fatalf("rune-discord-presence: %v", err)
	}
}

func extend(ctx context.Context, ws *extensionapi.Workspace, cfg config.Config) error {
	// Icon flavor/theme from user config. Defaults apply when unset;
	// invalid values fall back to defaults with a warning.
	flavor := "default"
	if v, err := cfg.GetString("flavor"); err == nil && v != "" {
		flavor = v
	}
	theme := "dark"
	if v, err := cfg.GetString("theme"); err == nil && v != "" {
		theme = v
	}
	if err := presence.SetIconTheme(flavor, theme); err != nil {
		log.Printf("config: %v; using defaults", err)
	} else {
		log.Printf("using flavor: %s theme: %s", flavor, theme)
	}

	if err := client.Login(discordAppID); err != nil {
		return err
	}
	log.Print("Discord logged in")

	p, err := presence.New(ctx, ws)
	if err != nil {
		return err
	}
	defer p.Close()

	h := &presence.EventHandler{Pres: p}
	events := []textapi.EventType{
		textapi.EventTypeOpen,
		textapi.EventTypeFocus,
		textapi.EventTypeUnfocus,
	}
	if err := ws.Editor(ctx).SubscribeEvents(events, h); err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}
