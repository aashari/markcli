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

### Shell Completion

`markcli` provides shell completion support for:
- bash
- zsh
- fish
- powershell

### Bash

```bash
# Load completions for current session
source <(markcli completion bash)

# Load completions for all sessions
markcli completion bash > /etc/bash_completion.d/markcli
```

### Zsh

```zsh
# Enable shell completion
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Load completions for current session
source <(markcli completion zsh)

# Load completions for all sessions (Linux)
markcli completion zsh > "${fpath[1]}/_markcli"

# Load completions for all sessions (macOS with Homebrew)
markcli completion zsh > $(brew --prefix)/share/zsh/site-functions/_markcli
```

### Fish

```fish
# Load completions for current session
markcli completion fish | source

# Load completions for all sessions
markcli completion fish > ~/.config/fish/completions/markcli.fish
```

### PowerShell

```powershell
# Load completions for current session
markcli completion powershell | Out-String | Invoke-Expression

# Load completions for all sessions
markcli completion powershell > markcli.ps1
```

After enabling completion, you can use <TAB> to auto-complete commands, flags, and arguments.

## Configuration

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

### Global Options
All commands support these global options:
- `--debug`: Enable debug mode for detailed logging
- `--help, -h`: Show help for any command

### Configuration Commands

1. **List Configurations**
   ```bash
   markcli config list
   ```
   Lists all configured platforms and their settings

### Atlassian Commands

#### Site Management

1. **Sites**
   ```bash
   # List all configured sites (default action)
   markcli atlassian sites

   # Set default site
   markcli atlassian sites set-default <site-name>
   ```
   - Default action: Lists all configured Atlassian sites
   - Subcommands:
     - `set-default`: Set the specified site as default
   - Flags:
     - `--help, -h`: Show help for sites command

#### Confluence Commands

1. **List Spaces**
   ```bash
   markcli atlassian confluence spaces [flags]
   ```
   List all available Confluence spaces
   - Flags:
     - `-a, --all`: Show all spaces (including personal and archived)
     - `--site string`: Atlassian site to use (defaults to the default site)

2. **Search Pages**
   ```bash
   markcli atlassian confluence pages search [flags]
   ```
   Search Confluence pages using CQL (Confluence Query Language)
   - Flags:
     - `-q, --query string`: Search query (required)
     - `-s, --space string`: Space key to search in
     - `-l, --limit int`: Number of results per page (default: 10)
     - `-p, --page int`: Page number (default: 1)
     - `--site string`: Atlassian site to use (defaults to the default site)

3. **Get Page Content**
   ```bash
   markcli atlassian confluence pages get [flags]
   ```
   Get a specific Confluence page by ID using Confluence API v2
   - Flags:
     - `--id string`: Page ID to retrieve
     - `--site string`: Atlassian site to use (defaults to the default site)

   Example:
   ```bash
   markcli atlassian confluence pages get --id 123456
   ```

#### Jira Commands

1. **List Projects**
   ```bash
   markcli atlassian jira projects [flags]
   ```
   List all available Jira projects
   - Flags:
     - `--site string`: Atlassian site to use (defaults to the default site)
     - `--sort string`: Sort projects by: key, name, type, or style (default: "key")

2. **Search Issues**
   ```bash
   markcli atlassian jira issues search [flags]
   ```
   Search Jira issues using text query with optional project filtering
   - Flags:
     - `-q, --query string`: Search query (required)
     - `-l, --limit int`: Number of results per page (default: 10)
     - `-p, --page int`: Page number (default: 1)
     - `-r, --project string`: Project key to filter issues
     - `--site string`: Atlassian site to use (defaults to the default site)

3. **Get Issue Details**
   ```bash
   markcli atlassian jira issues get [flags]
   ```
   Get a specific Jira issue by ID using Jira API v3
   - Flags:
     - `--id string`: Issue ID to retrieve
     - `--site string`: Atlassian site to use (defaults to the default site)

### Common Usage Patterns

1. **Using a Specific Site**
   All Atlassian commands accept the `--site` flag to specify which site configuration to use. If not provided, the default site is used.
   ```bash
   markcli atlassian jira projects --site="your-site-name"
   ```

2. **Pagination**
   Search commands support pagination with `--limit` and `--page` flags:
   ```bash
   markcli atlassian jira issues search -q "high priority" -l 20 -p 2
   ```

3. **Debug Mode**
   Add `--debug` to any command for detailed logging:
   ```bash
   markcli --debug atlassian confluence spaces
   ```

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
