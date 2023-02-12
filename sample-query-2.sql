WITH changes as (
    select count(*) n, filename from read_csv('git_history.csv', header=True, auto_detect=True) GROUP BY filename
)
select * from changes order by n desc limit 20;

