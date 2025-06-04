-- USER ROLES
CREATE TYPE user_role AS ENUM ('student', 'faculty', 'club', 'portal_admin');

-- FORM VISIBILITY
CREATE TYPE form_visibility AS ENUM ('public', 'private');

-- QUESTION TYPES
CREATE TYPE question_type AS ENUM (
  'text',
  'multiple_choice',
  'checkbox',
  'dropdown',
  'file_upload',
  'date_picker',
  'rating_scale'
);

-- OPTION TYPES
CREATE TYPE option_type AS ENUM ('text', 'number', 'array', 'upload');

-- ACCESS ACTIONS
CREATE TYPE access_action AS ENUM ('edited', 'deleted', 'submitted');
