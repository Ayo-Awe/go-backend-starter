# Go Backend Starter Template

This project serves as a starting point for my future backend side-projects in Go. It comes with a fairly well-organized structure, including services, concrete repository implementations, custom application errors, API boilerplate, and integration tests. A lot of things in this setup are still experimental and are likely to change overtime. If you happen to recognise any pattern that you believe might not scale well in production, feel free to open am issue and we can discuss it there.

## Prerequisites

Before getting started, ensure you have the following installed:

- **Go:** Version 1.23.1 or higher ([download](https://golang.org/dl/))
- **Task:** A command runner for automating tasks ([installation instructions](https://taskfile.dev/#/usage))
- **Docker:** For containerizing your environment ([download](https://www.docker.com/get-started))

## Project Structure

- **`cmd/`**  
  Contains the application entry points.

  - The [cmd/server/main.go](cmd/server/main.go) file starts the server.
  - Other command folders (e.g., [cmd/goose/main.go](cmd/goose/main.go)) handle tasks like database migrations.

- **`internal/`**

  - **`app/`**  
    Contains different services such as authentication ([internal/app/auth](internal/app/auth)) and user management ([internal/app/users](internal/app/users)).
  - **`db/`**  
    Provides the concrete repository implementations for interacting with the database. Check out the implementations in [internal/db/database.go](internal/db/database.go) and related files.
  - **`apperrors/`**  
    Contains custom application errors ([internal/apperrors/error.go](internal/apperrors/error.go)).
  - **`domain/`**  
    Holds foundational models and interfaces, including repositories and database interfaces, and paging logic.
  - **`oapi/`**  
    Contains generated boilerplate code from the [openapi.yml](openapi.yml) specification, ensuring consistency between your API and code.
  - **`server/`**  
    Implements HTTP handlers for incoming requests. Handlers often import request and response definitions from the `oapi` package and facilitate business logic.
  - Other folders such as `mail`, `rbac`, `vcs`, and `workers` provide additional supporting functionality.

- **`tests/`**  
  Includes integration tests for the project. The test environment setup is detailed in [tests/integration/testenv/testenv.go](tests/integration/testenv/testenv.go).

- **Other Files**
  - **`docker-compose.yml`** and Dockerfiles ([api.Dockerfile](api.Dockerfile), [docs.Dockerfile](docs.Dockerfile)) help configure and run containers.
  - **`Taskfile.yml`** defines tasks for code generation, migrations, and running the application.
  - **`.env`** and **`.env.local`** contain environment variables. Copy `.env.local` to `.env` to configure your local environment.

## Getting Started Locally

To run the project locally, follow these steps:

1. **Set Up Environment Variables:**  
   Copy over the `.env.local` file as `.env`:

   ```sh
   cp .env.local .env
   ```

2. **Start Docker Containers:**
   Bring up the necessary containers with:
   ```sh
   task docker/up
   ```
3. **Run the Application:**
   Finally, start the server:
   ```sh
   task run
   ```

## Code Generation and Migrations

- **API Generation:**  
  The Taskfile includes tasks to generate API models and boilerplate code from your OpenAPI definitions.

  ```sh
  task generate/api
  ```

- **Database Code Generation:**  
  Generate SQLC code for database queries and models:

  ```sh
  task generate/db
  ```

- **Database Migrations:**  
   Run migrations using Goose:
  ```sh
  task migrations/up
  ```

## WIP

There are still a couple of features that I'd like to implement and will do over the next few months

- Background workers
- Auth setup (OAuth and Session tokens)
- Automated CI/CD pipeline setup using Dagger

## Contributing

Feel free to fork this repository and submit pull requests. For any issues or questions, please open an issue in the repository.
