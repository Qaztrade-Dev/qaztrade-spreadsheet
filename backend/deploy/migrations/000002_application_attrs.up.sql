begin;

alter table "applications"
  add column "attrs" jsonb default '{}'::jsonb
;

commit;
