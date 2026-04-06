CREATE TABLE tb_auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE tb_refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    auth_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false,
    CONSTRAINT fk_user FOREIGN KEY (auth_id) REFERENCES tb_auth(id) ON DELETE cascade
);

CREATE TABLE tb_verify_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    auth_id UUID NOT NULL UNIQUE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_verify_auth FOREIGN KEY (auth_id) REFERENCES tb_auth(id) ON DELETE CASCADE
);

CREATE TABLE tb_short_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    own_by UUID NOT NULL,
    view INTEGER NOT NULL DEFAULT 0,
    target_url VARCHAR(255) NOT NULL,
    short_url VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    CONSTRAINT fk_short_link_auth FOREIGN KEY (own_by) REFERENCES tb_auth(id) ON DELETE CASCADE
);
