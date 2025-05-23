# Donations Implementation

This document outlines the implementation details of the Donations module in the Pet Paradise system.

## Domain Models

### Donation

The main entity representing a donation in the system:

- `ID` - Unique identifier for the donation
- `UserID` - ID of the user who made the donation
- `Amount` - The monetary amount of the donation
- `Status` - Current status of the donation (pending, completed, failed, refunded)
- `Created` - When the donation was created
- `Updated` - When the donation was last updated
- `Comment` - Optional comment from the donor
- `Anonymous` - Whether the donation should be shown as anonymous

### Status

An enumeration representing the possible statuses of a donation:
- `pending` - Donation has been initiated but not processed
- `completed` - Donation has been successfully processed
- `failed` - Donation processing failed
- `refunded` - Donation was refunded to the donor

## Architecture

The Donations module follows the hexagonal architecture pattern:

### Domain Layer

- Contains the core business logic and models
- Defines interfaces (ports) for external dependencies

### Application Layer

- Implements use cases using domain models
- Coordinates between domain and infrastructure layers
- Contains business logic orchestration

### Infrastructure Layer

#### Repository

- PostgreSQL implementation for storing and retrieving donations
- Handles database-specific operations

#### API

- HTTP handlers for donation-related endpoints
- Request validation and response formatting

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/donations | Create a new donation |
| GET | /api/donations | Get all donations |
| GET | /api/donations/:id | Get a specific donation by ID |
| GET | /api/donations/user/:userId | Get all donations for a specific user |
| PATCH | /api/donations/:id/status | Update a donation's status |
| DELETE | /api/donations/:id | Delete a donation |

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS donations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,
    comment TEXT,
    anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_donations_user_id ON donations(user_id);
CREATE INDEX IF NOT EXISTS idx_donations_status ON donations(status);
CREATE INDEX IF NOT EXISTS idx_donations_created ON donations(created);
```
