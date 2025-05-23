# Authentication Implementation for Pet Paradise

This document outlines the authentication system implemented for the Pet Paradise application.

## Overview

The authentication system uses JSON Web Tokens (JWT) to authenticate users and protect API endpoints. It includes:

1. Token generation during login
2. Token validation middleware
3. Role-based access control for protected routes
4. Token refresh mechanism
5. Token revocation system

## JWT Authentication

### Token Generation

When a user logs in with valid credentials, the system:
1. Validates the user's email and password
2. Creates an access token (short-lived) and refresh token (long-lived)
3. Sets appropriate expiration times (15 minutes for access token, 7 days for refresh token)
4. Signs both tokens with a secret key
5. Returns both tokens to the client

### Token Structure

The access token payload contains the following claims:

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

The refresh token payload contains:

```json
{
  "user_id": "user-uuid",
  "token_id": "unique-token-identifier",
  "exp": 1621754567,
  "iat": 1621148167,
  "sub": "user-uuid"
}
```

### Token Validation

For protected routes, the middleware:
1. Extracts the access token from the Authorization header
2. Checks if the token has been revoked (using an in-memory blacklist)
3. Verifies the token signature
4. Validates the token expiration
5. Extracts user information and adds it to the request context

### Token Refresh

When the access token expires, the client can:
1. Send the refresh token to the `/api/users/refresh` endpoint
2. Receive a new pair of access and refresh tokens
3. The old refresh token is automatically revoked
4. This process continues until the refresh token expires or is revoked

During token refresh, the system:
1. Validates the refresh token
2. Checks if the user still exists and is active
3. Generates a new token pair
4. Revokes the old refresh token
5. Returns the new token pair to the client

Example request to refresh tokens:
```json
POST /api/users/refresh
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Example response:
```json
{
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

### Token Revocation

The system supports revoking tokens before they expire:
1. When a user logs out, their access and refresh tokens are added to a blacklist
2. The blacklist is checked during token validation
3. Revoked tokens are rejected even if they haven't expired yet
4. Expired tokens are automatically removed from the blacklist through periodic cleanup

#### Token Blacklist System

The token blacklist is implemented with the following features:
1. In-memory storage with thread-safe access
2. Automatic cleanup of expired tokens every 15 minutes
3. Tokens are stored with their expiration time to limit the size of the blacklist
4. Logging of token revocation and cleanup activities

#### Revoking All User Tokens

The API supports revoking all tokens for a specific user with:
```
POST /api/users/:id/revoke-tokens
```

This is useful in the following scenarios:
- When a user changes password
- When an account may be compromised
- When an administrator needs to force re-authentication
- When a user's account is suspended or deactivated

## Middleware Implementation

Two middleware components are implemented:

1. **Protected Middleware**: Validates the JWT token and provides access to authenticated users
2. **Role-Required Middleware**: Checks if the authenticated user has the required role

## Route Protection

### Public Routes (No Authentication Required)

- `POST /api/users/register` - User registration
- `POST /api/users/login` - User login
- `POST /api/users/refresh` - Refresh tokens
- `GET /api/pets` - Get all pets
- `GET /api/pets/:id` - Get pet details
- `GET /api/pets/status` - Get pets by status

### Protected Routes (Authentication Required)

#### User Routes
- `GET /api/users/:id` - Get user details
- `PUT /api/users/:id` - Update user information
- `POST /api/users/:id/password` - Change password
- `POST /api/users/logout` - Logout
- `POST /api/users/:id/revoke-tokens` - Revoke all tokens for a user

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
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

## Security Considerations

1. **Secret Key Management**: In production, the JWT secret key should be stored in environment variables or a secure key management system.
2. **Token Expiration**: Access tokens expire after 15 minutes, refresh tokens after 7 days.
3. **Token Revocation**: Both access and refresh tokens can be revoked before expiration.
4. **HTTPS**: All API communication should be over HTTPS to prevent token interception.
5. **Token Storage**: Clients should store access tokens in memory and refresh tokens in secure storage.

## Future Improvements

1. **Persistent Token Storage**: Move blacklist to Redis or database for persistence across server restarts.
2. **Claims-Based Authorization**: Expand role-based access to more granular permission claims.
3. **Two-Factor Authentication**: Add support for 2FA for sensitive operations.
4. **Device Management**: Track tokens by device and allow users to manage active sessions.
5. **Rate Limiting**: Implement API rate limiting to prevent brute force attacks.
6. **Token Introspection**: Add an endpoint for clients to check if a token is still valid.
7. **User Activity Tracking**: Track when tokens are created, used, and revoked for audit purposes.
