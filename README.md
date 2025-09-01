# Chirpy - A Twitter-like Backend Service

Chirpy is a robust backend service written in Go that implements core functionality similar to Twitter. It provides a RESTful API for managing users, posts (called "chirps"), and user authentication.

## Features

- 🔐 User Authentication
  - JWT-based authentication
  - Refresh token support
  - Password hashing for security
  - Email/password update capability
  
- 📝 Chirps Management
  - Create and delete chirps
  - Get individual chirps
  - List all chirps with sorting options
  - Filter chirps by user
  
- 💎 Premium Features
  - Premium membership support
  - Webhook integration for membership upgrades
  - API key validation for webhooks
  
- 🔍 Additional Features
  - Health check endpoint
  - Metrics tracking
  - File serving capabilities
  - Database reset functionality

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Environment variables setup (see Configuration section)

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
DB_URL=postgresql://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your-jwt-secret-key
PLATFORM=your-platform-name
POLKA_KEY=your-polka-webhook-key
```

## API Endpoints

### Authentication & User Management
- `POST /api/users` - Create new user (sign up)
- `POST /api/login` - User login (returns access token)
- `POST /api/refresh` - Refresh access token using refresh token
- `POST /api/revoke` - Revoke refresh token
- `PUT /api/users` - Update user details (email/password)

### Chirps
- `POST /api/chirps` - Create a new chirp
- `GET /api/chirps` - Get all chirps with optional parameters:
  - `?author_id={userID}` - Filter chirps by user
  - `?sort={sortingMethod}` - Sort by creation date ("asc" or "desc")
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete a chirp (requires authentication)

### Premium Features
- `POST /api/polka/webhooks` - Webhook endpoint for premium membership upgrades (requires Polka API key)

### System & Metrics
- `GET /api/healthz` - Health check endpoint
- `GET /admin/metrics` - Get visitor metrics
- `POST /admin/reset` - Reset visitor counter
- `/app/*` - Static file server (with metrics tracking)

## Database Schema

The application uses PostgreSQL with the following main tables:

- `users` - Stores user information
- `chirps` - Stores all chirps
- `refresh_tokens` - Manages refresh tokens

Database migrations are handled using Goose.

## Security Features

- Password hashing using secure algorithms
- JWT-based authentication
- API key validation for webhooks
- Refresh token rotation
- Input validation and sanitization

## Running the Application

1. Set up your environment variables
2. Ensure PostgreSQL is running
3. Run the migrations
4. Build and start the server:
   ```bash
   # Build the executable
   go build .

   # Run the executable
   ./chirpy
   ```
   The server will start on port 8080 by default.

## Development

The project follows a clean architecture pattern with the following structure:

- `/internal` - Core business logic and internal packages
- `/sql` - Database queries and schema migrations
- Handlers are separated by functionality into different files
- Authentication logic is isolated in the auth package

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 