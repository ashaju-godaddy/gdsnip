-- Drop triggers
DROP TRIGGER IF EXISTS snippets_updated_at ON snippets;
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP TRIGGER IF EXISTS snippets_search_trigger ON snippets;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at();
DROP FUNCTION IF EXISTS snippets_search_update();

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS snippets;
DROP TABLE IF EXISTS users;
