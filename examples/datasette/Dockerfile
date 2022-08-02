FROM datasetteproject/datasette

RUN apt-get update && apt-get install -y \
    curl sqlite3 \
  && rm -rf /var/lib/apt/lists/*

# Download pre-build sqlite-html binary for linux (.so suffix)
ADD https://github.com/asg017/sqlite-html/releases/download/v0.1.0/html0.so /html0.so

# download some example webpages, save to a a SQLite database on disk
RUN mkdir -p /data \
  && curl -o /data/ca.gov.html https://www.ca.gov \
  && curl -o /data/whitehouse.gov.html https://www.whitehouse.gov/ \
  && curl -o /data/text.npr.org.html https://text.npr.org/ \
  && curl -o /data/blog.pythonlibrary.org.html https://blog.pythonlibrary.org/ \
  && curl -o /data/lacounty.html https://www.lavote.gov/home/voting-elections/current-elections/election-results/past-election-results

RUN sqlite3 /data/data.db 'create table pages as \
select \
  "ca.gov" as name, \
  cast(column2 as text) as page \
from(values \
  ("ca.gov", readfile("/data/ca.gov.html")), \
  ("whitehouse.gov", readfile("/data/whitehouse.gov.html")), \
  ("text.npr.org", readfile("/data/text.npr.org.html")), \
  ("lacounty-results", readfile("/data/lacounty.html")) \
)'

# Run datasette instance on that saved DB, and load the sqlite-html extension
CMD datasette /data/data.db \
  -p 8080 -h 0.0.0.0 --cors --load-extension=/html0 
