-- Expand the role CHECK constraint to allow tool_use and tool_result messages.
-- These are needed to persist multi-turn tool call context across requests.
ALTER TABLE messages
  DROP CONSTRAINT messages_role_check;

ALTER TABLE messages
  ADD CONSTRAINT messages_role_check
    CHECK (role IN ('user', 'assistant', 'system', 'tool_use', 'tool_result'));
