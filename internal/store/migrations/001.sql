CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     email TEXT UNIQUE NOT NULL,
                                     pass_hash TEXT NOT NULL,
                                     level INT NOT NULL DEFAULT 1,
                                     xp BIGINT NOT NULL DEFAULT 0,
                                     gold BIGINT NOT NULL DEFAULT 0,
                                     stats JSONB NOT NULL DEFAULT '{"focus":1,"discipline":1,"energy":1,"creativity":1}',
                                     streak INT NOT NULL DEFAULT 0,
                                     last_active_date DATE,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS quests (
                                      id UUID PRIMARY KEY,
                                      user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    weight INT NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'open',
    due_at TIMESTAMPTZ,
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS focus_runs (
                                          id UUID PRIMARY KEY,
                                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_id UUID REFERENCES quests(id) ON DELETE SET NULL,
    dungeon_rank TEXT NOT NULL,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ,
    target_minutes INT NOT NULL,
    result TEXT,
    xp_earned BIGINT NOT NULL DEFAULT 0,
    gold_earned BIGINT NOT NULL DEFAULT 0
    );

CREATE INDEX IF NOT EXISTS idx_focus_runs_user ON focus_runs(user_id);