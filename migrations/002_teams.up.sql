-- Teams table
CREATE TABLE teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    description     TEXT DEFAULT '',
    created_by      UUID REFERENCES users(id) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Team memberships with roles
CREATE TABLE team_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID REFERENCES teams(id) ON DELETE CASCADE NOT NULL,
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'member'
                    CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);

-- Indexes for teams
CREATE INDEX idx_teams_slug ON teams(slug);
CREATE INDEX idx_teams_created_by ON teams(created_by);

-- Indexes for team_members
CREATE INDEX idx_team_members_team ON team_members(team_id);
CREATE INDEX idx_team_members_user ON team_members(user_id);
CREATE INDEX idx_team_members_role ON team_members(team_id, role);

-- Composite index for snippets team filtering (uses existing columns)
CREATE INDEX idx_snippets_team_visibility ON snippets(owner_type, owner_id, visibility);

-- Updated_at trigger for teams (reuses existing function from 001_initial)
CREATE TRIGGER teams_updated_at
    BEFORE UPDATE ON teams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
