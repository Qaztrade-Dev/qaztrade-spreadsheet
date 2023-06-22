begin;

delete from "user_roles" where "value" in (
    'digital',
    'finance',
    'legal'
);
drop table "user_role_bindings";
drop table "assignment_results";
drop table "assignments";
drop table "batch_applications";
drop table "batches";

commit;
