-- Drop trigger
DROP TRIGGER IF EXISTS teams_updated_at ON teams;

-- Drop indexes
DROP INDEX IF EXISTS idx_snippets_team_visibility;
DROP INDEX IF EXISTS idx_team_members_role;
DROP INDEX IF EXISTS idx_team_members_user;
DROP INDEX IF EXISTS idx_team_members_team;
DROP INDEX IF EXISTS idx_teams_created_by;
DROP INDEX IF EXISTS idx_teams_slug;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
