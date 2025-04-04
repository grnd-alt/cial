CREATE TABLE IF NOT EXISTS user_follows(
    follower_id varchar(255) not null,
    followed_id varchar(255) not null,
    notification_type varchar(255),
    followed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    primary key (follower_id, followed_id)
);