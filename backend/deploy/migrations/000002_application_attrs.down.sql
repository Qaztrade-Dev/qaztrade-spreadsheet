begin;

alter table "applications"
  drop column "attrs"
;

commit;
