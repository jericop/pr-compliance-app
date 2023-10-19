-- TODO: Not sure if pr_id is unique across orgs/accounts. It may be best to use pull_request(id) as the foreign key in other tables
-- NOT NULL is used frequently to ensure that sqlc generates native go types rather than nullable pg types

CREATE TABLE IF NOT EXISTS installation (
  id INT NOT NULL,
  
  CONSTRAINT PK_installation_id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS repo (
  id INT NOT NULL,
  org TEXT NOT NULL,
  name TEXT NOT NULL,
  
  CONSTRAINT PK_repo_id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS gh_user (
  id INT NOT NULL,
  login TEXT NOT NULL,

  CONSTRAINT PK_gh_user_id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS pull_request (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  repo_id INT NOT NULL,
  pr_id INT NOT NULL,
  pr_number INT NOT NULL,
  opened_by INT NOT NULL ,
  installation_id INT NOT NULL,
  is_merged BOOLEAN NOT NULL DEFAULT false,
  
  UNIQUE(pr_id),
  UNIQUE(repo_id, pr_number),
  CONSTRAINT FK_pull_request_repo FOREIGN KEY (repo_id) REFERENCES repo(id),
  CONSTRAINT FK_pull_request_gh_user FOREIGN KEY (opened_by) REFERENCES gh_user(id),
  CONSTRAINT FK_pull_request_installation FOREIGN KEY (installation_id) REFERENCES installation(id)
);

CREATE TABLE IF NOT EXISTS pull_request_action (
  name TEXT NOT NULL,
  
  UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS pull_request_event (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  pr_id INT NOT NULL,
  action TEXT NOT NULL,
  sha TEXT NOT NULL,
  is_merged BOOLEAN NOT NULL DEFAULT false,
  last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
  
  CONSTRAINT FK_pull_request_event_pull_request FOREIGN KEY (pr_id) REFERENCES pull_request(pr_id),
  CONSTRAINT FK_pull_request_event_pull_request_action FOREIGN KEY (action) REFERENCES pull_request_action(name)
);

CREATE TABLE IF NOT EXISTS approval_schema (
  id INT GENERATED ALWAYS AS IDENTITY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  status_context TEXT NOT NULL,
  status_title TEXT NOT NULL,
  
  CONSTRAINT PK_approval_schema_id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS default_approval_schema (
  id INT NOT NULL,
  schema_id INT NOT NULL,

  CONSTRAINT PK_default_approval_schema UNIQUE (id),
  CONSTRAINT C_default_approval_schema_limit_1 CHECK (id = 1),
  CONSTRAINT FK_default_approval_schema_schema_id FOREIGN KEY (schema_id) REFERENCES approval_schema(id)
);

CREATE TABLE IF NOT EXISTS approval (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  schema_id INT NOT NULL,
  uuid TEXT NOT NULL,
  pr_id INT NOT NULL,
  sha TEXT NOT NULL,
  is_approved BOOLEAN NOT NULL DEFAULT false,
  last_updated TIMESTAMP NOT NULL DEFAULT NOW(),

  UNIQUE(schema_id, uuid, pr_id, sha),
  CONSTRAINT FK_approval_schema_id FOREIGN KEY (schema_id) REFERENCES approval_schema(id),
  CONSTRAINT FK_approval_pr_id FOREIGN KEY (pr_id) REFERENCES pull_request(pr_id)
);

CREATE OR REPLACE FUNCTION trigger_set_last_updated()
RETURNS TRIGGER AS $$
BEGIN
  NEW.last_updated = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER approval_last_updated
  BEFORE UPDATE ON approval
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_set_last_updated();

CREATE TABLE IF NOT EXISTS approval_yesno_question (
  id INT GENERATED ALWAYS AS IDENTITY,
  schema_id INT NOT NULL,
  question_text TEXT NOT NULL,
  
  -- CONSTRAINT PK_approval_yesno_question UNIQUE (schema_id, id),
  CONSTRAINT PK_approval_yesno_question_id PRIMARY KEY (id),
  CONSTRAINT FK_approval_yesno_question_schema_id FOREIGN KEY (schema_id) REFERENCES approval_schema(id)
);

-- Only questions with yes answers are recoreded, which means all other answers are false by default
CREATE TABLE IF NOT EXISTS approval_yes_answer (
  approval_id INT NOT NULL,
  question_id INT NOT NULL,
  
  CONSTRAINT PK_approval_yes_answer UNIQUE (approval_id, question_id),
  CONSTRAINT FK_approval_yes_answer_approval_id FOREIGN KEY (approval_id) REFERENCES approval(id),
  CONSTRAINT FK_approval_yes_answer_question_id FOREIGN KEY (question_id) REFERENCES approval_yesno_question(id)
);
