begin;

drop index applications_is_signed_idx;
drop index tnved_code_idx;
drop index tnved_level_idx;

alter table "assignments"
    drop column "total_rows",
    drop column "total_sum"
;

alter table "assignments"
  add column "rows_from" bigint default 0,
  add column "rows_until" bigint default 0
;

commit;
