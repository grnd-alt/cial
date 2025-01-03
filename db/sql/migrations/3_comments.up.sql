create table if not exists comments (
    id varchar(255) primary key,
    post_id varchar(255) not null,
    user_id varchar(255) not null,
    content text not null,
    created_at timestamp not null default now (),
    updated_at timestamp not null default now (),
    foreign key (post_id) references posts (id) on delete cascade
);
