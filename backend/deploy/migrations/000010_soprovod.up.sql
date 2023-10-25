begin;

alter table "assignments"
    add column "resolved_at" timestamptz,
    add column "countdown_duration" interval,
    add column "resolution_status_id" int
;

create table "assignment_messages" (
    "id" uuid primary key default gen_random_uuid(),
    "created_at" timestamptz default now(),
    "user_id" uuid,
    "assignment_id" int,
    "attrs" jsonb default '{}'::jsonb,
    "doodocs_document_id" text,
    "doodocs_signed_at" timestamptz,
    "doodocs_is_signed" boolean,
    foreign key ("user_id") references "users" ("id"),
    foreign key ("assignment_id") references "assignments" ("id")
);

commit;
