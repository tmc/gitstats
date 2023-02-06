# gitstats

gitstats is a tool to analyze git repositories.

This software is in early stages. Currently it produces CSV files that should be loaded into duckdb
for analysis.


## Samples Queries

Show the top 20 files by number of edits:
```sql
WITH changes as (
    select count(*) n, Filename from read_csv('git_history.csv', header=True, auto_detect=True) GROUP BY Filename
)
select * from changes order by n desc limit 20;
```
