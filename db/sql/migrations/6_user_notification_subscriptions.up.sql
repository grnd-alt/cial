CREATE TABLE IF NOT EXISTS user_subscriptions(
    user_id varchar(255) not null,
    subscription JSONB UNIQUE,
    created_at timestamp with time zone DEFAULT now(),
    primary key (user_id, subscription)
);