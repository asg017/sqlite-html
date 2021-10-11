.load dist/dom.so
.bail on
.headers on 
.mode csv

.param init

insert into sqlite_parameters(key, value) 
  values (':ca_data_breaches', readfile('./tests/data/ca-databreach.html'));


with breaches as (select 
  dom_trim(dom_$text(dom_table(html), 'td:nth-child(1)')) as organization,
  dom_attr_get(dom_table(html), 'td:nth-child(1) a', 'href') as report_url,
  dom_trim(dom_$text(dom_table(html), 'td:nth-child(2)')) as breach_dates,
  dom_trim(dom_$text(dom_table(html), 'td:nth-child(3)')) as report_date

from dom_$$(:ca_data_breaches, 'tbody tr'))
select * from breaches;

.exit 1;