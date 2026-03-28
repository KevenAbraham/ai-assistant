DELETE FROM messages WHERE role IN ('tool_use', 'tool_result');

ALTER TABLE messages
  DROP CONSTRAINT messages_role_check;

ALTER TABLE messages
  ADD CONSTRAINT messages_role_check
    CHECK (role IN ('user', 'assistant', 'system'));
