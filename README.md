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

### Authentication
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Users
- `POST /api/users` - Create new user
- `PUT /api/users` - Update user details

### Chirps
- `DELETE /api/chirps/{chirpID}` - Delete a chirp
- `POST /api/chirps` - Create a new chirp
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `GET /api/chirps?author_id={userID}` - Get chirps from specific user
- `GET /api/chirps?sort={sortingMethod}` - Get chirps sorted by their creation date. 
(sorting methods: "asc" for ascending, "desc" for descending)

### Premium Features
- `POST /api/polka/webhooks` - Webhook endpoint for premium membership upgrades

### System
- `GET /api/healthz` - Health check endpoint
- `GET /admin/metrics` - Get metrics
- `POST /admin/reset` - Reset metrics counter

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