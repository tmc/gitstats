.mode json
with ranked_files_by_commits as (
    select
        filename,
        count(sha) as total_commits,
        rank() over (order by count(sha) desc) as all_time_rank,
        sum(case when date > current_date - interval '1 year' then 1 else 0 end) as recent_commits,
        rank() over (order by sum(case when date > current_date - interval '1 year' then 1 else 0 end) desc) as recent_rank
    from read_csv('git_history.csv', header = true, auto_detect = true)
    group by filename
)

select all_time_rank, recent_rank, total_commits, recent_commits, filename
from ranked_files_by_commits
