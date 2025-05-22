# Football League Simulator

A Go-based REST API for simulating a football league with 4 teams.

## 🎯 Project Overview

This project simulates a football league consisting of 4 teams. The system simulates matches between teams, updates the league standings after each match, and predicts the final standings after all weeks are played. Match results are determined based on the relative strength of each team.

## 🛠️ Technologies Used

- **Programming Language**: Go (1.24)
- **Database**: PostgreSQL
- **Testing/Client**: Postman
- **API Framework**: Fiber
- **Documentation**: Swagger

## 📦 Architecture

- **MVC Pattern**: Uses Model-View-Controller architecture for clean separation of concerns
- **Interface-based Design**: Implements interfaces for key components to ensure loose coupling
- **Struct Composition**: Creates well-defined structures for entities like Team, Match, League, etc.
- **REST API Endpoints**: Provides a RESTful API to interact with the simulator
- **Swagger Documentation**: API endpoints are documented with Swagger for easy testing and integration

## 📁 Project Structure

```
league-simulator/
├── main.go                 # Application entry point
├── src/
│   ├── config/             # Configuration handling
│   ├── controller/         # HTTP request handlers
│   ├── database/           # Database connection and migrations
│   │   ├── migrations/     # SQL schema definitions
│   │   └── seed/           # Initial data for the database
│   ├── docs/               # Swagger documentation
│   ├── middleware/         # HTTP middleware
│   ├── model/              # Data models
│   ├── repository/         # Data access layer
│   └── service/            # Business logic
```

## 🔄 Core Functionality

1. Create and manage 4 teams with different strength attributes
2. Schedule and simulate matches between teams
3. Calculate and update league standings (points, wins, draws, losses)
4. Generate match results based on team strengths
5. Predict final standings after all weeks
6. Provide API endpoints to access all functionality

## 📊 Data Model

- **Team**: ID, Name, Strength, Points, etc.
- **Match**: ID, HomeTeam, AwayTeam, HomeScore, AwayScore, Week, etc.
- **League**: ID, Name, Teams, Matches, CurrentWeek, etc.
- **Standings**: Teams, Points, Wins, Draws, Losses, etc.

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Docker (optional)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/league-simulator.git
   cd league-simulator
   ```

2. Set up environment variables (or create a `.env` file):

   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=league_simulator
   DB_SSLMODE=disable
   SERVER_PORT=8080
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Generate Swagger documentation:

   ```bash
   # Install swag CLI tool if you don't have it
   go install github.com/swaggo/swag/cmd/swag@latest

   # Generate Swagger docs
   cd src
   swag init -g ../main.go -o ./docs
   cd ..
   ```

5. Run the application:

   ```bash
   go run main.go
   ```

6. Access the Swagger documentation:
   ```
   http://localhost:8080/swagger/
   ```

## 🔌 API Endpoints

All endpoints are available under both `/api` prefix and root path for backward compatibility.

### Teams

- `GET /api/teams` - List all teams
- `GET /api/teams/{id}` - Get a specific team
- `POST /api/teams` - Create a new team
- `PUT /api/teams/{id}` - Update a team
- `DELETE /api/teams/{id}` - Delete a team
- `POST /api/teams/initialize` - Create initial 4 teams

### Matches

- `GET /api/matches` - List all matches
- `GET /api/matches?week={week}` - List matches for a specific week
- `GET /api/matches/{id}` - Get a specific match
- `POST /api/matches` - Create a new match
- `PUT /api/matches/{id}` - Update a match

### League

- `POST /api/leagues` - Create a new league
- `GET /api/leagues/{id}` - Get a specific league
- `POST /api/leagues/{id}/simulate` - Simulate matches for the next week
- `GET /api/leagues/{id}/standings` - Get current standings

### Prediction

- `GET /api/leagues/{id}/predict` - Predict final standings

### Swagger Documentation

- `GET /swagger/` - Interactive API documentation

## 🧪 Testing

You can test the API using Swagger UI, Postman, or curl. Here's a sample workflow:

### Using Swagger UI

1. Open the Swagger documentation at `http://localhost:8080/swagger/`
2. Use the interactive UI to test each endpoint
3. Try the following workflow:
   - Create teams with the `/api/teams/initialize` endpoint
   - Create a league with the `/api/leagues` endpoint
   - Simulate matches with the `/api/leagues/{id}/simulate` endpoint
   - View standings with the `/api/leagues/{id}/standings` endpoint
   - Predict final standings with the `/api/leagues/{id}/predict` endpoint

### Using curl

1. Initialize teams:

   ```bash
   curl -X POST http://localhost:8080/api/teams/initialize
   ```

2. Create a league:

   ```bash
   curl -X POST http://localhost:8080/api/leagues -H "Content-Type: application/json" -d '{"name":"Premier League"}'
   ```

3. Simulate a week:

   ```bash
   curl -X POST http://localhost:8080/api/leagues/1/simulate
   ```

4. Get current standings:

   ```bash
   curl -X GET http://localhost:8080/api/leagues/1/standings
   ```

5. Predict final standings:
   ```bash
   curl -X GET http://localhost:8080/api/leagues/1/predict
   ```

## 🔧 Development

### Regenerating Swagger Documentation

If you make changes to the API endpoints or models, you'll need to regenerate the Swagger documentation:

```bash
cd src
swag init -g ../main.go -o ./docs
cd ..
```

### Database Migrations

The application automatically runs migrations when it starts. The migration files are located in `src/database/migrations/`.

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.
