begin;

drop table "assignment_messages";

alter table "assignments"
    drop column "resolved_at",
    drop column "countdown_duration",
    drop column "resolution_status_id"
;

commit;