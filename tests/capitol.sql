.load dist/dom.so
.bail on
.headers on 
.mode column

.param init

insert into sqlite_parameters(key, value) 
  values (':capitol_breach_cases', readfile('./tests/data/capitol-breach-cases.html'));

create table defendant_rows_html as 
  select 
    assert(dom_count( dom_table(html), 'td:nth-child(2) a') == 1)
      as assertions, 
    * 
  from dom_$$(:capitol_breach_cases, 'tbody tr');

create table defendants (
  case_number,
  name unique,
  name_url,
  charges,
  documents_html,
  location,
  status,
  last_updated
);

insert into defendants
  select 
    dom_trim(                                                                         
      dom_$text( dom_table(html), 'td:nth-child(1)')
    ) as case_number,
    dom_trim(                                                                         
      dom_$text( dom_table(html), 'td:nth-child(2)')
    ) as name,
    dom_attr_get(dom_table(html), 'td:nth-child(2) a', 'href')
      as name_url,
    dom_trim(                                                                         
      dom_$text( dom_table(html), 'td:nth-child(3)')
    ) as charges,                                                                     
    dom_$( dom_table(html), 'td:nth-child(4)')
      as documents_url,
    dom_trim(                                                                         
      dom_$text( dom_table(html), 'td:nth-child(5)')
    ) as location,
    dom_trim(                                                                         
      dom_$text( dom_table(html), 'td:nth-child(6)')
    ) as status,
    dom_trim(                                                                         
      dom_$text( dom_table(html), 'td:nth-child(7)')
    ) as last_updated
  from defendant_rows_html;


create table documents(
  defendant,
  name,
  url,
  foreign key (defendant) references defendants(name)
);

insert into documents
  select 
    dom_trim(                                                                         
      dom_$text( dom_table(defendant_rows_html.html), 'td:nth-child(2)')
    ) as defendant,
    documents.text as name, 
    dom_attr_get(documents.html, 'a', 'href') as url
  from defendant_rows_html, 
    dom_$$( dom_table(defendant_rows_html.html), "td:nth-child(4) a") as documents;


select * from documents;

.exit 1;