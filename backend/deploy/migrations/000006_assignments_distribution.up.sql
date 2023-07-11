begin;

--- Functions

CREATE FUNCTION safe_cast_to_int(text) RETURNS bigint
LANGUAGE plpgsql AS $$
DECLARE
    result bigint;
BEGIN
    BEGIN
        result := $1::bigint;
    EXCEPTION WHEN others THEN
        result := 0;
    END;
    RETURN result;
END; $$;


CREATE FUNCTION safe_cast_to_double(text) RETURNS DOUBLE PRECISION
LANGUAGE plpgsql AS $$
DECLARE
    result DOUBLE PRECISION;
BEGIN
    result := coalesce(
        nullif(
            REGEXP_REPLACE(
                REPLACE($1, ',', '.'),
                '[^0-9.]',
                '',
                'g'
            ),
            ''
        ),
        '0'
    ) :: DOUBLE PRECISION;
    RETURN result;
END; $$;


CREATE FUNCTION get_business_category_table(text, int, int)
RETURNS TABLE(
    spreadsheet_id text,
    sheet_title text,
    business_category text,
    expense double precision
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY 
    SELECT
        a.spreadsheet_id,
        s.value ->> 'title' as sheet_title,
        d.value ->> $2 as business_category,
        safe_cast_to_double(d.value ->> $3) as expense
    FROM
        applications a,
        jsonb_array_elements(a.attrs -> 'sheets') as s,
        jsonb_array_elements(s.value -> 'data') as d
    where
        a.is_signed
        and s.value ->> 'title' = $1
        and d.value ->> $2 <> '';
END; $$;

CREATE FUNCTION get_business_category_table_agg(view_name text)
RETURNS TABLE (
    spreadsheet_id text,
    sheet_title text,
    expenses_sum double precision
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY EXECUTE format(
    'SELECT
        e.spreadsheet_id,
        e.sheet_title,
        SUM(e.expense * b.coef) AS expenses_sum
    FROM
        %I e
    JOIN
        business_category_values_view b ON b.category = e.business_category
    GROUP BY
        e.spreadsheet_id, e.sheet_title', view_name);
END; $$;

--- Functions

CREATE VIEW business_category_values_view AS
select
    *
from
    (
        values
            ('крупный', 0.3),
            ('средний', 0.5),
            ('малый', 0.6)
    ) as temp_table(category, coef);


CREATE VIEW tnved_values_view AS
select
    *
from
    (
        values
            ('нижний', 1, 0.3),
            ('средний', 2, 0.5),
            ('высокий', 3, 0.8)
    ) as temp_table(level, index, coef);


CREATE VIEW tnved_view AS
select
    t.*,
    v.index,
    v.coef
from
    tnved t
    join tnved_values_view v on v.level = t.level;

CREATE VIEW logistics_values_view AS
SELECT
    title,
    coef
FROM
    (
        VALUES
            ('да', 1.05),
            ('нет', 1.0)
    ) AS temp_table(title, coef);

CREATE VIEW applicants_info_view AS
SELECT
    id,
    spreadsheet_id,
    least(
        safe_cast_to_int(attrs -> 'application' ->> 'tax_sum'),
        75000 * 3450
    ) as tax_sum
FROM
    applications
where
    is_signed;


--- Затраты на доставку транспортом

CREATE VIEW expenses_dostavka_view AS
SELECT
    a.spreadsheet_id,
    s.value ->> 'title' as sheet_title,
    split_part(d.value ->> 58, ' - ', 1) as tnved_code,
    coalesce(nullif(d.value ->> 37, ''), 'нет') as logistics_type,
    safe_cast_to_double(d.value ->> 124) as expense
FROM
    applications a,
    jsonb_array_elements(a.attrs -> 'sheets') as s,
    jsonb_array_elements(s.value -> 'data') as d
where
    a.is_signed
    and s.value ->> 'title' = 'Затраты на доставку транспортом'
    and d.value ->> 58 <> '';

--- agg

CREATE VIEW expenses_dostavka_view_agg AS
select
    e.spreadsheet_id,
    e.sheet_title,
    sum(e.expense * t.coef * n.coef) as expenses_sum,
    min(t.index) as tnved_index
from
    expenses_dostavka_view e
    join tnved_view t on t.code = e.tnved_code
    join logistics_values_view n on n.title = e.logistics_type
group by
    (e.spreadsheet_id, e.sheet_title);

--- Затраты на доставку транспортом

CREATE VIEW expenses_certifikazia_predpriyatia_view AS
select * from get_business_category_table('Затраты на сертификацию предприятия', 56, 58);

CREATE VIEW expenses_certifikazia_predpriyatia_view_agg AS
select * from get_business_category_table_agg('expenses_certifikazia_predpriyatia_view');


CREATE VIEW expenses_reklama_iku_view AS
select * from get_business_category_table('Затраты на рекламу ИКУ за рубежом', 57, 59);

CREATE VIEW expenses_reklama_iku_view_agg AS
select * from get_business_category_table_agg('expenses_reklama_iku_view');


CREATE VIEW expenses_perevod_iku_view AS
select * from get_business_category_table('Затраты на перевод каталога ИКУ', 55, 57);

CREATE VIEW expenses_perevod_iku_view_agg AS
select * from get_business_category_table_agg('expenses_perevod_iku_view');


CREATE VIEW expenses_arenda_iku_view AS
select * from get_business_category_table('Затраты на аренду помещения ИКУ', 52, 54);

CREATE VIEW expenses_arenda_iku_view_agg AS
select * from get_business_category_table_agg('expenses_arenda_iku_view');


CREATE VIEW expenses_certifikazia_iku_view AS
select * from get_business_category_table('Затраты на сертификацию ИКУ', 58, 60);

CREATE VIEW expenses_certifikazia_iku_view_agg AS
select * from get_business_category_table_agg('expenses_certifikazia_iku_view');


CREATE VIEW expenses_demonstrazia_iku_view AS
select * from get_business_category_table('Затраты на демонстрацию ИКУ', 52, 54);

CREATE VIEW expenses_demonstrazia_iku_view_agg AS
select * from get_business_category_table_agg('expenses_demonstrazia_iku_view');


CREATE VIEW expenses_franchaizing_view AS
select * from get_business_category_table('Затраты на франчайзинг', 52, 54);

CREATE VIEW expenses_franchaizing_view_agg AS
select * from get_business_category_table_agg('expenses_franchaizing_view');


CREATE VIEW expenses_registrazia_tovar_znakov_view AS
select * from get_business_category_table('Затраты на регистрацию товарных знаков', 64, 66);

CREATE VIEW expenses_registrazia_tovar_znakov_view_agg AS
select * from get_business_category_table_agg('expenses_registrazia_tovar_znakov_view');


CREATE VIEW expenses_arenda_view AS
select * from get_business_category_table('Затраты на аренду', 52, 54);

CREATE VIEW expenses_arenda_view_agg AS
select * from get_business_category_table_agg('expenses_arenda_view');


CREATE VIEW expenses_perevod_view AS
select * from get_business_category_table('Затраты на перевод', 55, 57);

CREATE VIEW expenses_perevod_view_agg AS
select * from get_business_category_table_agg('expenses_perevod_view');


CREATE VIEW expenses_reklaman_view AS
select * from get_business_category_table('Затраты на рекламу товаров за рубежом', 57, 59);

CREATE VIEW expenses_reklaman_view_agg AS
select * from get_business_category_table_agg('expenses_reklaman_view');


CREATE VIEW expenses_uchastie_vystavka_view AS
select * from get_business_category_table('Затраты на участие в выставках', 71, 73);

CREATE VIEW expenses_uchastie_vystavka_view_agg AS
select * from get_business_category_table_agg('expenses_uchastie_vystavka_view');


CREATE VIEW expenses_uchastie_vystavka_iku_view AS
select * from get_business_category_table('Затраты на участие в выставках ИКУ', 71, 73);

CREATE VIEW expenses_uchastie_vystavka_iku_view_agg AS
select * from get_business_category_table_agg('expenses_uchastie_vystavka_iku_view');


--- Затраты на соответствие товаров требованиям

CREATE VIEW expenses_sootvetstvie_tovara_view AS
SELECT
    a.spreadsheet_id,
    s.value ->> 'title' as sheet_title,
    split_part(d.value ->> 62, ' - ', 1) as tnved_code,
    safe_cast_to_double(d.value ->> 64) as expense
FROM
    applications a,
    jsonb_array_elements(a.attrs -> 'sheets') as s,
    jsonb_array_elements(s.value -> 'data') as d
where
    a.is_signed
    and s.value ->> 'title' = 'Затраты на соответствие товаров требованиям'
    and d.value ->> 62 <> '';

--- agg

CREATE VIEW expenses_sootvetstvie_tovara_view_agg AS
select
    e.spreadsheet_id,
    e.sheet_title,
    sum(e.expense * t.coef) as expenses_sum,
    min(t.index) as tnved_index
from
    expenses_sootvetstvie_tovara_view e
    join tnved_view t on t.code = e.tnved_code
group by
    (e.spreadsheet_id, e.sheet_title);

--- Затраты на соответствие товаров требованиям

commit;
