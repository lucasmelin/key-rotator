# Key Rotator 🔑

Key Rotator is a Go program designed to manage and rotate secrets for GitHub repositories. It reads a YAML configuration file specifying multiple secrets and their destinations, prompts the user for the secret values, encrypts them using the repository's public key, and updates the secrets in the specified GitHub repositories.

## Prerequisites

- Go installed. See the [official installation guide](https://golang.org/doc/install) for instructions.
- A GitHub personal access token with the necessary permissions to update secrets.
- The `GITHUB_TOKEN` environment variable set with your GitHub token.

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

1. Clone the repository and navigate to the project directory.

2. Ensure you have the `GITHUB_TOKEN` environment variable set:

   ```sh
   export GITHUB_TOKEN=your_github_token
   ```

3. Run the program with the path to your YAML configuration file:

   ```sh
   go run . path/to/your/key.yaml
   # Or optionally, use the dry-run flag to see what
   # changes would be made without updating the secrets.
   go run . --dry-run path/to/your/key.yaml
   ```

4. Follow the prompts to enter the secret values.

## License

This project is licensed under the MIT License. See the [`LICENSE` file](./LICENSE) for details.
