# github-org-prs

CLI tool to get info on open PRs for a given org


# Flags

- `--org` Github organization shortname
- `--key` Github API key
- `--format`, `-f` (_`table`, `json`_) The format to print to screen
- `--orderby`, `-o` (`UpdatedAt__desc`) Order the results as a list columnName_asc|desc
- `--columns`, `-c`: (`URL`,`Approved`) List of columns to return

# Column Names

- `RepoName`
- `URL`
- `CreatedAt`
- `UpdatedAt`
- `Author`
- `TotalReviews`
- `Approved`
- `ChangesRequested`
- `Reviewers`


# Fitering Options

TODO
