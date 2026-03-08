-- Users table
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    username        VARCHAR(50) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Snippets table
CREATE TABLE snippets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    slug            VARCHAR(100) NOT NULL,
    description     TEXT DEFAULT '',
    content         TEXT NOT NULL,
    variables       JSONB DEFAULT '[]',
    tags            TEXT[] DEFAULT '{}',
    visibility      VARCHAR(20) DEFAULT 'private' CHECK (visibility IN ('public', 'private', 'team')),
    owner_type      VARCHAR(20) NOT NULL CHECK (owner_type IN ('user', 'team')),
    owner_id        UUID NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    created_by      UUID REFERENCES users(id),
    version         INTEGER DEFAULT 1,
    pull_count      INTEGER DEFAULT 0,
    search_vector   tsvector,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(owner_type, owner_id, slug)
);

-- Full-text search function
CREATE FUNCTION snippets_search_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.name, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(array_to_string(NEW.tags, ' '), '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Full-text search trigger
CREATE TRIGGER snippets_search_trigger
    BEFORE INSERT OR UPDATE ON snippets
    FOR EACH ROW EXECUTE FUNCTION snippets_search_update();

-- Indexes for performance
CREATE INDEX idx_snippets_slug ON snippets(slug);
CREATE INDEX idx_snippets_namespace_slug ON snippets(namespace, slug);
CREATE INDEX idx_snippets_owner ON snippets(owner_type, owner_id);
CREATE INDEX idx_snippets_visibility ON snippets(visibility);
CREATE INDEX idx_snippets_tags ON snippets USING GIN(tags);
CREATE INDEX idx_snippets_search ON snippets USING GIN(search_vector);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Updated_at auto-update function
CREATE FUNCTION update_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Updated_at triggers
CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER snippets_updated_at
    BEFORE UPDATE ON snippets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
