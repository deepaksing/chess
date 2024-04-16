-- Create User table
CREATE TABLE IF NOT EXISTS "User" (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Match table
CREATE TABLE IF NOT EXISTS "Match" (
    match_id SERIAL PRIMARY KEY,
    white_player_id INT REFERENCES "User"(user_id),
    black_player_id INT REFERENCES "User"(user_id),
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    winner_id INT,
    CONSTRAINT check_players_different CHECK (white_player_id != black_player_id)
);

CREATE TABLE IF NOT EXISTS Queue (
    queue_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL REFERENCES "User"(username),
    join_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
