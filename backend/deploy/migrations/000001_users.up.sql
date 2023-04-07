begin;

create table "user_roles" (
  "id" int primary key,
  "value" text
);

insert into "user_roles" 
  ("id", "value")
values
  (1, 'user'),
  (2, 'manager')
;

create table "users" (
  "id" uuid primary key default gen_random_uuid(),
  "created_at" timestamptz default now(),
  "email" text not null unique,
  "hashed_password" text not null,
  "attrs" jsonb default '{}'::jsonb,
  "role_id" int,
  foreign key ("role_id") references "user_roles" ("id")
);

create index users_id_idx on "users" using hash ("id");
create index users_email_idx on "users" using hash ("email");

comment on column users.attrs is 'Attributes of the user, i.e. "{"org_name": "OpenAI Inc."}"';

create table "application_statuses" (
  "id" int primary key,
  "value" text
);

insert into "application_statuses" 
  ("id", "value")
values
  (1, 'user_filling'), -- Пользователь заполняет
  (2, 'manager_reviewing'), -- Менеджер проверяет
  (3, 'user_fixing'), -- Пользователь исправляет
  (4, 'completed'), -- Завершен
  (5, 'rejected') -- Отклонен
;

create table "applications" (
  "id" uuid primary key default gen_random_uuid(),
  "created_at" timestamptz default now(),
  "user_id" uuid,
  "status_id" int,
  "spreadsheet_id" text,
  "link" text,
  "sign_link" text,
  "is_signed" boolean default false,
  foreign key ("user_id") references "users" ("id"),
  foreign key ("status_id") references "application_statuses" ("id")
);

create index applications_id_idx on "applications" using hash ("id");
create index applications_user_id_idx on "applications" using hash ("user_id");

create table "oauth2_tokens" (
  "id" serial primary key,
  "token" text
);



commit;
