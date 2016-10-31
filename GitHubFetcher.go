package repoclone

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/oauth2"
)

var debug bool

// SetDebug turns on debugging output on this library
func SetDebug(d bool) {
	debug = d
}

// ListPRs will list all the pull requests for an organization
func ListPRs(orgname string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repo", "Created", "Author", "Title", "Link"})

	client := getClient()
	allRepos := getAllRepos(orgname)
	for _, repo := range allRepos {
		if debug {
			log.Printf("Analying repo: %s", *repo.Name)
		}

		opt := &github.PullRequestListOptions{State: "open", Direction: "asc"}
		owner := "AKQASF"
		prs, _, err := client.PullRequests.List(owner, *repo.Name, opt)
		check(err)
		if debug {
			log.Printf("Number of PRs found: %d", len(prs))
		}
		for _, pr := range prs {
			table.Append(formatPR(pr, repo))
		}
	}
	table.Render()
}

func formatPR(pr *github.PullRequest, repo github.Repository) []string {
	var formatedTime string

	warnAfter := time.Now().AddDate(0, 0, -3)
	errorAfter := time.Now().AddDate(0, 0, -7)

	format := "2006-01-02"

	if pr.CreatedAt.Before(errorAfter) {
		formatedTime = color.RedString(pr.CreatedAt.Format(format))
	} else if pr.CreatedAt.Before(warnAfter) {
		formatedTime = color.YellowString(pr.CreatedAt.Format(format))
	} else {
		formatedTime = color.GreenString(pr.CreatedAt.Format(format))
	}

	//fmt.Printf("%s %s %s %s %s\n", *repo.Name, formatedTime, *pr.User.Login, *pr.Title, *pr.URL)
	return []string{*repo.Name, formatedTime, *pr.User.Login, *pr.Title, *pr.URL}
}

// Sync pull down all repos for an orgnaization
func Sync(orgname string, path string, clone bool, update bool, remove bool) {
	allRepos := getAllRepos(orgname)
	allReposMap := make(map[string]bool)
	for _, repo := range allRepos {
		if !repoExistsLocally(repo, path) {
			if clone {
				doClone(repo, path)
			}
		} else if update {
			doUpdate(repo, path)
		}
		allReposMap[*repo.Name] = true
	}
	if remove {
		cleanup(path, allReposMap)
	}
}

func cleanup(path string, repos map[string]bool) {
	files := getDirectories(path)

	for _, directory := range files {
		if !directory.IsDir() {
			if debug {
				log.Printf("Skipping %s as it is not a directory", directory.Name())
			}
			continue
		}
		// check if the directory exists in the map
		if _, b := repos[directory.Name()]; !b {
			if debug {
				log.Printf("Removing %s", directory.Name())
			}
			os.RemoveAll(filepath.Join(path, directory.Name()))
		} else if debug {
			log.Printf("Skipping %s as it is found in the organization", directory.Name())
		}
	}
}

func getDirectories(path string) []os.FileInfo {

	// Get a list of directories off this
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Error reading the directory: " + path)
	}
	return files
}

func getAllRepos(orgname string) []github.Repository {
	client := getClient()

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(orgname, opt)
		if err != nil {
			return nil
		}
		for _, repo := range repos {
			allRepos = append(allRepos, *repo)
		}
		//allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	log.Printf("Found %d repo(s) for the organization %s", len(allRepos), orgname)
	return allRepos
}

func getClient() *github.Client {
	var tc *http.Client
	envToken := os.Getenv("GITHUB_TOKEN")
	if len(envToken) > 0 {
		token := oauth2.Token{AccessToken: envToken}
		ts := oauth2.StaticTokenSource(&token)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}
	client := github.NewClient(tc)
	return client
}

func repoExistsLocally(repo github.Repository, path string) bool {
	fullPath := filepath.Join(path, *repo.Name)
	_, err := os.Stat(fullPath)
	return err == nil
}

func doUpdate(repo github.Repository, path string) (int, error) {
	if debug {
		log.Printf("Updating %s", *repo.Name)
	}
	directory := filepath.Join(path, *repo.Name)
	return run(directory, "git", "pull")
}

func doClone(repo github.Repository, path string) (int, error) {
	if debug {
		log.Printf("Cloning %s (%s)", *repo.Name, *repo.SSHURL)
	}
	return run(path, "git", "clone", *repo.SSHURL)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
