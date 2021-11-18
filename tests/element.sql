.load ./dist/html0.so

select 1, html_element('p', null, 'hello', html_element("b", null, "Alex"));

select 1, html_element('b', null, '<script>escape me</script>');

select 1, html_element('p', null, html('<b>dude'));

select 
  html_element('html', null, 
    html_element('head', null, 
      html_element('meta', html_attributes(json_object('charset', 'utf-8'))), 
      html_element('title', null, 'My test page')
    ),
    html_element('body', null,
      html_element('img', html_attributes(json_object(
        'src', 'images/firefox-icon.png', 
        'alt', 'My test image'
        )))
    )
  );

.exit 0;