-- TODO: remove this after testing as migrations should run as a specific user and connect to a specific database
-- \connect pr_compliance

CREATE TABLE IF NOT EXISTS repo (
  id INT NOT NULL,
  org TEXT NOT NULL,
  name TEXT NOT NULL,
  
  -- UNIQUE(org, name)
  CONSTRAINT PK_user_id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS gh_user (
  id INT NOT NULL,
  login TEXT NOT NULL,

  CONSTRAINT PK_gh_user_id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS pull_request (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  repo_id INT,
  pr_id INT NOT NULL,
  pr_number INT NOT NULL,
  opened_by INT,
  is_merged BOOLEAN DEFAULT false,
  
  UNIQUE(pr_id),
  UNIQUE(repo_id, pr_number),
  CONSTRAINT FK_pull_request_repo FOREIGN KEY (repo_id) REFERENCES repo(id),
  CONSTRAINT FK_pull_request_gh_user FOREIGN KEY (opened_by) REFERENCES gh_user(id)
);

CREATE TABLE IF NOT EXISTS pull_request_action (
  name VARCHAR(128) NOT NULL,
  
  UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS pull_request_event (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  pr_id INT,
  action VARCHAR(128),
  sha VARCHAR(40),
  is_merged BOOLEAN DEFAULT false,
  last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
  
  CONSTRAINT FK_pull_request_event_pull_request FOREIGN KEY (pr_id) REFERENCES pull_request(pr_id),
  CONSTRAINT FK_pull_request_event_pull_request_action FOREIGN KEY (action) REFERENCES pull_request_action(name)
);

CREATE TABLE IF NOT EXISTS approval (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  uuid VARCHAR(36) NOT NULL,
  pr_id INT,
  sha VARCHAR(40) NOT NULL,
  approved_on TIMESTAMP NOT NULL DEFAULT NOW(),

  UNIQUE(uuid, pr_id, sha),
  CONSTRAINT FK_approval_pull_request_id FOREIGN KEY (pr_id) REFERENCES pull_request(pr_id)
);

-- TODO: create trigger to update pull_request(is_merged) when pull_request_event action is received that has that field set to true.