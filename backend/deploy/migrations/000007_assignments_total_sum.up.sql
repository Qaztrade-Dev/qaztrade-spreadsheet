begin;

alter table "assignments"
    drop column "rows_from",
    drop column "rows_until"
;

alter table "assignments"
  add column "total_rows" bigint default 0,
  add column "total_sum" bigint default 0
;

CREATE INDEX applications_is_signed_idx ON applications (is_signed);

DROP INDEX applications_spreadsheet_id_idx;
CREATE INDEX applications_spreadsheet_id_idx ON applications (spreadsheet_id);

commit;
