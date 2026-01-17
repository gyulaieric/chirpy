# Chirpy

A lightweight web server and social media backend built with Go.

## Description

Chirpy is a Twitter-like API that allows users to create accounts, post short messages (chirps), and interact with content. This project demonstrates RESTful API design, user authentication, and database management in Go.

## Features

- User registration and authentication
- CRUD functionality for chirps
- JWT-based authentication with refresh tokens
- PostgreSQL database integration

## Installation

1. Clone the repository:
```bash
git clone https://github.com/gyulaieric/chirpy.git
cd chirpy
```

2. Install Dependencies:
```bash
go mod download
```

3. Set up your environment variables (create a .env file):
```bash
DB_URL=your_db_url
JWT_SECRET=your_jwt_secret
PLATFORM="dev"
POLKA_KEY=your_polka_api_key
```

4. Run the server:
```bash
go build -o chirpy && ./chirpy
```

# Usage

The server runs on http://localhost:8080 by default. Click [here](/docs/endpoints) for documentation on the available endpoints.
