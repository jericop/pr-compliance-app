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
  installation(id)
VALUES
  (8675309)
;

INSERT INTO pull_request (pr_id, pr_number, repo_id, is_merged, installation_id, opened_by)
VALUES
    (991, 1, 1, 'false', 8675309, 1),
    (992, 2, 1, 'false', 8675309, 2),
    (993, 1, 2, 'true', 8675309, 3), -- already merged
    (994, 2, 2, 'false', 8675309, 3)
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
    (993, 'synchronize', '4bcfe98e640c8284511312660fb8709b0afa888e', 'false'),
    
    -- This is a new PR that gets a new commit pushed and then it is closed without being merged
    (994, 'opened', '88efe98e640c8284511312660fb8709b0afa888e', 'false'),
    (994, 'synchronize', 'abcde798228d17af2d34fce4cfbdf35556832472', 'false')
;

-- This is inserting data into the compliance approval table, not be confused with an approval for the pull request.
INSERT INTO approval(schema_id, uuid, pr_id, sha, is_approved)
VALUES
  -- This would be an example where approvals were done for each commit
  (1, '82a844d8-fccc-47e1-a3fd-008b17b67510', 991, '78981922613b2afb6025042ff6bd878ac1994e85', true),
  (1, '82a844d8-fccc-47e1-a3fd-008b17b67510', 991, '61780798228d17af2d34fce4cfbdf35556832472', true),
  
  -- This is an example where the approval was done for the last commit before being merged
  (1, 'a0b3adb0-174c-4be5-984e-3005aeffbf65', 992, 'f2ad6c76f0115a6ba5b00456a849810e7ec0af20', true),
  
  -- This was already merged
  (1, '7829491b-9bdc-4167-87e3-73a334fb5916', 993, '4bcfe98e640c8284511312660fb8709b0afa888e', true),
  
  (1, 'b329487e-8qv4-8b17-87e3-79bd34fb59bd', 994, '4bcfe98e640c8284511312660fb8709b0afa888e', false)
;

-- select id, question_text from approval_question where schema_id = (select id from approval_schema where name = 'v1') group by id,question_text having id = min(id);

-- Get questions (GROUP BY also orders by id column in ascending order)
SELECT id, question_text FROM approval_yesno_question
WHERE schema_id = (select id from approval_schema where name = 'v1') 
GROUP BY id, question_text;

-- Add yes answers for the first and last questions in the default schema
WITH sorted_questions AS (
  SELECT schema_id, id, question_text FROM approval_yesno_question
  WHERE schema_id = (SELECT schema_id FROM default_approval_schema LIMIT 1)
  GROUP BY schema_id, id, question_text
)
INSERT INTO approval_yes_answer (approval_id, question_id)
VALUES
(
  (SELECT id FROM approval WHERE uuid = 'b329487e-8qv4-8b17-87e3-79bd34fb59bd'),
  (SELECT id FROM sorted_questions LIMIT 1)
),
(
  (SELECT id FROM approval WHERE uuid = 'b329487e-8qv4-8b17-87e3-79bd34fb59bd'),
  (SELECT MAX(id) FROM sorted_questions LIMIT 1)
)
;
