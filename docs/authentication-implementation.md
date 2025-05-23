# Authentication Implementation for Pet Paradise

This document outlines the authentication system implemented for the Pet Paradise application.

## Overview

The authentication system uses JSON Web Tokens (JWT) to authenticate users and protect API endpoints. It includes:

1. Token generation during login
2. Token validation middleware
3. Role-based access control for protected routes

## JWT Authentication

### Token Generation

When a user logs in with valid credentials, the system:
1. Validates the user's email and password
2. Creates a JWT token containing the user's ID, email, and role
3. Sets an expiration time (24 hours by default)
4. Signs the token with a secret key
5. Returns the token to the client

### Token Structure

The JWT payload contains the following claims:

```json
{
  "user_id": "user-uuid",
  "email": "user@example.com",
  "role": "admin",
  "exp": 1621234567,
  "iat": 1621148167,
  "sub": "user-uuid"
}
```

### Token Validation

For protected routes, the middleware:
1. Extracts the JWT token from the Authorization header
2. Verifies the token signature
3. Validates the token expiration
4. Extracts user information and adds it to the request context

## Middleware Implementation

Two middleware components are implemented:

1. **Protected Middleware**: Validates the JWT token and provides access to authenticated users
2. **Role-Required Middleware**: Checks if the authenticated user has the required role

## Route Protection

### Public Routes (No Authentication Required)

- `POST /api/users/register` - User registration
- `POST /api/users/login` - User login
- `GET /api/pets` - Get all pets
- `GET /api/pets/:id` - Get pet details
- `GET /api/pets/status` - Get pets by status

### Protected Routes (Authentication Required)

#### User Routes
- `GET /api/users/:id` - Get user details
- `PUT /api/users/:id` - Update user information
- `POST /api/users/:id/password` - Change password
- `POST /api/users/logout` - Logout

#### Pets Routes
- `POST /api/pets` - Create a pet (staff only)
- `PUT /api/pets/:id` - Update pet details (staff only)
- `PATCH /api/pets/:id/status` - Update pet status (staff only)
- `DELETE /api/pets/:id` - Delete a pet (staff only)

#### Adoptions Routes
- `POST /api/adoptions` - Create an adoption request
- `GET /api/adoptions/:id` - Get adoption details
- `GET /api/adoptions/user/:userId` - Get user's adoptions
- `GET /api/adoptions` - Get all adoptions (staff only)
- `PUT /api/adoptions/:id` - Update adoption details (staff only)
- `DELETE /api/adoptions/:id` - Delete an adoption (staff only)

#### Donations Routes
- `POST /api/donations` - Make a donation
- `GET /api/donations/user/:userId` - Get user's donations
- `GET /api/donations/:id` - Get donation details
- `GET /api/donations` - Get all donations (admin only)
- `PATCH /api/donations/:id/status` - Update donation status (admin only)
- `DELETE /api/donations/:id` - Delete a donation (admin only)

## Role-Based Access Control

The system implements the following user roles with different permission levels:

1. **Admin** (`admin`): Full access to all system features
2. **User** (`user`): Access to personal data and basic functionality
3. **Volunteer** (`volunteer`): Manage pets and adoptions
4. **Vet** (`vet`): Manage pets, medical records, and adoptions

## Usage Examples

### Authentication Header

The client must include the token in the Authorization header for protected routes:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Login Response

```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "john@example.com",
    "status": "active",
    "role": "admin"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Security Considerations

1. **Secret Key Management**: In production, the JWT secret key should be stored in environment variables or a secure key management system.
2. **Token Expiration**: JWT tokens expire after 24 hours, requiring users to log in again.
3. **HTTPS**: All API communication should be over HTTPS to prevent token interception.
4. **Token Storage**: Clients should store the token securely (e.g., in HttpOnly cookies or secure browser storage).

## Future Improvements

1. **Refresh Tokens**: Implement refresh tokens for extending sessions without requiring re-login.
2. **Token Revocation**: Add a mechanism to revoke tokens before they expire (e.g., during logout).
3. **Claims-Based Authorization**: Expand role-based access to more granular permission claims.
4. **Two-Factor Authentication**: Add support for 2FA for sensitive operations.
5. **Rate Limiting**: Implement API rate limiting to prevent brute force attacks.
