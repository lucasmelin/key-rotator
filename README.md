# Key Rotator 🔑

Key Rotator is a Go program designed to manage and rotate secrets for GitHub repositories. It reads a YAML configuration file specifying multiple secrets and their destinations, prompts the user for the secret values, encrypts them using the repository's public key, and updates the secrets in the specified GitHub repositories.

## Prerequisites

- Go installed. See the [official installation guide](https://golang.org/doc/install) for instructions.
- [`libsodium`](https://libsodium.gitbook.io/doc/) installed. On macOS, you can install it using Homebrew:

  ```sh
  brew install libsodium
  ```
- [`pkg-config`](https://www.freedesktop.org/wiki/Software/pkg-config/) installed. On macOS, you can install it using Homebrew:

  ```sh
  brew install pkg-config
  ```
- A GitHub personal access token with the necessary permissions to update secrets.
- The `GITHUB_TOKEN` environment variable set with your GitHub token.

## Configuration

The configuration file is a YAML file that specifies the secrets and their destinations. Below is an example `key.yaml` file:

```yaml
secrets:
  - name: "SERVICE_ACCOUNT_KEY"
    description: "Secret key for the service account"
    destinations:
      - description: "SERVICE_ACCOUNT_KEY for lucasmelin/key-rotator in GitHub Actions"
        type: "github-repository"
        repo: "lucasmelin/key-rotator"
        name: "SERVICE_ACCOUNT_KEY"
      - description: "SERVICE_ACCOUNT_KEY for lucasmelin/key-rotator in Dependabot"
        type: "github-repository-dependabot"
        repo: "lucasmelin/key-rotator"
        name: "SERVICE_ACCOUNT_KEY"
      - description: "SERVICE_ACCOUNT_KEY for lucasmelin/key-rotator in Dependabot"
        type: "github-repository-environment"
        repo: "lucasmelin/key-rotator"
        environment: "prod"
        name: "SERVICE_ACCOUNT_KEY"
  - name: "SERVICE_ACCOUNT_USERNAME"
    description: "Secret username for the service account"
    destinations:
      - description: "SERVICE_ACCOUNT_USERNAME for lucasmelin/key-rotator in GitHub Actions"
        type: "github-repository"
        repo: "lucasmelin/key-rotator"
        name: "SERVICE_ACCOUNT_USERNAME"
```

## Usage

1. Clone the repository and navigate to the project directory.

2. Ensure you have the `GITHUB_TOKEN` environment variable set:

   ```sh
   export GITHUB_TOKEN=your_github_token
   ```

3. Run the program with the path to your YAML configuration file:

   ```sh
   go run . path/to/your/key.yaml
   ```

4. Follow the prompts to enter the secret values.

## License

This project is licensed under the MIT License. See the [`LICENSE` file](./LICENSE) for details.
