# Key Rotator 🔑

Key Rotator is a Go program designed to manage and rotate secrets for GitHub repositories. It reads a YAML configuration file specifying multiple secrets and their destinations, prompts the user for the secret values, encrypts them using the repository's public key, and updates the secrets in the specified GitHub repositories.

## Installation

### Homebrew

```shell
brew install lucasmelin/tap/key-rotator
```

### Build from source

#### Prerequisites

- Go installed. See the [official installation guide](https://golang.org/doc/install) for instructions.

#### Build steps

```shell
# Clone the repository via HTTPS
git clone https://github.com/lucasmelin/key-rotator.git

# Change into the repository directory
cd key-rotator

# Download the dependencies
go mod tidy

# Install key-rotator in your $GOPATH/bin directory
go install .
```

## Configuration

The configuration file is a YAML file that specifies the secrets and their destinations. Below is an example `key.yaml` file:

```yaml
secrets:
  - name: "SERVICE_ACCOUNT_KEY"
    description: "Secret key for the service account"
    destinations:
      - name: "SERVICE_ACCOUNT_KEY"
        type: "github-repository"
        repo: "lucasmelin/key-rotator"
      - name: "SERVICE_ACCOUNT_KEY"
        type: "github-repository-dependabot"
        repo: "lucasmelin/key-rotator"
      - name: "SERVICE_ACCOUNT_KEY"
        type: "github-repository-environment"
        repo: "lucasmelin/key-rotator"
        environment: "prod"
  - name: "SERVICE_ACCOUNT_USERNAME"
    description: "Secret username for the service account"
    destinations:
      - name: "SERVICE_ACCOUNT_USERNAME"
        type: "github-repository"
        repo: "lucasmelin/key-rotator"
```

## Usage

1. Navigate to the directory containing your YAML configuration file.

2. Ensure you have the `GITHUB_TOKEN` environment variable set:

   ```sh
   export GITHUB_TOKEN=your_github_token
   ```

3. Run `key-rotator rotate` with the path to your YAML configuration file:

   ```sh
   key-rotator rotate path/to/your/key.yaml
   # Or optionally, use the dry-run flag to see what
   # changes would be made without updating the secrets.
   key-rotator rotate --dry-run path/to/your/key.yaml
   ```
   
4. Follow the prompts to rotate all the secrets defined in your configuration file. To cancel the program, press <kbd>Ctrl</kbd>+<kbd>c</kbd>.

## License

This project is licensed under the MIT License. See the [`LICENSE` file](./LICENSE) for details.
