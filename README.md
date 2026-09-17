# Rune Discord Presence

<img width="734" height="310" alt="rune-discord-presence" src="https://github.com/user-attachments/assets/a8644f7f-ff1f-4ec6-a467-6d58beaa2537" />

Discord Rich Presence extension for [Rune](https://rune.build)

## Install

Open up the `console` from the Rune Command Prompt

```
pkg install github.com/ravener/rune-discord-presence
```

## Configuration

Icons are used from [vyfor/icons](https://github.com/vyfor/icons) and the used flavor/theme combination can be configured in your Rune's `config.yaml`

```yaml
extensions:
  rune-discord-presence:
    path: ...
    config:
      flavor: default # default, atom, catppuccin, classic, minecraft, void
      theme: dark # dark, light, accent
```

## Contributing

This is still an early release and may not behave perfectly, please open an issue if you find any bugs, Pull Requests are also welcome.

## License

[MIT](LICENSE)
