-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    strength INTEGER NOT NULL CHECK (strength >= 1 AND strength <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    home_score INTEGER DEFAULT 0,
    away_score INTEGER DEFAULT 0,
    week INTEGER NOT NULL,
    played BOOLEAN DEFAULT FALSE,
    played_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (home_team_id != away_team_id)
);

-- Create standings_history table
CREATE TABLE IF NOT EXISTS standings_history (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id),
    week INTEGER NOT NULL,
    points INTEGER DEFAULT 0,
    played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    draws INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    goals_for INTEGER DEFAULT 0,
    goals_against INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (team_id, week)
);

-- Create leagues table
CREATE TABLE IF NOT EXISTS leagues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    current_week INTEGER DEFAULT 0,
    total_weeks INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create function to update timestamps
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS update_teams_timestamp ON teams;
DROP TRIGGER IF EXISTS update_matches_timestamp ON matches;
DROP TRIGGER IF EXISTS update_standings_history_timestamp ON standings_history;
DROP TRIGGER IF EXISTS update_leagues_timestamp ON leagues;

-- Create triggers for updated_at columns
CREATE TRIGGER update_teams_timestamp
BEFORE UPDATE ON teams
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_matches_timestamp
BEFORE UPDATE ON matches
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_standings_history_timestamp
BEFORE UPDATE ON standings_history
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_leagues_timestamp
BEFORE UPDATE ON leagues
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();
