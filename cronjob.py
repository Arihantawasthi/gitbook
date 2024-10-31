import subprocess
import os

GIT_DIR = os.path.expanduser("~/Documents/projects/gitbook-repos/")

def main():
    data = {}
    data["num_of_files"] = 0
    data["num_of_lines"] = 0
    data["num_of_commits"] = 0
    repos = get_all_repos()
    data["num_of_repos"] = len(repos)

    for repo in repos:
        repo_path = GIT_DIR+repo
        result = subprocess.run(
            ["git", "--git-dir", repo_path, "ls-tree", "-r", "HEAD", "--name-only"],
            capture_output=True,
            text=True,
            check=True
        )
        files = result.stdout.splitlines()

        data["num_of_files"] += len(files)
        data["num_of_lines"] += num_of_lines(files, repo_path)
        data["num_of_commits"] += num_of_commits(repo_path)

    return data

def num_of_commits(repo_path):
    result = subprocess.run(
        ["git", "--git-dir", repo_path, "rev-list", "--all", "--count"],
        capture_output=True,
        text=True,
        check=True
    )
    total_commits = int(result.stdout.strip())
    return total_commits

def get_all_repos():
    result = subprocess.run(
        ["ls", GIT_DIR],
        capture_output=True,
        text=True,
        check=True
    )
    return result.stdout.splitlines()

def num_of_lines(file_names, repo_path):
    total_lines = 0
    for name in file_names:
        file_content = subprocess.run(
            ["git", "--git-dir", repo_path, "cat-file", "-p", f"HEAD:{name}"],
            capture_output=True,
            text=True,
            errors="ignore", # ignore decoding errors
            check=True
        )
        total_lines += len(file_content.stdout.splitlines())

    return total_lines

if __name__ == "__main__":
    data = main()
    print(data)
