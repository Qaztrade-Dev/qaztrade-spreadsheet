begin;

alter table "applications"
  drop column "no"
;

drop table "application_signings";

commit;
