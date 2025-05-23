# Users Implementation

This document outlines the implementation details of the Users module in the Pet Paradise system.

## Domain Models

### User

The main entity representing a user in the system:

- `ID` - Unique identifier for the user
- `Name` - User's full name
- `Email` - User's email address (unique)
- `Password` - Hashed password for authentication
- `Status` - Current status of the user (active, inactive, suspended, pending)
- `Created` - When the user was created
- `Updated` - When the user was last updated
- `Role` - User's role in the system (admin, user, volunteer, vet)
- `Address` - User's physical address
- `Phone` - User's phone number
- `Documents` - Array of document references/links

### Status

An enumeration representing the possible statuses of a user:
- `active` - User is active and can use the system
- `inactive` - User has been deactivated
- `suspended` - User has been temporarily suspended
- `pending` - User registration is pending approval

### Role

An enumeration representing the possible roles of a user:
- `admin` - System administrator with full access
- `user` - Regular user with limited access
- `volunteer` - Volunteer with special access to certain features
- `vet` - Veterinarian with access to medical features

## Architecture

The Users module follows the hexagonal architecture pattern:

### Domain Layer

- Contains the core business logic and models
- Defines interfaces (ports) for external dependencies

### Application Layer

- Implements use cases using domain models
- Coordinates between domain and infrastructure layers
- Contains business logic orchestration

### Infrastructure Layer

#### Repository

- PostgreSQL implementation for storing and retrieving users
- Handles database-specific operations

#### API

- HTTP handlers for user-related endpoints
- Request validation and response formatting
- Authentication and authorization handling

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/users | Create a new user |
| GET | /api/users | Get all users |
| GET | /api/users/:id | Get a specific user by ID |
| GET | /api/users/email | Get a user by email |
| GET | /api/users/status | Get users by status |
| PUT | /api/users/:id | Update a user's information |
| PATCH | /api/users/:id/role | Update a user's role |
| PATCH | /api/users/:id/status | Update a user's status |
| POST | /api/users/:id/password | Change a user's password |
| DELETE | /api/users/:id | Delete a user |
| POST | /api/users/login | Authenticate a user |
| POST | /api/users/logout | Log out a user |

## Security

- Passwords are hashed using bcrypt before storage
- Authentication is handled through JWT tokens (placeholder implementation)
- User statuses are used to control access (only active users can log in)

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,
    role VARCHAR(20) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    documents JSONB
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
```

## Future Improvements

- Implement proper JWT token generation and validation
- Add refresh token mechanism
- Implement password reset functionality
- Add email verification
- Enhance validation for user inputs (email format, password strength)
