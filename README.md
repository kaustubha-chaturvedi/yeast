## Features

- 🚀 **Simple Plugin System**: Create, install, and run plugins with ease
- 📦 **Plugin Registry**: Browse and install plugins from [yeast.kaustubha.work](https://yeast.kaustubha.work)
- 🔍 **Auto-Discovery**: Automatically discovers plugins in your PATH
- 🛠️ **Plugin Authoring**: Create new plugins with a single command
- 📤 **Publishing**: Publish your plugins to the registry

## Installation

### Quick Install (Recommended)

Run this command in your terminal:

**Linux/macOS:**
```bash
curl -fsSL https://yeast.kaustubha.work/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://yeast.kaustubha.work/install.ps1 | iex
```

The script will:
- Detect your OS and architecture
- Download the latest release
- Install `yst` to a directory in your PATH
- Make it executable

### Manual Installation

1. Download the binary for your platform from [Releases](https://github.com/kaustubha-chaturvedi/yeast/releases)
2. Create `~/.yeast` directory (or `%USERPROFILE%\.yeast` on Windows)
3. Rename the binary to `yst` (or `yst.exe` on Windows) and place it in `~/.yeast`
4. Add `~/.yeast` to your PATH:
   - **Linux/macOS**: Add `export PATH="$PATH:$HOME/.yeast"` to `~/.bashrc` or `~/.zshrc`
   - **Windows**: Add `%USERPROFILE%\.yeast` to your system PATH
5. Make it executable (Linux/macOS): `chmod +x ~/.yeast/yst`

## Usage

### Running Plugins

Once installed, you can run plugins directly using their alias:

```bash
yst <plugin-alias> [args...]
```

For example:
```bash
yst compress-image input.jpg output.jpg
```

### Plugin Management

#### List Installed Plugins

```bash
yst plugins list
```

#### Install a Plugin

```bash
yst plugins install <author:plugin-alias>
```

Example:
```bash
yst plugins install aadarshshrivastava:yst-image-tools
```

#### Publish Your Plugin

From your plugin repository:
```bash
yst plugins publish
```

This will:
- Read metadata from your plugin's `main.go`
- Validate the plugin (name, alias, domain, etc.)
- Check if the alias is available
- Publish to the registry

### Creating a New Plugin

```bash
yst plugins create-new <name> --alias <alias> --domain <domain>
```

Example:
```bash
yst plugins create-new "Image Tools" --alias image-tools --domain image
```

This creates a new plugin skeleton with:
- `main.go` with embedded metadata
- `go.mod` file
- `.github/workflows/release.yaml` for automated releases
- Basic command structure

## Plugin Development

### Plugin Structure

A YEAST plugin is a Go binary that:
- Has embedded metadata (name, alias, domain, version)
- Responds to `__yst_metadata` command with JSON metadata
- Uses Cobra for command structure

### Metadata Requirements

Your plugin's `main.go` must include:
- `name`: Display name of the plugin
- `alias`: Short identifier (used in commands)
- `domain`: Category/domain (e.g., "image", "video", "audio")
- `version`: Semantic version

### Example Plugin

```go
var embeddedMetadata = `{"name":"Image Tools","domain":"image","alias":"image-tools","version":"1.0.0"}`

func getMetadata() map[string]interface{} {
    var meta map[string]interface{}
    json.Unmarshal([]byte(embeddedMetadata), &meta)
    return meta
}
```

### Building and Releasing

Plugins created with `yst plugins create-new` include a GitHub Actions workflow that:
- Builds binaries for Linux, macOS, and Windows (amd64, arm64)
- Creates releases automatically
- Uploads binaries to GitHub Releases

## How It Works

1. **Discovery**: YEAST scans your PATH for binaries matching the pattern `yst-*`
2. **Indexing**: For each found binary, it runs `yst-<alias> __yst_metadata` to get plugin info
3. **Execution**: When you run `yst <alias>`, it finds the plugin binary and executes it with your arguments

## Requirements

- Go 1.25.4+ (for building from source)
- A directory in your PATH (for plugin discovery)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See LICENSE file for details.

## Links

- **Website**: [yeast.kaustubha.work](https://yeast.kaustubha.work)
- **Plugin Registry**: [yeast.kaustubha.work/plugins](https://yeast.kaustubha.work/plugins)
- **GitHub**: [github.com/kaustubha-chaturvedi/yeast](https://github.com/kaustubha-chaturvedi/yeast)
