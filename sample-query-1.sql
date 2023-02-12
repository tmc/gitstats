.mode markdown
WITH ranked_files AS (
  SELECT filename, COUNT(sha) AS total_commits,
         RANK() OVER (ORDER BY COUNT(sha) DESC) AS rank
  FROM read_csv('git_history.csv', header=True, auto_detect=True)
  WHERE date >= '2022-01-01'
  GROUP BY filename
)
SELECT rank, total_commits, filename
FROM ranked_files
WHERE rank <= 20;
