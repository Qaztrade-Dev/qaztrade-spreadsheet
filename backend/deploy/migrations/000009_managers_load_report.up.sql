begin;

CREATE VIEW managers_load_report AS
with
managers as (
	select
		u.id,
		u.attrs->>'full_name' as full_name
	from users u
	join user_role_bindings urb on urb.user_id = u.id
	join user_roles ur on ur.id = urb.role_id
	where 
		ur.value = 'manager'
)
select
	m.id,
	m.full_name,
	a.type,
	count(*) as total_rows,
	sum(a.total_rows) as total_rows_sum,
	sum(a.total_sum) as total_sum_sum
from managers m
join assignments a on a.user_id = m.id
group by m.id, m.full_name, a.type
order by m.id, a.type
;

commit;