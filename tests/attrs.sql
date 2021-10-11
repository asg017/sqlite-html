.load dist/dom.so
.bail on

.param init


insert into sqlite_parameters(key, value) 
values 
  (':a', '<a href="https://google.com"> hello</a>');


select dom_attr_get(:a, 'a', 'href');

.exit 1