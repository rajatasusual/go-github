# Go-Github

![Go](https://img.shields.io/badge/Go-1.23-blue) ![Docker](https://img.shields.io/badge/Docker-20.10-blue) ![Kubernetes](https://img.shields.io/badge/Kubernetes-1.21-blue)

![Go-Github](https://raw.githubusercontent.com/rajatasusual/go-github/refs/heads/main/favicon.ico)

**Go-Github** is a Golang-based project that retrieves and displays GitHub user profiles, commit histories, and statistics. It combines templated HTML with Chart.js for visualizations and supports multiple database connectors: a local SQLite database and a managed Xata database. The project is containerized for easy deployment in Kubernetes environments.

## Features

- **Golang-based Backend**: Written entirely in Go, with a clean architecture.
- **Templated HTML with Templ**: Uses `templ` to render and serve HTML content dynamically.
- **Commit History Visualization**: Displays commit history using Chart.js for user-friendly data visualization.
- **Database Storage**:
  - **SQLite**: Local lightweight database, perfect for isolated environments and testing.
  - **Xata**: Managed database connector, allowing you to switch to a hosted solution.
- **GitHub SDK Integration**: Fetches user data and repositories directly from GitHub.
- **Containerization**: Easily deployable Docker container.
- **Kubernetes Deployment**: Configured for smooth Kubernetes deployment with automated builds.
- **File Change Watcher**: `watcher.sh` script monitors code changes, automatically rebuilding the container for rapid development.

## Project Structure

- **`cmd`**: Contains handlers and routes to handle HTTP requests and route them appropriately.
- **`helper`**: Holds database connectors for Xata and SQLite, managing database setup, connection pooling, and CRUD functions.
- **`service`**: Acts as a controller for database operations, processing and directing logic between handlers and database helpers.
- **`views`**: Includes templated HTML files served by the application, presenting the frontend for users.

## Installation

1. **Clone the Repository**

    ```bash
    git clone https://github.com/rajatasusual/go-github.git
    cd go-github
    ```

2. **Install Dependencies**

    Ensure Go is installed (Go 1.18 or later) and run:

    ```bash
    go mod tidy
    ```

3. **Environment Configuration**

    Create a `.env` file with the following environment variables:

    ```dotenv
    XATA_DATABASE_NAME=go-github
    XATA_DATABASE_URL=<Xata database URL>
    XATA_TABLE_NAME=user-details
    XATA_API_KEY=<Your Xata API key>
    GITHUB_TOKEN=<Your GitHub token>
    SQLLITE_DB_PATH=./app.db
    ```

4. **Build and Run Locally**

    ```bash
    go build -o go-github ./cmd
    ./go-github
    ```

5. **Docker Deployment**

    Pull the latest Docker image:

    ```bash
    docker pull rajatasusual/main:2.0.0
    ```

    Run the Docker container:

    ```bash
    docker run -d --env-file .env -p 8080:8080 rajatasusual/main:2.0.0
    ```

6. **Kubernetes Deployment**

    Make sure your Kubernetes cluster is running, then create the deployment using your preferred Kubernetes YAML configuration or Helm chart. Here’s a basic example:

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: go-github
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: go-github
      template:
        metadata:
          labels:
            app: go-github
        spec:
          containers:
            - name: go-github
              image: rajatasusual/main:2.0.0
              envFrom:
                - configMapRef:
                    name: go-github-config
              ports:
                - containerPort: 8080
    ```

    Create a `ConfigMap` for environment variables or use a `.env` file when deploying to Kubernetes.

## Configuration

The application uses environment variables for configuration, defined in a `.env` file:

```dotenv
XATA_DATABASE_NAME=go-github           # Database name for Xata
XATA_DATABASE_URL=<URL>                 # Xata database URL
XATA_TABLE_NAME=user-details            # Table name for user data
XATA_API_KEY=<Xata API key>             # Xata API key
GITHUB_TOKEN=<GitHub API token>         # Token for GitHub API access
SQLLITE_DB_PATH=./app.db                # Path for local SQLite database
```

## Contributing

Contributions are welcome! To get started:

	1.	Fork the repository.
	2.	Create a new branch for your feature: git checkout -b feature-name.
	3.	Make your changes and commit them with clear messages.
	4.	Submit a pull request.

## License

This project is licensed under the MIT License.

## Acknowledgements

	•	Chart.js for easy-to-use data visualizations.
	•	Xata for managed database services.
	•	Go SDK for GitHub for GitHub data access.
	•	templ for fast, secure templating in Go.

Enjoy exploring GitHub profiles and contributions with Go-Github! Reach out for support at Rajatasusual.
