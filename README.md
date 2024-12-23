# markcli

A command-line tool for managing markdown content across different platforms. Currently supports Atlassian Confluence and Jira.

## Support

Currently markcli supports the following OSes:

* macOS
  * 64bit (amd64)
  * Arm (Apple Silicon)
* Linux
  * 64bit (amd64)
  * Arm64
* Windows
  * 64bit (amd64)

## Installation

### Automatic

#### Using Homebrew (macOS and Linux)

```bash
brew install aashari/tap/markcli
```

### Manual

1. Check out markcli into any path (here is `${HOME}/.markcli`)

```bash
git clone --depth=1 https://github.com/aashari/markcli.git ~/.markcli
```

2. Add `~/.markcli/bin` to your `$PATH` any way you like

For bash:
```bash
echo 'export PATH="$HOME/.markcli/bin:$PATH"' >> ~/.bash_profile
```

For zsh:
```bash
echo 'export PATH="$HOME/.markcli/bin:$PATH"' >> ~/.zprofile
```

For fish:
```bash
set -Ux fish_user_paths $HOME/.markcli/bin $fish_user_paths
```

**Alternative**: You can make symlinks into a directory that is already in your `$PATH` (e.g. `/usr/local/bin`) *OSX/Linux Only!*

```bash
# For macOS and Linux users
sudo ln -s ~/.markcli/markcli /usr/local/bin/markcli

# For Ubuntu/Debian users (local user installation)
mkdir -p ~/.local/bin
ln -s ~/.markcli/markcli ~/.local/bin/markcli
. ~/.profile
```

### Binary Download

1. Download the latest version from the [releases page](https://github.com/aashari/markcli/releases) for your platform:
   - `markcli-linux` for Linux
   - `markcli-mac` for macOS
   - `markcli.exe` for Windows

2. Make it executable (Linux/macOS):
```bash
chmod +x ./markcli-*
```

3. Move to a directory in your `$PATH`:

Linux/macOS:
```bash
# System-wide installation (requires sudo)
sudo mv markcli-* /usr/local/bin/markcli

# Or, user-local installation
mkdir -p ~/.local/bin
mv markcli-* ~/.local/bin/markcli
```

Windows:
- Rename the downloaded file to `markcli.exe`
- Move it to a directory in your system's PATH (e.g., `C:\Windows`)
- Or add its location to your PATH environment variable

### Verify Installation

After installation, verify that markcli is properly installed:

```bash
markcli --version
```

## Quick Start

1.  **Prerequisites**
    *   Go 1.22.1 or later (only needed for building from source)
    *   Git
    *   Atlassian account with API token ([How to get an API token](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens/))

2.  **Configuration**

    `markcli` stores its configuration in `~/.config/markcli/config.json`. You can add new configurations using the `config add` command. For example:

    ```bash
    markcli config add atlassian
    ```

    This will prompt you for your Atlassian site URL or name, email, and API token.

    You can also manually create or edit the `config.json` file directly. A sample configuration would look like this:

    ```json
    {
      "atlassian": {
        "sitename": {
          "site_name": "sitename",
          "base_url": "https://your-site.atlassian.net",
          "email": "your-email@example.com",
          "token": "your-api-token"
        }
      },
      "default_atlassian_site": "sitename"
    }
    ```

## Command Reference

### Global Flags
- `--debug`: Enable debug mode for detailed logging

### Configuration Commands
```bash
# Add a new configuration
markcli config add [platform]     # Currently supports: atlassian

# List all configurations
markcli config list

# Remove a configuration
markcli config remove [platform] [site-name]
```

### Atlassian Commands

#### Site Management
```bash
# List all configured sites
markcli atlassian sites list

# Set default site
markcli atlassian sites set-default [site-name]
```

#### Confluence Commands

1. **List Spaces**
   ```bash
   # List all spaces
   markcli atlassian confluence spaces
   
   # Include personal and archived spaces
   markcli atlassian confluence spaces --all
   
   # Use specific site
   markcli atlassian confluence spaces --site "your-site"
   ```

2. **Search Pages**
   ```bash
   # Basic search
   markcli atlassian confluence pages search --query "your search query"
   
   # Search in specific space with pagination
   markcli atlassian confluence pages search \
     --query "your search query" \
     --space "SPACE_KEY" \
     --limit 20 \
     --page 2
   ```

3. **Get Page Content**
   ```bash
   # Get specific page by ID
   markcli atlassian confluence pages get --id "page-id"
   ```

#### Jira Commands

1. **List Projects**
   ```bash
   # List all projects
   markcli atlassian jira projects
   
   # Sort projects by name
   markcli atlassian jira projects --sort name
   
   # Available sort options: key, name, type, style
   ```

2. **Search Issues**
   ```bash
   # Basic search
   markcli atlassian jira issues search --query "high priority"
   
   # Search with project filter and pagination
   markcli atlassian jira issues search \
     --query "high priority" \
     --project "PROJ" \
     --limit 20 \
     --page 2
   ```

3. **Get Issue Details**
   ```bash
   # Get specific issue
   markcli atlassian jira issues get --id "PROJ-123"
   ```

### Common Options

All Atlassian commands accept these common flags:
- `--site`: Specify which Atlassian site configuration to use (defaults to the default site)
- `--debug`: Enable debug output for troubleshooting

### Search Command Options

All search commands (Confluence pages, Jira issues) support:
- `--query, -q`: Search query (required)
- `--limit, -l`: Results per page (default: 10)
- `--page, -p`: Page number (default: 1)

## Features

*   **Atlassian Confluence Support:**
    *   List Confluence spaces with filtering options
    *   Search Confluence pages using CQL (Confluence Query Language)
    *   Get the content of a specific Confluence page by its ID

*   **Atlassian Jira Support:**
    *   List Jira projects with sorting options
    *   Search Jira issues using text queries and project filters
    *   Get detailed information about specific Jira issues

*   **Atlassian Site Management:**
    *   Manage configurations for multiple Atlassian sites
    *   Set and use default Atlassian site
    *   Easy switching between different sites

*   **Common Features:**
    *   Consistent markdown output across all supported platforms
    *   Support for pagination in list and search operations
    *   Robust error handling and informative error messages
    *   Debug logging for troubleshooting

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
