begin;

create table "users" (
  "id" uuid primary key default gen_random_uuid(),
  "created_at" timestamptz default now(),
  "email" text not null,
  "hashed_password" text not null,
  "attrs" jsonb default {}
);

comment on column "users.attrs" is 'Attributes of the user, i.e. "{"org_name": "OpenAI Inc."}"';

create table "application_statuses" (
  "id" int primary key,
  "value" text
);

insert into "application_statuses" 
  ("id", "value")
values
  (1, "user_filling"), -- Пользователь заполняет
  (2, "manager_reviewing"), -- Менеджер проверяет
  (3, "completed"), -- Завершен
  (4, "rejected") -- Отклонен
;

create table "applications" (
  "id" uuid primary key default gen_random_uuid(),
  "created_at" timestamptz default now(),
  "user_id" uuid,
  "status_id" int,
  "spreadsheet_id" text,
  foreign key ("user_id") references "users" ("id"),
  foreign key ("status_id") references "application_statuses" ("id"),
);

end;
