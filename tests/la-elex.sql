.load dist/dom.so
.bail on
.headers on 
.mode csv

.param init

insert into sqlite_parameters(key, value) 
  values (':la_election_results', readfile('./tests/data/la-election-results.html'));
