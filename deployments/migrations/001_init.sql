-- +goose Up
-- +goose StatementBegin

CREATE TYPE role_in_department AS ENUM ('trainee', 'activist');

CREATE TABLE organization
(
    id   UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE department
(
    id   UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE "user"
(
    id              UUID PRIMARY KEY,
    full_name       VARCHAR(255)        NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    hash_password   TEXT                NOT NULL,
    nice            INT,
    is_admin        BOOLEAN,
    role            role_in_department,
    organization_id UUID,
    department_id   UUID,
    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE SET NULL,
    FOREIGN KEY (department_id) REFERENCES department (id) ON DELETE SET NULL
);

CREATE TYPE equipment_status AS ENUM (
    'available',
    'issued',
    'under_maintenance',
    'unavailable'
    );

CREATE TABLE equipment
(
    id                   UUID PRIMARY KEY,
    inventory_number     VARCHAR(100) UNIQUE NOT NULL,
    name                 VARCHAR(255)        NOT NULL,
    short_name           VARCHAR(100),
    category             VARCHAR(100)        NOT NULL,
    available_to_trainee BOOLEAN DEFAULT FALSE,
    status               equipment_status,
    department_id        UUID,
    FOREIGN KEY (department_id) REFERENCES department (id) ON DELETE SET NULL
);

CREATE TYPE equipment_invocation_status AS ENUM (
    'created',
    'under_review',
    'changes_required',
    'approved',
    'equipment_issued',
    'equipment_returned',
    'completed',
    'cancelled'
    );

CREATE TABLE equipment_invocation
(
    id                    UUID PRIMARY KEY,
    event_name            VARCHAR(255) NOT NULL,
    start_time            TIMESTAMP    NOT NULL,
    end_time              TIMESTAMP    NOT NULL,
    equipment_return_time TIMESTAMP    NOT NULL,
    sd_card_return_time   TIMESTAMP    NOT NULL,
    status                equipment_invocation_status,
    curator_comment       TEXT,

    equipment_id          UUID         NOT NULL,
    organization_id       UUID         NOT NULL,
    department_id         UUID         NOT NULL,
    user_id               UUID         NOT NULL,
    admin_id              UUID         NOT NULL,
    FOREIGN KEY (equipment_id) REFERENCES equipment (id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE SET NULL,
    FOREIGN KEY (department_id) REFERENCES department (id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE SET NULL,
    FOREIGN KEY (admin_id) REFERENCES "user" (id) ON DELETE SET NULL
);

CREATE TABLE message_equipment_invocation
(
    id            UUID PRIMARY KEY,
    content       TEXT NOT NULL,
    sent_time     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    invocation_id UUID NOT NULL,
    sender_id     UUID NOT NULL,
    recipient_id  UUID NOT NULL,
    FOREIGN KEY (invocation_id) REFERENCES equipment_invocation (id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TYPE studio_invocation_status AS ENUM (
    'created',
    'under_review',
    'changes_required',
    'approved',
    'completed',
    'cancelled'
    );

CREATE TABLE studio_invocation
(
    id                   UUID PRIMARY KEY,
    event_name           VARCHAR(255)             NOT NULL,
    shooting_description TEXT                     NOT NULL,
    start_time           TIMESTAMP                NOT NULL,
    end_time             TIMESTAMP                NOT NULL,
    needs_chromakey      BOOLEAN DEFAULT FALSE,
    needs_cyclorama      BOOLEAN DEFAULT FALSE,
    needs_black_fabric   BOOLEAN DEFAULT FALSE,
    status               studio_invocation_status NOT NULL,
    curator_comment      TEXT,

    organization_id      UUID                     NOT NULL,
    department_id        UUID                     NOT NULL,
    user_id              UUID                     NOT NULL,
    admin_id             UUID                     NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE SET NULL,
    FOREIGN KEY (department_id) REFERENCES department (id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE SET NULL,
    FOREIGN KEY (admin_id) REFERENCES "user" (id) ON DELETE SET NULL
);


CREATE TABLE message_studio_invocation
(
    id            UUID PRIMARY KEY,
    content       TEXT NOT NULL,
    sent_time     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    invocation_id UUID NOT NULL,
    sender_id     UUID NOT NULL,
    recipient_id  UUID NOT NULL,
    FOREIGN KEY (invocation_id) REFERENCES studio_invocation (id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_email ON "user" (email);
CREATE INDEX idx_equipment_inventory ON equipment (inventory_number);
CREATE INDEX idx_message_studio_invocation_sender ON message_studio_invocation (sender_id);
CREATE INDEX idx_message_studio_invocation_recipient ON message_studio_invocation (recipient_id);
CREATE INDEX idx_message_equipment_invocation_sender ON message_equipment_invocation (sender_id);
CREATE INDEX idx_message_equipment_invocation_recipient ON message_equipment_invocation (recipient_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_equipment_inventory;
DROP INDEX IF EXISTS idx_message_studio_invocation_sender;
DROP INDEX IF EXISTS idx_message_studio_invocation_recipient;
DROP INDEX IF EXISTS idx_message_equipment_invocation_sender;
DROP INDEX IF EXISTS idx_message_equipment_invocation_recipient;


DROP TABLE IF EXISTS message_studio_invocation;
DROP TABLE IF EXISTS message_equipment_invocation;
DROP TABLE IF EXISTS studio_invocation;
DROP TABLE IF EXISTS equipment_invocation;
DROP TABLE IF EXISTS equipment;
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS department;
DROP TABLE IF EXISTS organization;

DROP TYPE IF EXISTS studio_invocation_status;
DROP TYPE IF EXISTS equipment_invocation_status;
DROP TYPE IF EXISTS equipment_status;
DROP TYPE IF EXISTS role_in_department;
-- +goose StatementEnd