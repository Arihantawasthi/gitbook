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
