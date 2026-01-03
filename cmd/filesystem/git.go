package filesystem

import (
	"github.com/go-git/go-git/v5"
)

// Using go-git because it doesn't require git binary
// TODO: Traverse path and find out if have .git
func IsPathInGitRepo(pathInfo *PathInfo) (bool, error) {
	// DetectDotGit: Allow path to be files or folders
	_, err := git.PlainOpenWithOptions(pathInfo.path, &git.PlainOpenOptions{DetectDotGit: true})
	if err == nil {
		return true, nil
	}

	if err == git.ErrRepositoryNotExists {
		return false, nil
	}

	return false, err
}
