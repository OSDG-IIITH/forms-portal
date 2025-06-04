-- USERS
CREATE TABLE users (
    uid TEXT PRIMARY KEY NOT NULL, -- ULID vs first part of the mail id?
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE CHECK (email ~* '^[^@]+@([a-zA-Z0-9.-]+\\.)?iiit\\.ac\\.in$'),
    role user_role NOT NULL
);

-- FORM SETTINGS
CREATE TABLE form_settings (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_uid UUID REFERENCES forms(uid) ON DELETE CASCADE,
    anonymity BOOLEAN DEFAULT FALSE,
    max_responses INTEGER,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    allowed_roles user_role[],
    allow_editing BOOLEAN DEFAULT FALSE,
    time_limit_seconds INTEGER
);

-- FORMS
CREATE TABLE forms (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_uid TEXT REFERENCES users(uid) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    visibility form_visibility NOT NULL DEFAULT 'private', -- accessible only to the list shared
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    settings_id UUID REFERENCES form_settings(id)
);

-- FORM PERMISSIONS - TO DISCUSS AND EDIT
CREATE TABLE form_permissions (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_uid UUID REFERENCES forms(uid) ON DELETE CASCADE,
    user_uid TEXT REFERENCES users(uid) ON DELETE CASCADE,
    role TEXT CHECK (role IN ('respondent', 'editor', 'owner', 'commentator')),
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- FORM ACCESS LOGS
CREATE TABLE form_access_logs (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_uid UUID REFERENCES forms(uid) ON DELETE CASCADE,
    user_uid TEXT REFERENCES users(uid),
    action access_action NOT NULL,
    performed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- QUESTIONS
CREATE TABLE questions (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_uid UUID REFERENCES forms(uid) ON DELETE CASCADE,
    type question_type NOT NULL,
    text TEXT NOT NULL,
    required BOOLEAN DEFAULT FALSE,
    display_order INTEGER,
    conditional_logic JSONB
);

-- OPTIONS
CREATE TABLE options (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_uid UUID REFERENCES questions(uid) ON DELETE CASCADE,
    type option_type NOT NULL,
    text TEXT,
    validation JSONB
);

-- RESPONSES
CREATE TABLE responses (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_uid UUID REFERENCES forms(uid) ON DELETE CASCADE,
    user_uid TEXT REFERENCES users(uid),
    submitted_at TIMESTAMP
);

-- FORM ANSWERS
CREATE TABLE form_answers (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    response_uid UUID REFERENCES responses(uid) ON DELETE CASCADE,
    form_uid UUID REFERENCES forms(uid),
    user_uid TEXT REFERENCES users(uid),
    question_uid UUID REFERENCES questions(uid),
    answer_value TEXT,
    answer_array TEXT[],
    submitted_at TIMESTAMP
);

-- FILES
CREATE TABLE files (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    response_uid UUID REFERENCES responses(uid) ON DELETE CASCADE,
    question_uid UUID REFERENCES questions(uid),
    path TEXT NOT NULL,
    original_name TEXT,
    mime_type TEXT,
    file_size BIGINT
);
