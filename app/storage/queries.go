package storage

import (
	"fmt"
	"gitbook/app/types"
	"time"
)

func GetLatestStats() (types.AggStats, error) {
	row := DBConn.QueryRow(
		`SELECT
            num_of_lines, num_of_commits,
            num_of_files, num_of_repos
        FROM stats WHERE date = (SELECT MAX(date) FROM stats) LIMIT 1;`,
	)

	var stats types.AggStats
	if err := row.Scan(
		&stats.NumOfLines, &stats.NumOfCommits,
		&stats.NumOfFiles, &stats.NumOfRepos,
	); err != nil {
		return stats, err
	}

	return stats, nil
}

func GetStatsForADate(date time.Time) (types.AggStats, error) {
    row := DBConn.QueryRow(
        `SELECT
            num_of_lines, num_of_commits, num_of_files, num_of_repos
        FROM stats WHERE date = $1;
        `,
        date,
    )
    var stats types.AggStats
    err := row.Scan(
		&stats.NumOfLines, &stats.NumOfCommits,
		&stats.NumOfFiles, &stats.NumOfRepos,
    )
    if err != nil {
        return stats, err
    }

    return stats, nil
}

func GetRepos(limit, offset int) ([]types.RepoDetails, error) {
    rows, err:= DBConn.Query(
        `SELECT
            name, description, is_pinned, default_branch, author, created_at, last_commit_at
        FROM repos WHERE visibility = 'public' ORDER BY is_pinned DESC, last_commit_at DESC
        LIMIT $1 OFFSET $2;
        `,
        limit, offset,
    )
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var repos []types.RepoDetails
    for rows.Next() {
        var repo types.RepoDetails
        err := rows.Scan(
            &repo.Name, &repo.Desc, &repo.IsPinned, &repo.DefaultBranch,
            &repo.Author, &repo.CreatedAt, &repo.LastCommitAt,
        )
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        repos = append(repos, repo)
    }

    return repos, nil
}

func UpdateLastCommit(repoName, authorName string, lastCommitAt time.Time) error {
    query := `
        UPDATE repos
        SET last_commit_at = $1
        WHERE name = $2 and author = $3;
    `

    _, err := DBConn.Exec(query, lastCommitAt, repoName, authorName)

    return err
}
