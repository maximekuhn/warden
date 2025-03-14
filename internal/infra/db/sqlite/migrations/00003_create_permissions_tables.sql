CREATE TABLE user_policy_plan (
    user_id  TEXT PRIMARY KEY,
    user_plan INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES auth(user_id) ON DELETE CASCADE
);

CREATE TABLE user_role_server (
    user_id TEXT NOT NULL,
    server_id TEXT NOT NULL,
    user_role INTEGER NOT NULL,
    PRIMARY KEY (user_id, server_id),
    FOREIGN KEY (user_id) REFERENCES user_policy_plan(user_id) ON DELETE CASCADE
);
