-- Create User table
CREATE TABLE IF NOT EXISTS "User" (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Match table
CREATE TABLE IF NOT EXISTS match (
    match_id SERIAL PRIMARY KEY,
    white_player_username VARCHAR(255) NOT NULL REFERENCES "User"(username),
    black_player_username VARCHAR(255) NOT NULL REFERENCES "User"(username),
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    winner_id INT,
    CONSTRAINT check_players_different CHECK (white_player_username != black_player_username)
);

CREATE TABLE IF NOT EXISTS queue (
    queue_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL REFERENCES "User"(username),
    join_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS status (
    status_id SERIAL PRIMARY KEY,
    match_id INT REFERENCES "match"(match_id),
    username VARCHAR(255) REFERENCES "User"(username),
    opponent VARCHAR(255) REFERENCES "User"(username),
    isPlaying BOOLEAN
);