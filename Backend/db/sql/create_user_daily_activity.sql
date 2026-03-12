PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS "UserDailyActivity" (
    user_id INTEGER NOT NULL,
    activity_date DATE NOT NULL,
    enrolled INTEGER NOT NULL DEFAULT 0,
    unenrolled INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, activity_date),
    FOREIGN KEY (user_id) REFERENCES "User"(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_daily_activity_user_date
    ON "UserDailyActivity" (user_id, activity_date);

COMMIT;
