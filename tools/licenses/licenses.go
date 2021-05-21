package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/andrewarchi/browser/jsonutil"
	"github.com/andrewarchi/ws-corpus/tools"
	"github.com/pelletier/go-toml"
)

// Uses SPDX license IDs
// https://spdx.org/licenses/

// Stack Exchange content licenses
// https://meta.stackexchange.com/help/licensing
// CC BY-SA 2.5 before 2011-04-08
// CC BY-SA 3.0 from 2011-04-08 up to but not including 2018-05-02
// CC BY-SA 4.0 on or after 2018-05-02

var ghRepo = regexp.MustCompile("^https://github.com/[^/]+/[^/]+$")
var ghToken = os.Getenv("GITHUB_ACCESS_TOKEN")

func main() {
	var projects []tools.Project
	try(jsonutil.DecodeFile("projects.json", &projects))
	for i := range projects {
		p := &projects[i]
		l := p.License
		var err error
		if l == "" && len(p.Source) > 0 && ghRepo.MatchString(p.Source[0]) {
			fmt.Fprintf(os.Stderr, "Getting license for %s from GitHub\n", p.Path)
			repo := strings.TrimPrefix(p.Source[0], "https://github.com/")
			l, err = getGitHubLicense(repo)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				break
			}
		}
		if (l == "" || l == "not found" || l == "none") && p.Path != "" {
			l, err = getPackageJSONLicense(p.Path)
			if l != "" {
				fmt.Fprintf(os.Stderr, "Got license for %s from package.json\n", p.Path)
			} else if err == nil {
				l, err = getCargoTOMLLicense(p.Path)
				if l != "" {
					fmt.Fprintf(os.Stderr, "Got license for %s from Cargo.toml\n", p.Path)
				}
			}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if l != "" {
			p.License = l
		}
	}
	e := json.NewEncoder(os.Stdout)
	e.SetEscapeHTML(false)
	try(e.Encode(projects))
}

func getGitHubLicense(repo string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/license", repo), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if ghToken != "" {
		req.Header.Add("Authorization", "token "+ghToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var l struct {
		Message string `json:"message"`
		License struct {
			ID string `json:"spdx_id"`
		} `json:"license"`
	}
	if err := jsonutil.DecodeAllowUnknownFields(resp.Body, &l); err != nil {
		return "", err
	}
	if l.License.ID != "" {
		if l.License.ID == "NOASSERTION" {
			return "other", nil
		}
		return l.License.ID, nil
	}
	if l.Message != "" {
		if l.Message == "Not Found" {
			return "not found", nil
		}
		return "", fmt.Errorf("message: %s", l.Message)
	}
	return "", nil
}

func getPackageJSONLicense(path string) (string, error) {
	filename := path + "/package.json"
	if stat, err := os.Stat(filename); err != nil || stat.IsDir() {
		return "", nil
	}
	var pack struct {
		License string `json:"license"`
	}
	if err := jsonutil.DecodeFileAllowUnknownFields(filename, &pack); err != nil {
		return "", err
	}
	return pack.License, nil
}

func getCargoTOMLLicense(path string) (string, error) {
	filename := path + "/Cargo.toml"
	if stat, err := os.Stat(filename); err != nil || stat.IsDir() {
		return "", nil
	}
	tree, err := toml.LoadFile(filename)
	if err != nil {
		return "", err
	}
	if l, ok := tree.Get("package.license").(string); ok {
		return l, nil
	}
	return "", nil
}

func try(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
