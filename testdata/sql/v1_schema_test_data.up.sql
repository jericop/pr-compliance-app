-- TODO: remove this after testing as migrations should run as a specific user and connect to a specific database
-- \connect pr_compliance;

INSERT INTO
  repo(org, name, id)
VALUES
  ('org1', 'web-app-api', 1),
  ('org1', 'web-app-ui', 2)
;

INSERT INTO
  gh_user(id, login)
VALUES
  (1, 'user1'),
  (2, 'user2'),
  (3, 'user3')
;

INSERT INTO
  pull_request_action(name)
VALUES
  ('opened'),
  ('closed'),
  ('reopened'),
  ('synchronize')
;

INSERT INTO pull_request (pr_id, pr_number, repo_id, is_merged, opened_by)
VALUES
    (991, 1, 1, 'false', 1),
    (992, 2, 1, 'false', 2),
    (993, 1, 2, 'true', 3) -- already merged
;

INSERT INTO pull_request_event (pr_id, action, sha, is_merged)
VALUES
    -- This PR is opened, closed, and reopened again
    (991, 'opened', '78981922613b2afb6025042ff6bd878ac1994e85', 'false'),
    (991, 'closed', '78981922613b2afb6025042ff6bd878ac1994e85', 'false'),
    (991, 'reopened', '78981922613b2afb6025042ff6bd878ac1994e85', 'false'),
    -- Then a new commit is pushed before it is merged and closed
    (991, 'synchronize', '61780798228d17af2d34fce4cfbdf35556832472', 'false'),
    (991, 'closed', '61780798228d17af2d34fce4cfbdf35556832472', 'true'),
    
    -- This is an existing PR that has a new commit pushed and then it is closed without being merged
    (992, 'synchronize', 'f2ad6c76f0115a6ba5b00456a849810e7ec0af20', 'false'),
    (992, 'closed', 'f2ad6c76f0115a6ba5b00456a849810e7ec0af20', 'true'),

    -- This is a new PR that gets a new commit pushed and then it is closed without being merged
    (993, 'opened', '4bcfe98e640c8284511312660fb8709b0afa888e', 'false'),
    (993, 'synchronize', '4bcfe98e640c8284511312660fb8709b0afa888e', 'false')
;

-- This is inserting data into the compliance approval table, not be confused with an approval for the pull request.
INSERT INTO approval(uuid, pr_id, sha, approved_on)
VALUES
  -- This would be an example where approvals were done for each commit
  ('82a844d8-fccc-47e1-a3fd-008b17b67510', 991, '78981922613b2afb6025042ff6bd878ac1994e85', NOW()),
  ('82a844d8-fccc-47e1-a3fd-008b17b67510', 991, '61780798228d17af2d34fce4cfbdf35556832472', NOW()),
  
  -- This is an example where the approval was done for the last commit before being merged
  ('a0b3adb0-174c-4be5-984e-3005aeffbf65', 992, 'f2ad6c76f0115a6ba5b00456a849810e7ec0af20', NOW()),
  
  -- This is an example where the PR
  ('7829491b-9bdc-4167-87e3-73a334fb5916', 993, '4bcfe98e640c8284511312660fb8709b0afa888e', NOW())
;
