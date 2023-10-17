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
  name VARCHAR(128) NOT NULL,
  
  UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS pull_request_event (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  pr_id INT NOT NULL,
  action VARCHAR(128) NOT NULL,
  sha VARCHAR(40) NOT NULL,
  is_merged BOOLEAN NOT NULL DEFAULT false,
  last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
  
  CONSTRAINT FK_pull_request_event_pull_request FOREIGN KEY (pr_id) REFERENCES pull_request(pr_id),
  CONSTRAINT FK_pull_request_event_pull_request_action FOREIGN KEY (action) REFERENCES pull_request_action(name)
);

CREATE TABLE IF NOT EXISTS approval (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  uuid VARCHAR(36) NOT NULL,
  pr_id INT NOT NULL,
  sha VARCHAR(40) NOT NULL,
  is_approved BOOLEAN NOT NULL DEFAULT false,
  last_updated TIMESTAMP NOT NULL DEFAULT NOW(),

  UNIQUE(uuid, pr_id, sha),
  CONSTRAINT FK_approval_pull_request_id FOREIGN KEY (pr_id) REFERENCES pull_request(pr_id)
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

-- TODO: Not sure if pr_id is unique across orgs/accounts. It may be best to use pull_request(id) as the foreign key in other tables