INSERT INTO
  pull_request_action(name)
VALUES
  ('opened'),
  ('closed'),
  ('reopened'),
  ('synchronize')
;

INSERT INTO
  approval_schema(name, description, status_context, status_title)
VALUES
  ('v1', '', 'Pull Request Compliance', 'User Review Required')
;

INSERT INTO
  default_approval_schema(id, schema_id)
VALUES
  (1, (SELECT id FROM approval_schema WHERE name = 'v1'))
;

INSERT INTO
  approval_yesno_question (question_text, schema_id)
VALUES
(
  'Have you read the contribution guidelines',
  (SELECT id FROM approval_schema WHERE name = 'v1')
),
(
  'Do you like coffee',
  (SELECT id FROM approval_schema WHERE name = 'v1')
),
(
  'Have you signed the license agreement',
  (SELECT id FROM approval_schema WHERE name = 'v1')
)
;