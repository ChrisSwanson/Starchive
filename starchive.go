package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

// Starchive maintains configuration attributes to be used and passed throughout
// the gathering and downloading states of determining starred repos, which
// directory and github token to be used.
type Starchive struct {

	// logger to control output and log leves of the client.
	logger *zerolog.Logger

	// Dir, the directory to archive git repos in.
	Dir string

	// Repos, the list of repository (names and urls) to clone/pull after
	// gathering the starred repositories
	Repos []Repo

	// Token, the github user access token used to identify the user and access
	// the github api for starred repos with.
	Token string
}

// Repo holds specific attributes to be utilized in the Starchive struct - the
// name of the repository, as well as the cloning url
type Repo struct {

	// Name is the name of the repo
	Name string

	// CloneURL is the url in which to clone from
	CloneURL string
}

// NewStarchive creates and returns a new, empty` starchive struct.
func NewStarchive() Starchive {

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	return Starchive{
		logger: &log,
		Dir:    "",
		Repos:  []Repo{},
		Token:  "",
	}
}

func (s *Starchive) cloneRepos() {

	// If the directory doesn't exist, create it.
	if _, err := os.Stat(s.Dir); os.IsNotExist(err) {
		os.MkdirAll(s.Dir, 0744)
	}

	var repoDir string

	// for each starred repo in the list, git clone or pull
	for _, repo := range s.Repos {

		// s.Dir base directory + the name of the repo
		repoDir = filepath.Join(s.Dir, repo.Name)

		if _, err := os.Stat(repoDir); os.IsNotExist(err) {
			// If the repo directory does not exist git clone the repo
			s.logger.Debug().Str("CloneURL", repo.CloneURL).Str("Name", repo.Name).Msg("cloning repo")

			// clone the repo to the repoDir
			_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
				URL: repo.CloneURL,
			})

			if err != nil {
				s.logger.Warn().Err(err).Str("CloneURL", repo.CloneURL).Msgf("error cloning repo")
			}

		} else {
			// Else if the repo directory already exists, git pull the repo
			s.logger.Debug().Str("CloneURL", repo.CloneURL).Str("Name", repo.Name).Msg("pulling repo")

			// open the git directory for the repo directory
			r, err := git.PlainOpen(repoDir)
			if err != nil {
				s.logger.Warn().Err(err).Str("CloneURL", repo.CloneURL).Msg("error opening repo")
			}

			// gather the working tree for the git repo
			workTree, err := r.Worktree()
			if err != nil {
				s.logger.Warn().Err(err).Str("CloneURL", repo.CloneURL).Msg("error with working tree")
			}

			// git pull
			if err = workTree.Pull(&git.PullOptions{}); err != nil {
				s.logger.Warn().Err(err).Str("CloneURL", repo.CloneURL).Msg("error pulling from git repo")
			}

		}

	}

}

func (s *Starchive) getStarred() {

	// the default results per page is already 30, but in case we want to modify
	// this we can change the value here.
	var resultsPerPage int = 30

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// utilize the oauth2 client for github calls
	client := github.NewClient(tc)

	// configure the list starred options with the results per page attribute
	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: resultsPerPage},
	}

	// iterate over the github paginated results until there is no more results
	for {

		//
		repos, resp, err := client.Activity.ListStarred(ctx, "", opt)
		if err != nil {
			log.Fatal(err)
		}

		// append the relevant repository information to the repo slice
		for _, repo := range repos {
			s.logger.Debug().Str("CloneURL", *repo.Repository.CloneURL).Msg("found repo")

			s.Repos = append(s.Repos, Repo{
				Name:     *repo.Repository.Name,
				CloneURL: *repo.Repository.CloneURL,
			})

		}

		// if there is no more pages, break
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage

	}

}

// Run determines if the appropriate variables are provided, and set to runtime
// default values, gathers the starred repository list and git clones or pulls
// (depending on current state on disk).
func (s *Starchive) Run() {

	s.logger.Debug().Msg("test")

	// if the token provided is nil, this won't work.  Exit fatally.
	if s.Token == "" {
		s.logger.Fatal().Msg("no github user access token provided")
	}

	// if the directory is nil, use the current directory
	if s.Dir == "" {
		path, err := os.Getwd()

		if err != nil {
			s.logger.Fatal().Err(err).Msg("no valid directory provided")
		}

		s.Dir = path
	}

	s.getStarred()
	s.cloneRepos()

}
