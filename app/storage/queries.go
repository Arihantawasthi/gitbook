package storage

import "gitbook/app/types"

func GetStats() ([]types.AggStats, error) {
	rows, err := DBConn.Query(
		`SELECT
            SUM(num_of_lines), SUM(num_of_commits),
            SUM(num_of_files), SUM(num_of_repos)
        FROM stats`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statsLift []types.AggStats
	for rows.Next() {
		var stats types.AggStats
		if err := rows.Scan(&stats.NumOfLines, &stats.NumOfCommits,
			&stats.NumOfFiles, &stats.NumOfRepos,
		); err != nil {
			return nil, err
		}
        statsLift = append(statsLift, stats)
	}
    if err := rows.Err(); err != nil {
        return nil, err
    }
	return statsLift, nil
}
