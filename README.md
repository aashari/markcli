# markcli

A command-line tool for managing markdown content across different platforms. Currently supports Atlassian Confluence and Jira.

## Quick Start

1.  **Prerequisites**
    *   Go 1.22.1 or later (only needed for building from source)
    *   Git
    *   Atlassian account with API token ([How to get an API token](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens/))

2.  **Installation**

    Choose one of the following installation methods:

    **Option 1: Download Pre-built Binary (Recommended)**
    
    1. Navigate to the [releases page](https://github.com/andifg/markcli/releases) and download the appropriate executable for your operating system:
       - `markcli-linux` for Linux
       - `markcli-mac` for macOS
       - `markcli.exe` for Windows

    2. For Linux & macOS:
       ```bash
       # Make the downloaded file executable
       chmod +x markcli-linux  # or markcli-mac for macOS
       
       # Move to a directory in your PATH
       sudo mv markcli-linux /usr/local/bin/markcli  # or markcli-mac for macOS
       ```

    3. For Windows:
       - Rename the downloaded file to `markcli.exe` if needed
       - Move it to a directory in your system's PATH (e.g., `C:\Windows`)
       - Or add its location to your PATH environment variable

    **Option 2: Build from Source**
    ```bash
    git clone https://github.com/andifg/markcli.git
    cd markcli
    go build ./cmd/markcli
    ```

3.  **Configuration**

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

## Development

For detailed development instructions, see [DEVELOPMENT.md](DEVELOPMENT.md).

Quick development workflow:

```bash
# Run without building
go run cmd/markcli/main.go [command]

# Run tests
go test ./...

# Build for development
go build ./cmd/markcli
```

## Documentation

-   [DEVELOPMENT.md](DEVELOPMENT.md) - Comprehensive development guide
-   [COMMAND.md](COMMAND.md) - Complete command reference (Coming Soon)

## Features

*   **Atlassian Confluence Support:**
    *   List Confluence spaces with filtering options.
    *   Search Confluence pages using CQL (Confluence Query Language).
    *   Get the content of a specific Confluence page by its ID.

*   **Atlassian Jira Support:**
    *   List Jira projects.
    *   Search Jira issues using a text query, with optional project filtering.
    *   Get the details of a specific Jira issue by its ID.

*   **Atlassian Site Management:**
    *   Manage configurations for multiple Atlassian sites.
    *   Set a default Atlassian site.

*   **Common Features:**
    *   Consistent markdown output across all supported platforms.
    *   Support for pagination when listing spaces or searching pages/issues.
    *   Robust error handling and informative error messages.
    *   Debug logging is available using the `--debug` flag.

## Basic Usage

1.  **List Confluence Spaces**

    ```bash
    markcli atlassian confluence spaces
    ```

    You can also list all spaces including personal and archived spaces by using the following command

    ```bash
    markcli atlassian confluence spaces -a
    ```

2.  **Search Confluence Pages**

    ```bash
    markcli atlassian confluence pages search --query "your search query" --space "space key"
    ```

    Replace `your search query` with your query and `space key` with the space in which you want to search.

3.  **Get Confluence Page Content**

    ```bash
    markcli atlassian confluence pages get --id <page-id>
    ```

    Replace `<page-id>` with the ID of the page you want to retrieve.

4. **List Jira Projects**

   ```bash
   markcli atlassian jira projects
    ```
You can sort the output using `--sort` flag, for example use `--sort name` to sort by name, or `--sort type` to sort by project type.

   ```bash
  markcli atlassian jira projects --sort name
   ```

5.  **Search Jira Issues**

   ```bash
   markcli atlassian jira issues search --query "your text query"
   ```

    Replace `"your text query"` with your search term.

   You can also filter the search by a project key:
   ```bash
  markcli atlassian jira issues search --query "your text query" -r <project key>
   ```
    Use `-r` flag to filter the search in a project using its `project key`.

   You can specify pagination using `--page` and `--limit`:

   ```bash
   markcli atlassian jira issues search -q "your text query" -l 20 -p 2
   ```
    Use `-l` to limit the number of records to show, and `-p` to go to a specific page.

6.  **Get Jira Issue Details**
    ```bash
    markcli atlassian jira issues get --id <issue-id>
    ```
    Replace `<issue-id>` with the ID of the Jira issue you want to retrieve.

**Common Options**

*  All commands accept a `--site` parameter to specify which Atlassian site configuration to use. If not provided, the tool will try to use the default site.
   ```bash
  markcli atlassian jira projects  --site="your-site-name"
  ```

*  All commands also accepts a `--debug` flag which can be useful for debugging.

**Configuration Examples**

* To add a new Atlassian configuration:
 ```bash
 markcli config add atlassian
 ```
* To list all available configurations:
 ```bash
  markcli config list
 ```
* To remove a configuration:
```bash
  markcli config remove atlassian <site-name>
```

* Set the default Atlassian site:
```bash
  markcli atlassian sites set-default <site-name>
```
* List Atlassian sites:
 ```bash
   markcli atlassian sites list
 ```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
