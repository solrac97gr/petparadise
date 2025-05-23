# Pet Paradise - An Open Source Pet Refuge Management System

> WORK IN PROGRESS 

Pet Paradise is a microservices-based application for pet adoption and donation management, designed to help manage pet refuges. It provides a user-friendly interface for handling pet information, adoptions, and donations.

## Project Structure

The project follows a clean hexagonal architecture approach:

```
backend/
  ├── cmd/               # Application entry points
  ├── internal/          # Application core modules
  │   ├── adoptions/     # Adoption module
  │   ├── donations/     # Donations module
  │   ├── pets/          # Pets module
  │   └── users/         # Users module
  ├── pkg/               # Shared packages
  └── deploy/            # Deployment configurations

frontend/
  ├── src/               # Frontend source code
  └── deploy/            # Frontend deployment configurations
```

Each module in the internal directory follows the hexagonal architecture pattern:
- **Domain**: Contains the business logic, models, and interfaces (ports)
- **Application**: Contains the use cases and service implementations
- **Infrastructure**: Contains the external adapters (REST API, repositories)

## Features

- User authentication and authorization
- Pet management (add, edit, delete pets)
- Adoption management (view, approve, reject adoptions)
- Donation management (view, add, delete donations)

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.24 or later
- Node.js 18 or later
- PostgreSQL (for local development without Docker)

### Development

1. Clone the repository:
   ```bash
   git clone https://github.com/solrac97gr/petparadise.git
   cd petparadise
   ```

2. Start all services with Docker Compose:
   ```bash
   make run
   ```

3. Or run services individually:
   ```bash
   # Start just the database
   make db-start
   
   # Run the backend
   make backend-run
   
   # Run the frontend
   make frontend-run
   ```

4. Stop all services:
   ```bash
   make stop
   ```

### API Endpoints

#### Adoptions
- `GET /api/adoptions` - Get all adoptions
- `GET /api/adoptions/:id` - Get adoption by ID
- `POST /api/adoptions` - Create a new adoption
- `PUT /api/adoptions/:id` - Update an adoption
- `DELETE /api/adoptions/:id` - Delete an adoption
- `GET /api/adoptions/user/:userId` - Get adoptions by user ID

#### Users, Pets, and Donations
- Similar endpoint structures for each module

## License

MIT
