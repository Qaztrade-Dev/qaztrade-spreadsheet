begin;

create table "assignments" (
  "id" uuid primary key default gen_random_uuid(),
  "created_at" timestamptz default now(),
  "user_id" uuid,
  "application_id" uuid,
  "sheet_title" text,
  "sheet_id" bigint,
  "rows_from" bigint,
  "rows_to" bigint,
  foreign key ("user_id") references "users" ("id"),
  foreign key ("application_id") references "applications" ("id")
);

create index assignments_id_idx on "assignments" using hash ("id");
create index assignments_user_id_idx on "assignments" using hash ("user_id");
create index assignments_application_id_idx on "assignments" using hash ("application_id");

create table "assignment_results" (
  "id" uuid primary key default gen_random_uuid(),
  "created_at" timestamptz default now(),
  "assignment_id" uuid,
  "total_completed" bigint,
  foreign key ("assignment_id") references "assignments" ("id")
);

create index assignment_results_assignment_id_idx on "assignment_results" using hash ("assignment_id");

commit;
