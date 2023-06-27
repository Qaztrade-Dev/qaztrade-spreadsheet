begin;

create table "batches" (
  "id" serial primary key,
  "created_at" timestamptz default now(),
  "step" int default 0,
  "is_completed" boolean default false,
  "completed_at" timestamptz default null
);
create index batches_id_idx on "batches" using hash ("id");

create table "batch_applications" (
  "created_at" timestamptz default now(),
  "batch_id" int,
  "application_id" uuid,
  "is_completed" boolean default false,
  "completed_at" timestamptz default null,
  foreign key ("batch_id") references "batches" ("id") on delete cascade,
  foreign key ("application_id") references "applications" ("id") on delete cascade,
  primary key ("batch_id", "application_id")
);
create index batch_applications_batch_id_idx on "batch_applications" using hash ("batch_id");
create index batch_applications_application_id_idx on "batch_applications" using hash ("application_id");

create table "assignments" (
  "id" serial primary key,
  "created_at" timestamptz default now(),
  "user_id" uuid,
  "application_id" uuid,
  "type" text, -- 'digital', 'finance', 'legal'
  "sheet_title" text,
  "sheet_id" bigint,
  "rows_from" bigint,
  "rows_until" bigint,
  "rows_total" bigint GENERATED ALWAYS AS ("rows_until" - "rows_from" +1) STORED,
  "is_completed" boolean default false,
  "completed_at" timestamptz default null,
  "last_result_id" uuid default null,
  foreign key ("user_id") references "users" ("id") on delete set null,
  foreign key ("application_id") references "applications" ("id") on delete cascade
);

create index assignments_id_idx on "assignments" using hash ("id");
create index assignments_type_idx on "assignments" using hash ("type");
create index assignments_user_id_idx on "assignments" using hash ("user_id");
create index assignments_application_id_idx on "assignments" using hash ("application_id");

create table "assignment_results" (
  "id" uuid primary key default gen_random_uuid(),
  "created_at" timestamptz default now(),
  "assignment_id" int,
  "total_completed" bigint default 0,
  foreign key ("assignment_id") references "assignments" ("id") on delete cascade
);
create index assignment_results_id_idx on "assignment_results" using hash ("id");
create index assignment_results_assignment_id_idx on "assignment_results" using hash ("assignment_id");

create table "user_role_bindings" (
  "created_at" timestamptz default now(),
  "user_id" uuid,
  "role_id" int,
  foreign key ("user_id") references "users" ("id") on delete cascade,
  foreign key ("role_id") references "user_roles" ("id") on delete cascade,
  primary key ("user_id", "role_id")
);
create index user_role_bindings_user_id_idx on "user_role_bindings" using hash ("user_id");

insert into "user_roles" 
  ("id", "value")
values
  (3, 'admin'),
  (4, 'digital'),
  (5, 'finance'),
  (6, 'legal')
;

insert into "user_role_bindings" ("user_id", "role_id")
select u.id, u.role_id
from "users" u;

commit;
