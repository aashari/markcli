# markcli: Your Markdown CLI Companion for Atlassian

`markcli` is a powerful command-line interface (CLI) tool designed to streamline your interactions with Atlassian products, specifically Confluence and Jira. It empowers you to manage and access markdown content across these platforms directly from your terminal.

## Key Features

- **Cross-Platform Compatibility:** Runs seamlessly on macOS, Linux, and Windows.
- **Unified Markdown Output:** Consistent markdown formatting across Confluence and Jira.
- **Atlassian Confluence Integration:**
  - List, search, and filter Confluence spaces.
  - Retrieve detailed information and markdown content of Confluence pages.
  - Search Confluence pages using CQL (Confluence Query Language).
- **Atlassian Jira Integration:**
  - List Jira projects and sort results by key, name, type, or style.
  - Search for Jira issues using text queries and project filters.
  - Get detailed information and comments for specific Jira issues.
- **Global Search:** Search across both Confluence pages and Jira issues.
- **Multiple Site Support:** Easily manage and switch between multiple Atlassian site configurations.
- **Pagination:** Efficiently handle large datasets with pagination for list and search commands.
- **Debug Mode:** Enable detailed logging for troubleshooting.
- **Shell Completion:** Available for Bash, Zsh, Fish, and PowerShell.

## Installation

### Binary Download

1.  Download the appropriate binary for your operating system from the [releases page](https://github.com/aashari/markcli/releases).

    - `markcli-linux` for Linux
    - `markcli-mac` for macOS
    - `markcli.exe` for Windows

2.  **Make the binary executable (Linux/macOS):**

    ```bash
    chmod +x ./markcli-*
    ```

3.  **Move the binary to a directory in your `$PATH`:**

    - **Linux/macOS:**

      ```bash
      # System-wide installation (requires sudo)
      sudo mv markcli-* /usr/local/bin/markcli

      # User-local installation
      mkdir -p ~/.local/bin
      mv markcli-* ~/.local/bin/markcli
      ```

    - **Windows:**

      - Rename the downloaded file to `markcli.exe`.
      - Move it to a directory included in your system's `PATH` environment variable (e.g., `C:\Windows`), or add the directory containing the binary to the `PATH`.

### Verify Installation

Confirm the installation by running:

```bash
markcli --version
```

### Shell Completion

`markcli` offers convenient shell completion to speed up your command entry:

#### Bash

```bash
# Load completion for current session
source <(markcli completion bash)

# Load completion for all sessions
markcli completion bash > /etc/bash_completion.d/markcli
```

#### Zsh

```zsh
# Enable completion (add to ~/.zshrc and restart terminal)
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Load completion for current session
source <(markcli completion zsh)

# Load completion for all sessions (Linux)
markcli completion zsh > "${fpath[1]}/_markcli"

# Load completion for all sessions (macOS with Homebrew)
markcli completion zsh > $(brew --prefix)/share/zsh/site-functions/_markcli
```

#### Fish

```fish
# Load completion for current session
markcli completion fish | source

# Load completion for all sessions
markcli completion fish > ~/.config/fish/completions/markcli.fish
```

#### PowerShell

```powershell
# Load completion for current session
markcli completion powershell | Out-String | Invoke-Expression

# Load completion for all sessions
markcli completion powershell > markcli.ps1
```

## Configuration

`markcli` stores its configuration data in `~/.config/markcli/config.json`.

### Adding New Configurations

You can add or modify configurations using the `config` command:

```bash
markcli config add atlassian
```

Follow the prompts to configure a new Atlassian site. You will be asked for:

- The Atlassian site URL or name
- Your email address
- Your API token

### Manual Configuration

You can also directly edit the `config.json` file. Example:

```json
{
  "atlassian": {
    "my_site": {
      "site_name": "my_site",
      "base_url": "https://my_site.atlassian.net",
      "email": "user@example.com",
      "token": "your-api-token"
    }
  },
   "default_atlassian_site": "my_site"
}
```

**Important:** Do not share this file, as it contains your API token.

## Command Reference

### Global Options

- `--debug`: Enable debug logging.
- `--help`, `-h`: Display help for any command.

### Configuration Commands

- **`markcli config add <platform>`**: Add a new platform configuration.
- **`markcli config list`**: List configured platforms and settings.
- **`markcli config remove <platform> <site-name (optional)>`**: Remove platform configuration.

  - For example `markcli config remove atlassian my_site`

### Atlassian Commands

#### Site Management

- **`markcli atlassian sites`**: Lists all configured Atlassian sites.

- **`markcli atlassian sites set-default <site-name>`**: Set the default Atlassian site.

#### Global Search

- **`markcli atlassian search [flags]`**: Search across Confluence and Jira.

  **Flags:**

  - `-t, --text <string>`: Search text (required).
  - `-l, --limit <int>`: Number of results per page (default: 100).
  - `-p, --page <int>`: Page number (default: 1).
  - `--site <string>`: Atlassian site to use (defaults to the default site).

  **Example:**

  ```bash
  markcli atlassian search -t "meeting notes" -l 5
  markcli atlassian search -t "API documentation" -l 10 -p 2 --site my_site
  ```

#### Confluence Commands

- **`markcli atlassian confluence spaces [flags]`**: List Confluence spaces.

  **Flags:**

  - `-a, --all`: Show all spaces, including personal and archived ones.
  - `--site <string>`: Atlassian site to use (defaults to the default site).

- **`markcli atlassian confluence pages [flags]`**: List and search Confluence pages.

  **Flags:**

  - `--space <string>`: Space key to list pages from (e.g., "TEAM").
  - `-t, --text <string>`: Search text (for search subcommand).
  - `-l, --limit <int>`: Number of results per page (default: 100).
  - `-p, --page <int>`: Page number (default: 1).
  - `--site <string>`: Atlassian site to use (defaults to the default site).

  **Examples:**

  ```bash
  # List pages from a space
  markcli atlassian confluence pages --space TEAM

  # Search pages with text
  markcli atlassian confluence pages search -t "deployment" -l 20
  ```

- **`markcli atlassian confluence pages get [flags]`**: Get a specific Confluence page.

  **Flags:**

  - `--id <string>`: Page ID to retrieve.
  - `--site <string>`: Atlassian site to use (defaults to the default site).

  **Example:**

  ```bash
  markcli atlassian confluence pages get --id 123456
  ```

#### Jira Commands

- **`markcli atlassian jira projects [flags]`**: List Jira projects.

  **Flags:**

  - `--site <string>`: Atlassian site to use (defaults to the default site).
  - `--sort <string>`: Sort by `key`, `name`, `type`, or `style` (default: `key`).

- **`markcli atlassian jira issues [flags]`**: List and search Jira issues.

  **Flags:**

  - `--project <string>`: Project key to list issues from (e.g., "SHOP").
  - `-t, --text <string>`: Search text (for search subcommand).
  - `-l, --limit <int>`: Number of results per page (default: 100).
  - `-p, --page <int>`: Page number (default: 1).
  - `--site <string>`: Atlassian site to use (defaults to the default site).

  **Note:** Issues are ordered by last updated date in descending order, and statuses "Abandoned" and "Done" are excluded by default.

  **Examples:**

  ```bash
  # List issues from a project
  markcli atlassian jira issues --project SHOP

  # Search issues with text
  markcli atlassian jira issues search -t "deployment" -l 20
  ```

- **`markcli atlassian jira issues get [flags]`**: Get a specific Jira issue.

  **Flags:**

  - `--id <string>`: Issue ID to retrieve.
  - `--site <string>`: Atlassian site to use (defaults to the default site).

  **Example:**

  ```bash
  markcli atlassian jira issues get --id PROJ-123
  ```

## Usage Patterns

- **Specifying a Site:**

  ```bash
  markcli atlassian jira projects --site="your_site_name"
  ```

- **Pagination:**

  ```bash
  markcli atlassian jira issues search -t "bug" -l 20 -p 2
  ```

- **Debug Mode:**

  ```bash
  markcli --debug atlassian confluence spaces
  ```

## Contributing

Contributions are welcome! To contribute:

1.  Fork the repository.
2.  Create a new branch for your feature: (`git checkout -b feature/your-feature`).
3.  Commit your changes: (`git commit -m 'Add your feature'`).
4.  Push to your branch: (`git push origin feature/your-feature`).
5.  Create a pull request.

We appreciate your contributions!
