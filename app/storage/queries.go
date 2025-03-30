package storage

import (
	"gitbook/app/types"
)

func GetStats() (*types.AggStats, error) {
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
		return nil, err
	}

	return &stats, nil
}

func GetRepos(limit int) (*types.RepoDetails, error) {
    row := DBConn.QueryRow(
        `SELECT
            name, description, is_pinned, default_branch, author, created_at, last_commit_at
        FROM repos WHERE visibility = "public" ORDER BY is_pinned DESC, last_commit_at DESC;
        `,
    )
    var repos types.RepoDetails

    if err := row.Scan(
        &repos.Name, &repos.Desc, &repos.IsPinned, &repos.DefaultBranch,
        &repos.Author, &repos.CreatedAt, &repos.LastCommitAt,
    ); err != nil {
        return nil, err
    }

    return  &repos, nil
}
