create table calls
(
    id       bigint unsigned auto_increment
        primary key,
    chat     text      null,
    start_at timestamp null,
    end_at   timestamp null
);

create table users
(
    id         bigint unsigned auto_increment
        primary key,
    created_at timestamp default CURRENT_TIMESTAMP null,
    updated_at timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    deleted_at timestamp                           null,
    kakao_id   bigint unsigned                     null,
    name       varchar(100)                        null,
    phone_no   varchar(100)                        null
);

create table parents
(
    id         bigint unsigned auto_increment
        primary key,
    created_at timestamp default CURRENT_TIMESTAMP null,
    updated_at timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    deleted_at timestamp                           null,
    user_id    bigint unsigned                     null,
    name       varchar(100)                        null,
    phone_no   varchar(100)                        null,
    constraint parents_ibfk_1
        foreign key (user_id) references users (id)
);

create index user_id
    on parents (user_id);

create table reports
(
    id                bigint unsigned auto_increment
        primary key,
    recipient_user_id bigint unsigned null,
    title             varchar(255)    null,
    body              text            null,
    parent_id         bigint unsigned null,
    status            varchar(255)    null,
    created_at        timestamp       null,
    call_id           bigint unsigned null,
    constraint reports_ibfk_1
        foreign key (recipient_user_id) references users (id),
    constraint reports_ibfk_2
        foreign key (call_id) references calls (id),
    constraint reports_ibfk_3
        foreign key (parent_id) references parents (id)
);

create index call_id
    on reports (call_id);

create index recipient_user_id
    on reports (recipient_user_id);

create table subscriptions
(
    id         bigint unsigned auto_increment
        primary key,
    created_at timestamp default CURRENT_TIMESTAMP null,
    updated_at timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    deleted_at timestamp                           null,
    user_id    bigint unsigned                     null,
    expires_at date                                null,
    constraint subscriptions_ibfk_1
        foreign key (user_id) references users (id)
);

create index user_id
    on subscriptions (user_id);

