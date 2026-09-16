package presence

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// File icon resolution based on the vyfor/icons asset set (also used by
// cord.nvim). Icons are served straight from the repository, so no assets
// have to be uploaded to the Discord application.
//
// URL layout:
//   https://raw.githubusercontent.com/vyfor/icons/master/icons/<theme>/<flavor>/<name>.png

const (
	iconsBase = "https://raw.githubusercontent.com/vyfor/icons/master/icons/"

	// iconVersion is appended as a cache-busting query, mirroring
	// cord.nvim's ICONS_VERSION. Bump when the upstream set changes.
	iconVersion = "?v=26"

	// fallbackIcon is used for extensions we don't map.
	fallbackIcon = "book"

	// runeLogo is the small-image asset. If you upload a Rune logo to the
	// Discord application's Rich Presence assets, put its key here; it can
	// also be an external URL.
	runeLogo = "rune-default-dark"
)

// iconTheme and iconFlavor select which icon set is served. Both are
// runtime-configurable via SetIconTheme; they must only be read/written
// while holding iconsMu (or before any concurrent use).
var (
	iconsMu    sync.Mutex
	iconTheme  = "dark"
	iconFlavor = "default"
)

// validFlavors are the icon packs published by vyfor/icons.
var validFlavors = map[string]bool{
	"atom":       true,
	"catppuccin": true,
	"classic":    true,
	"default":    true,
	"minecraft":  true,
	"void":       true,
}

// validThemes are the color variants published per flavor.
var validThemes = map[string]bool{
	"dark":   true,
	"light":  true,
	"accent": true,
}

// SetIconTheme switches the icon flavor/theme used for new pushes.
// Errors report which key was invalid so the caller can log it.
func SetIconTheme(flavor, theme string) error {
	if !validFlavors[flavor] {
		return fmt.Errorf("invalid icon flavor %q (valid: %s)", flavor, strings.Join(sortedKeys(validFlavors), ", "))
	}
	if !validThemes[theme] {
		return fmt.Errorf("invalid icon theme %q (valid: %s)", theme, strings.Join(sortedKeys(validThemes), ", "))
	}
	iconsMu.Lock()
	defer iconsMu.Unlock()
	iconFlavor = flavor
	iconTheme = theme
	return nil
}

func sortedKeys(m map[string]bool) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// iconURL returns the CDN URL for an icon name.
func iconURL(name string) string {
	iconsMu.Lock()
	defer iconsMu.Unlock()
	return iconsBase + iconFlavor + "/" + iconTheme + "/" + name + ".png" + iconVersion
}

// filenameIcons maps well-known filenames (basename, lowercased) to icons.
var filenameIcons = map[string]string{
	"makefile":           "make",
	"dockerfile":         "docker",
	"justfile":           "just",
	"cmakelists.txt":     "cmake",
	"license":            "license",
	"cargo.lock":         "cargo",
	"cargo.toml":         "cargo",
	"go.mod":             "go",
	"go.sum":             "go",
	"package.json":       "npm",
	"package-lock.json":  "npm",
	"tsconfig.json":      "typescript",
	"gemfile":            "ruby",
	"gemfile.lock":       "ruby",
	"rakefile":           "ruby",
	".gitignore":         "git",
	".gitattributes":     "git",
	".gitmodules":        "git",
	".env":               "database",
	"requirements.txt":   "python",
	"pyproject.toml":     "python",
	"setup.py":           "python",
	"gradlew":            "gradle",
	"build.gradle":       "gradle",
	"settings.gradle":    "gradle",
	"pom.xml":            "maven",
	".editorconfig":      "editorconfig",
	"flake.nix":          "nix",
	"flake.lock":         "nix",
	"vite.config.ts":     "typescript",
	"tailwind.config.js": "css",
	"astro.config.mjs":   "astro",
	"svelte.config.js":   "svelte",
	"vue.config.js":      "vue",
	"next.config.js":     "react",
	"init.lua":           "neovim",
	"init.vim":           "neovim",
}

// extIcons maps file extensions (without dot, lowercased) to icons.
// Names must exist in the vyfor/icons set (theme "default").
var extIcons = map[string]string{
	// Ada
	"adb": "ada", "ads": "ada",
	// Assembly
	"asm": "assembly", "s": "assembly", "nasm": "assembly",
	// Astro
	"astro": "astro",
	// AutoHotkey
	"ahk": "autohotkey", "ahk2": "autohotkey",
	// AWK
	"awk": "awk",
	// Bolt
	"bolt": "bolt", "bsl": "bolt",
	// C / C++
	"c": "c", "h": "c",
	"cc": "cpp", "cpp": "cpp", "cxx": "cpp", "hpp": "cpp", "hh": "cpp",
	// Clojure
	"clj": "clojure", "cljs": "clojure", "cljc": "clojure", "edn": "clojure",
	// CMake
	"cmake": "cmake",
	// Crystal
	"cr": "crystal",
	// C# / F#
	"cs": "csharp", "csx": "csharp",
	"fs": "fsharp", "fsi": "fsharp", "fsx": "fsharp",
	// CSS / Sass / Less
	"css": "css", "scss": "sass", "sass": "sass", "less": "less", "postcss": "postcss", "pcss": "postcss",
	// D
	"d": "d", "di": "d",
	// Dart
	"dart": "dart",
	// Django (Python web framework templates)
	"djt": "django",
	// Docker
	"dockerfile": "docker",
	// Elixir
	"eex": "elixir", "ex": "elixir", "exs": "elixir",
	// Elm
	"elm": "elm",
	// Erlang
	"erl": "erlang", "hrl": "erlang",
	// Fennel
	"fnl": "fennel",
	// Fish shell
	"fish": "fishshell",
	// Fortran
	"f90": "fortran", "f95": "fortran", "f03": "fortran",
	// GameMaker
	"gml": "gamemaker",
	// Git
	"gitignore": "git", "gitconfig": "git",
	// Gleam
	"gleam": "gleam",
	// Go
	"go": "go",
	// Godot
	"gd": "godot", "godot": "godot", "tres": "godot", "tscn": "godot",
	// Gradle / Groovy
	"gradle": "gradle", "groovy": "groovy",
	// GraphQL
	"graphql": "graphql", "gql": "graphql",
	// Haskell
	"hs": "haskell", "lhs": "haskell",
	// Haxe
	"hx": "haxe",
	// HTML
	"html": "html", "htm": "html", "xhtml": "html",
	// Java
	"java": "java",
	// JavaScript / TypeScript
	"js": "javascript", "mjs": "javascript", "cjs": "javascript", "jsx": "javascript",
	"ts": "typescript", "tsx": "typescript", "mts": "typescript", "cts": "typescript",
	// JSON
	"json": "json", "jsonc": "json", "json5": "json",
	// Julia
	"jl": "julia",
	// Jupyter
	"ipynb": "jupyter",
	// Kotlin
	"kt": "kotlin", "kts": "kotlin",
	// LaTeX / Typst / Tex
	"tex": "latex", "sty": "latex", "typ": "typst",
	// Lock files
	"lock": "lock",
	// Logs
	"log": "log",
	// Lua / Luau
	"lua": "lua", "luau": "luau",
	// Markdown
	"md": "markdown", "markdown": "markdown", "mdx": "markdown",
	// MATLAB
	"mat": "matlab", "m": "matlab",
	// MCFunction (Minecraft)
	"mcfunction": "mcfunction",
	// Nim
	"nim": "nim", "nims": "nim",
	// Nix
	"nix": "nix",
	// Objective-C
	"mm": "objc", "m+": "objc",
	// OCaml
	"ml": "ocaml", "mli": "ocaml",
	// Odin
	"odin": "odin",
	// Org
	"org": "org",
	// Pascal
	"pas": "pascal", "pp": "pascal",
	// Perl
	"pl": "perl", "pm": "perl",
	// PHP
	"php": "php",
	// Powershell
	"ps1": "powershell", "psm1": "powershell", "psd1": "powershell",
	// Prisma
	"prisma": "prisma",
	// Python
	"py": "python", "pyi": "python", "pyw": "python",
	// R / Racket
	"r": "r", "rmd": "r",
	"rkt": "racket", "rktl": "racket",
	// Ruby
	"rb": "ruby", "erb": "ruby", "rake": "ruby", "gemspec": "rubygems",
	// Rust
	"rs": "rust",
	// Scala
	"sbt": "scala", "scala": "scala",
	// Scheme / Lisp
	"scm": "scheme", "ss": "scheme", "lisp": "lisp", "lsp": "lisp", "cl": "lisp",
	// Shader
	"glsl": "shader", "frag": "shader", "vert": "shader", "comp": "shader", "wgsl": "shader",
	// Shell
	"sh": "shell", "bash": "shell", "zsh": "shell", "ksh": "shell", "nushell": "nushell", "nu": "nushell",
	// Solidity
	"sol": "solidity",
	// SQL / databases
	"sql": "database", "db": "database", "sqlite": "database", "sqlite3": "database",
	// Squirrel
	"nut": "squirrel",
	// Svelte
	"svelte": "svelte",
	// SVG
	"svg": "svg",
	// Swift
	"swift": "swift",
	// Terraform / HCL
	"tf": "terraform", "tfvars": "terraform", "hcl": "terraform",
	// TOML / YAML / XML
	"toml": "toml", "yaml": "yaml", "yml": "yaml", "xml": "xml",
	// Vim / Neovim
	"vim": "vim", "viml": "viml",
	// V
	"v": "v", "vsh": "v",
	// Vala
	"vala": "vala",
	// Vue
	"vue": "vue",
	// WebAssembly
	"wasm": "wasm", "wat": "wasm",
	// Zig
	"zig": "zig",
}

// fileIcon resolves the icon URL for a file path. Exact filenames win over
// extensions (so Cargo.toml gets the cargo icon, not the generic toml one).
func fileIcon(path string) string {
	base := strings.ToLower(filepath.Base(path))
	if name, ok := filenameIcons[base]; ok {
		return iconURL(name)
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	if name, ok := extIcons[ext]; ok {
		return iconURL(name)
	}
	return iconURL(fallbackIcon)
}

// runeLogoURL is the URL used for the small-image key. Point this at a
// hosted Rune logo, or upload an asset to the Discord application and
// return its key instead.
func runeLogoURL() string {
	return runeLogo
}
