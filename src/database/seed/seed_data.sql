-- Seed data for teams
INSERT INTO teams (name, strength) VALUES
    ('Manchester United', 85),
    ('Liverpool', 88),
    ('Chelsea', 82),
    ('Arsenal', 80)
ON CONFLICT (id) DO NOTHING;

-- Seed data for a league
INSERT INTO leagues (name, current_week, total_weeks)
VALUES ('Premier League', 0, 3)
ON CONFLICT (id) DO NOTHING;

-- Seed data for matches (round-robin tournament for 4 teams)
-- Week 1
INSERT INTO matches (home_team_id, away_team_id, week, played)
VALUES 
    (1, 2, 1, false),
    (3, 4, 1, false)
ON CONFLICT (id) DO NOTHING;

-- Week 2
INSERT INTO matches (home_team_id, away_team_id, week, played)
VALUES 
    (1, 3, 2, false),
    (2, 4, 2, false)
ON CONFLICT (id) DO NOTHING;

-- Week 3
INSERT INTO matches (home_team_id, away_team_id, week, played)
VALUES 
    (1, 4, 3, false),
    (2, 3, 3, false)
ON CONFLICT (id) DO NOTHING;

-- Initialize standings for each team at week 0
INSERT INTO standings_history (team_id, week, points, played, wins, draws, losses, goals_for, goals_against)
VALUES
    (1, 0, 0, 0, 0, 0, 0, 0, 0),
    (2, 0, 0, 0, 0, 0, 0, 0, 0),
    (3, 0, 0, 0, 0, 0, 0, 0, 0),
    (4, 0, 0, 0, 0, 0, 0, 0, 0)
ON CONFLICT (team_id, week) DO NOTHING;
