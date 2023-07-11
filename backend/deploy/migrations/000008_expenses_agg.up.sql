CREATE VIEW expenses_agg AS
with expenses_tmp as(
select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_dostavka_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_certifikazia_predpriyatia_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_reklama_iku_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_perevod_iku_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_arenda_iku_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_certifikazia_iku_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_demonstrazia_iku_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_franchaizing_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_registrazia_tovar_znakov_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_arenda_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_perevod_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_reklaman_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_uchastie_vystavka_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_uchastie_vystavka_iku_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id

union

select
	a.id,
	a.attrs->'application'->>'bin' as bin,
    a.attrs->'application'->>'from' as company,
    a.attrs->'application'->>'fact_addr' as fact_addr,
    a.attrs->'application'->>'industry' as industry,
	e.sheet_title,
	safe_cast_to_int(s.value->>'rows') as total_rows,
	least(e.expenses_sum, info.tax_sum) as total_sum
from applications a
cross join jsonb_array_elements(a.attrs -> 'sheets') as s
join expenses_sootvetstvie_tovara_view_agg e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'
join applicants_info_view info on info.id = a.id
)
select * from expenses_tmp
