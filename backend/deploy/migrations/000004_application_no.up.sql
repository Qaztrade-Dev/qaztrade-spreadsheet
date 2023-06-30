begin;

create table "application_signings" (
  "id" serial primary key,
  "created_at" timestamptz default now(),
  "application_id" uuid,
  foreign key ("application_id") references "applications" ("id") on delete cascade
);
create index application_signings_application_id_idx on "application_signings" using hash ("application_id");

alter table "applications"
  add column "no" int default -1
;

insert into "application_signings" (application_id)
select
    id
from applications
where
    is_signed = true
order by 
    sign_at asc
;

update "applications" as app
set
    "no" = appsig.id
from application_signings appsig
where 
    appsig.application_id = app.id
;

commit;
