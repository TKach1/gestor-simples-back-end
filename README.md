# Gestor Simples API

An API developed in Go for managing users, products, and sales, providing a comprehensive solution for business administration.

## Features

-   **User Management**: CRUD operations for users (administrators and sellers).
-   **Product Management**: CRUD operations for products, including stock control.
-   **Sales Management**: Record new sales, update product stock, and view sales history.
-   **Authentication & Authorization**: Secure access using JWT (JSON Web Tokens) with role-based authorization (Admin/Seller).
-   **Dashboard Summary**: Provides aggregated data for quick business insights.

## Technologies Used

-   **Go**: The main programming language.
-   **Gorilla Mux**: Powerful URL router and dispatcher.
-   **PostgreSQL**: Relational database for data storage.
-   **JWT (JSON Web Tokens)**: For secure authentication.
-   **Bcrypt**: For secure password hashing.
-   **Docker**: For containerization and easy deployment.
-   **godotenv**: For loading environment variables from a `.env` file.

## Project Structure

-   `main.go`: Entry point of the application, sets up routes and starts the server.
-   `internal/database`: Contains database connection and migration logic.
-   `internal/models`: Defines data structures (structs) for users, products, sales, etc.
-   `pkg/auth`: Handles authentication logic, JWT generation, and password hashing.
-   `API_DOCUMENTATION.md`: Detailed documentation of all API endpoints.
-   `database_diagram.md`: Description and ER diagram of the database schema.
-   `Dockerfile`: Defines the Docker image for the application.

## Getting Started (Local Development)

### Prerequisites

Before running the application locally, ensure you have the following installed:

-   **Go**: Version 1.25.4 or higher.
-   **PostgreSQL**: A running PostgreSQL instance.
-   **make**: (Optional) For running convenient scripts.

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/gestor-simples-lambda.git
cd gestor-simples-lambda
```

### 2. Environment Variables

Create a `.env` file in the root directory of the project based on the `.env.example` (if provided, otherwise create it manually) and configure your database connection and JWT secret.

```
# .env file example
DATABASE_URL="postgres://user:password@host:port/dbname?sslmode=disable"
JWT_SECRET="your_super_secret_jwt_key"
```

### 3. Database Setup

Ensure your PostgreSQL database is running and accessible. You'll need to create the database and run the schema migrations. The database schema is detailed in `database_diagram.md`. You might need a tool like `migrate` or custom SQL scripts to set up your tables.

**Example SQL to create tables (based on `database_diagram.md`):**

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    price REAL NOT NULL DEFAULT 0.0
);

CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sales_items (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER NOT NULL REFERENCES sales(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL
);
```

### 4. Run the Application

```bash
go run main.go
```

The API will start on `http://localhost:8080`.

## Running with Docker

### Prerequisites

-   **Docker**: Ensure Docker is installed and running on your system.

### 1. Build the Docker Image

Navigate to the root directory of the project and build the Docker image:

```bash
docker build -t gestor-simples-api .
```

This command builds an image named `gestor-simples-api`.

### 2. Run the Docker Container

You'll need to link your PostgreSQL database to the Docker container or provide the `DATABASE_URL` via environment variables.

**Example using `docker run`:**

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:password@host.docker.internal:5432/dbname?sslmode=disable" \
  -e JWT_SECRET="your_super_secret_jwt_key" \
  gestor-simples-api
```

**Note**: `host.docker.internal` is a special DNS name that resolves to the internal IP address of the host machine from within a Docker container. Use this if your PostgreSQL is running directly on your host machine. If your database is in another Docker container, you might need to use Docker Compose or link containers.

The API will be accessible on `http://localhost:8080` on your host machine.

## API Documentation

For detailed information on all available API endpoints, request/response formats, and authentication requirements, please refer to:

-   [`API_DOCUMENTATION.md`](./API_DOCUMENTATION.md)

## Database Diagram

The relational database schema is described and visualized in:

-   [`database_diagram.md`](./database_diagram.md)

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.
