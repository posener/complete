package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/posener/complete"
)

var modEditReplaceRegexp = regexp.MustCompile(strings.Replace(strings.Replace(strings.Replace(modE_R_RE, " ", "", -1), "\n", "", -1), "\t", "", -1))
var modE_R_RE = `(?:-replace=)?
(?:
	(?P<old>[^@=]+)
	(?:
		(?:(@)(?P<v_old>[^=]+))?
		(?:=
			(?:
				(?P<new>[^@]+)
				(?:(@)(?P<v_new>.+)?)?
			)?
		)?
	)? 
)?`

type replaceArgs struct {
	Old   string
	OldAt bool
	OldV  string
	Eq    bool
	New   string
	NewAt bool
	NewV  string
}

func (rs replaceArgs) String() string {
	s := rs.Old
	if rs.OldV != "" {
		s += "@" + rs.OldV
	}
	s += "=" + rs.New
	if rs.NewV != "" {
		s += "@" + rs.NewV
	}
	return s
}

func parseReplace(last string) (ret replaceArgs) {
	defer func() {
		if x := recover(); x != nil {
			ret = replaceArgs{}
		}
	}()
	if !modEditReplaceRegexp.MatchString(last) {
		panic("")
	}
	subs := modEditReplaceRegexp.FindAllStringSubmatch(last, 5)
	return replaceArgs{
		Old:   subs[0][1],
		OldAt: subs[0][2] != "",
		OldV:  subs[0][3],
		Eq:    strings.Contains(last, "="),
		New:   subs[0][4],
		NewAt: subs[0][5] != "",
		NewV:  subs[0][6],
	}
}

func sliceContains(src []string, match string) bool {
	for _, st := range src {
		if st == match {
			return true
		}
	}
	return false
}

// var modEditReplaceRegexp = regexp.MustCompile("-replace=[^@=]+?P<v_old>@[^=]+)=)?) )?")

// from go help mod edit
// The -replace=old[@v]=new[@v] and -dropreplace=old[@v] flags
// add and drop a replacement of the given module path and version pair.
// If the @v in old@v is omitted, the replacement applies to all versions
// with the old module path. If the @v in new@v is omitted, the new path
// should be a local module root directory, not a module path.
// Note that -replace overrides any existing replacements for old[@v].
func predictModEditReplace(args complete.Args) (prediction []string) {
	if sliceContains(args.Completed, "-replace") {
		cline := strings.Split(os.Getenv("COMP_LINE"), " ")
		rArgs := parseReplace(cline[len(cline)-1])
		log.Printf("XXXX %#+v\n", rArgs)
		if rArgs.Old == "" || (rArgs.OldV == "" && !rArgs.Eq) {
			return rArgs.old()
		}
		if /* rArgs.Old != "" && */ rArgs.Eq && !rArgs.NewAt {
			return rArgs.newPart()
		}
	}
	return complete.PredictFiles("go.mod").Predict(args)
}

func (rArgs *replaceArgs) old() []string {
	log.Printf("YYYY %#+v\n", rArgs)
	cmd := exec.Command("go", "mod", "edit", "-json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("go mod error %v\n", err)
		return nil
	}
	gomod := GoMod{}
	err = json.Unmarshal(out, &gomod)
	if err != nil {
		log.Printf("go mod error %v\n%v\n", string(out), err)
		return nil
	}
	req := make([]string, 0)
	for i := range gomod.Require {
		req = append(req, gomod.Require[i].Path+"=")
		if rArgs.OldV == "" {
			req = append(req, gomod.Require[i].Path+"@")
		}
	}
	return req
}

func (rArgs *replaceArgs) newPart() []string {
	if rArgs.New != "" {
		lastPath := rArgs.New
		usr, err := user.Current()
		if strings.HasPrefix(rArgs.New, "~") {
			if err != nil {
				return nil
			}
			lastPath = filepath.Join(usr.HomeDir, lastPath[2:])
		}
		fd, err := os.Open(lastPath)
		if err == nil {
			fi, err := fd.Stat()
			if err == nil {
				log.Printf("stat %+v\n", fi)
				if fi.IsDir() {
					log.Printf("is dir %+v\n", fd.Name())
					req := make([]string, 0)
					regexStr := fmt.Sprintf(`\n*module\s+%s(\s+//.*)?\n`, rArgs.Old)
					log.Printf("regexp %s\n", regexStr)
					re, err := regexp.Compile(regexStr)
					if err != nil {
						log.Printf("regex err %v\n", err)
					} else {
						err = filepath.Walk(fd.Name(), func(path string, info os.FileInfo, err error) error {
							if strings.HasSuffix(path, "go.mod") {
								if gomodF, err := ioutil.ReadFile(path); err != nil {
									log.Printf("go.mod read err %v\n", err)
								} else {
									if re.Match(gomodF) {
										if strings.HasPrefix(rArgs.New, "~") {
											gomodP := "~" + path[len(usr.HomeDir):len(path)-len("/go.mod")]
											log.Printf("return %v\n", gomodP)
											req = append(req, gomodP)
										} else {
											gomodP := path[:len(path)-len("go.mod")]
											req = append(req, gomodP)
											if strings.HasPrefix(rArgs.New, "../") {
												if strings.HasSuffix(rArgs.New, "/.") {
													req = append(req, rArgs.New+"./")
													req = append(req, rArgs.New+"./../")
												} else if strings.HasSuffix(rArgs.New, "/.") {
													req = append(req, rArgs.New+"/")
													req = append(req, rArgs.New+"/../")
												} else {
													req = append(req, rArgs.New+"../")
													req = append(req, rArgs.New+"../../")
												}
											}
										}
										log.Printf("found go.mod %s\n", path)
									} else {
										// log.Printf("ignoring go.mod %s\n", path)
									}
								}
							}
							return nil
						})
						if len(req) != 0 {
							return req
						}
					}
				}
			} else {
				log.Printf("stat err %+v\n", err)
			}
		} else {
			log.Printf("open err %+v\n", err)
		}
	}
	log.Printf("lastcompleted %+v\n", rArgs)
	if strings.HasPrefix(rArgs.New, "../") {
		return []string{rArgs.New + "../", rArgs.New + "../../"}
	}
	if strings.HasPrefix(rArgs.Old, "github.com/") {
		cl := http.Client{Timeout: 3 * time.Second}
		repo := rArgs.Old[len("github.com/"):]
		ppid := os.Getppid()
		fkd := forked{}
		tempf := fmt.Sprintf("%s/go-gomod-replace-dsts-%d", os.TempDir(), ppid)
		log.Printf("tempf %s\n", tempf)
		tf, err := os.Open(tempf)
		if err == nil {
			log.Printf("opened cache\n")
			if fi, err := tf.Stat(); err == nil {
				if !fi.ModTime().Add(time.Minute).After(time.Now()) {
					cache, err := ioutil.ReadFile(tempf)
					if err == nil {
						err = json.Unmarshal(cache, &fkd)
						if err != nil {
							log.Printf("cache err %v\n", err)
						}
						log.Printf("got from cache\n")
					}
				} else {
					log.Printf("cache stale %v\n", err)
				}
			} else {
				log.Printf("stat error %v\n", err)
			}
		} else {
			log.Printf("cache read error %v\n", err)
		}
		if len(fkd) == 0 {
			log.Printf("getting from api\n")
			resp, err := cl.Get(fmt.Sprintf("https://api.github.com/repos/%s/forks", repo))
			if err != nil {
				log.Printf("github api call err %v\n", err)
				return nil
			}
			// bod := resp.Body()
			de := json.NewDecoder(resp.Body)
			err = de.Decode(&fkd)
			if err != nil {
				log.Printf("github api call json err %v\n", err)
				return nil
			}
			data, _ := json.Marshal(&fkd)
			ioutil.WriteFile(fmt.Sprintf("%s/go-gomod-replace-dsts-%d", os.TempDir(), ppid), data, os.ModePerm)
		}
		log.Printf("fkd %v\n", fkd)
		req := make([]string, len(fkd))
		for i := range fkd {
			req[i] = "github.com/" + fkd[i].FullName + "@master"
		}
		return req
	}
	return nil
}

type Module struct {
	Path    string
	Version string
}

type GoMod struct {
	Module  Module
	Require []Require
	Exclude []Module
	Replace []Replace
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

type forked []struct {
	// ID       int    `json:"id"`
	// NodeID   string `json:"node_id"`
	// Name     string `json:"name"`
	FullName string `json:"full_name"`
	// Private  bool   `json:"private"`
	// Owner    struct {
	// 	Login             string `json:"login"`
	// 	ID                int    `json:"id"`
	// 	NodeID            string `json:"node_id"`
	// 	AvatarURL         string `json:"avatar_url"`
	// 	GravatarID        string `json:"gravatar_id"`
	// 	URL               string `json:"url"`
	// 	HTMLURL           string `json:"html_url"`
	// 	FollowersURL      string `json:"followers_url"`
	// 	FollowingURL      string `json:"following_url"`
	// 	GistsURL          string `json:"gists_url"`
	// 	StarredURL        string `json:"starred_url"`
	// 	SubscriptionsURL  string `json:"subscriptions_url"`
	// 	OrganizationsURL  string `json:"organizations_url"`
	// 	ReposURL          string `json:"repos_url"`
	// 	EventsURL         string `json:"events_url"`
	// 	ReceivedEventsURL string `json:"received_events_url"`
	// 	Type              string `json:"type"`
	// 	SiteAdmin         bool   `json:"site_admin"`
	// } `json:"owner"`
	// HTMLURL          string      `json:"html_url"`
	// Description      string      `json:"description"`
	// Fork             bool        `json:"fork"`
	// URL              string      `json:"url"`
	// ForksURL         string      `json:"forks_url"`
	// KeysURL          string      `json:"keys_url"`
	// CollaboratorsURL string      `json:"collaborators_url"`
	// TeamsURL         string      `json:"teams_url"`
	// HooksURL         string      `json:"hooks_url"`
	// IssueEventsURL   string      `json:"issue_events_url"`
	// EventsURL        string      `json:"events_url"`
	// AssigneesURL     string      `json:"assignees_url"`
	// BranchesURL      string      `json:"branches_url"`
	// TagsURL          string      `json:"tags_url"`
	// BlobsURL         string      `json:"blobs_url"`
	// GitTagsURL       string      `json:"git_tags_url"`
	// GitRefsURL       string      `json:"git_refs_url"`
	// TreesURL         string      `json:"trees_url"`
	// StatusesURL      string      `json:"statuses_url"`
	// LanguagesURL     string      `json:"languages_url"`
	// StargazersURL    string      `json:"stargazers_url"`
	// ContributorsURL  string      `json:"contributors_url"`
	// SubscribersURL   string      `json:"subscribers_url"`
	// SubscriptionURL  string      `json:"subscription_url"`
	// CommitsURL       string      `json:"commits_url"`
	// GitCommitsURL    string      `json:"git_commits_url"`
	// CommentsURL      string      `json:"comments_url"`
	// IssueCommentURL  string      `json:"issue_comment_url"`
	// ContentsURL      string      `json:"contents_url"`
	// CompareURL       string      `json:"compare_url"`
	// MergesURL        string      `json:"merges_url"`
	// ArchiveURL       string      `json:"archive_url"`
	// DownloadsURL     string      `json:"downloads_url"`
	// IssuesURL        string      `json:"issues_url"`
	// PullsURL         string      `json:"pulls_url"`
	// MilestonesURL    string      `json:"milestones_url"`
	// NotificationsURL string      `json:"notifications_url"`
	// LabelsURL        string      `json:"labels_url"`
	// ReleasesURL      string      `json:"releases_url"`
	// DeploymentsURL   string      `json:"deployments_url"`
	// CreatedAt        time.Time   `json:"created_at"`
	// UpdatedAt        time.Time   `json:"updated_at"`
	// PushedAt         time.Time   `json:"pushed_at"`
	// GitURL           string      `json:"git_url"`
	// SSHURL           string      `json:"ssh_url"`
	// CloneURL         string      `json:"clone_url"`
	// SvnURL           string      `json:"svn_url"`
	// Homepage         string      `json:"homepage"`
	// Size             int         `json:"size"`
	// StargazersCount  int         `json:"stargazers_count"`
	// WatchersCount    int         `json:"watchers_count"`
	// Language         string      `json:"language"`
	// HasIssues        bool        `json:"has_issues"`
	// HasProjects      bool        `json:"has_projects"`
	// HasDownloads     bool        `json:"has_downloads"`
	// HasWiki          bool        `json:"has_wiki"`
	// HasPages         bool        `json:"has_pages"`
	// ForksCount       int         `json:"forks_count"`
	// MirrorURL        interface{} `json:"mirror_url"`
	// Archived         bool        `json:"archived"`
	// OpenIssuesCount  int         `json:"open_issues_count"`
	// License          struct {
	// 	Key    string `json:"key"`
	// 	Name   string `json:"name"`
	// 	SpdxID string `json:"spdx_id"`
	// 	URL    string `json:"url"`
	// 	NodeID string `json:"node_id"`
	// } `json:"license"`
	// Forks         int    `json:"forks"`
	// OpenIssues    int    `json:"open_issues"`
	// Watchers      int    `json:"watchers"`
	// DefaultBranch string `json:"default_branch"`
}
