CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_reservations_room_id ON reservations(room_id);
CREATE INDEX IF NOT EXISTS idx_reservations_time ON reservations(room_id, start_time, end_time);