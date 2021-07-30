-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="Pacific/Tahiti";

CREATE TABLE IF NOT EXISTS pwtoshare (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    pw TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    days_limit INT NOT NULL,
    views_remaining INT NOT NULL
);
