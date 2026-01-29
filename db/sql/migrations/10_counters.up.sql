CREATE TABLE counters(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  icon VARCHAR(255)
);

CREATE TABLE counters_users(
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(255) not null,
  counter_id int not null,
  token VARCHAR(255),
  access_type VARCHAR(255),
  UNIQUE(user_id, counter_id),
  FOREIGN KEY (counter_id) REFERENCES counters(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE counters_users_events(
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(255) not null,
  counter_id int not null,
  created_at timestamp with time zone DEFAULT now(),
  FOREIGN KEY (counter_id) REFERENCES counters(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id, counter_id) REFERENCES counters_users(user_id, counter_id) ON DELETE CASCADE
);
