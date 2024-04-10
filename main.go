package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var (
	repo    string
	filter  string
	exclude string
	all     bool
)

type tagResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func init() {
	flag.StringVar(&repo, "repo", "", "Docker Hub Repository such as alpine, busybox, and jupyter/base-notebook")
	flag.StringVar(&filter, "filter", "", "Filter to use to look for the latest tag including the specified pattern")
	flag.StringVar(&exclude, "exclude", "", "Exclude pattern")
	flag.BoolVar(&all, "all", false, "Write the list of tags found in the ./output.txt file")
	flag.Parse()
}

func getUrl(repo string) string {
	var fullUrl string

	hasSlash := strings.Contains(repo, "/")
	if hasSlash {
		fullUrl = "https://registry.hub.docker.com/api/content/v1/repositories/public/" + repo + "/tags?page_size=30"
	} else {
		fullUrl = "https://registry.hub.docker.com/api/content/v1/repositories/public/library/" + repo + "/tags?page_size=30"
	}

	return fullUrl
}

func fetchTags(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseTags(data []byte) ([]string, error) {
	var response tagResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	tags := make([]string, len(response.Results))
	for i, result := range response.Results {
		tags[i] = result.Name
	}
	return tags, nil
}

func versionLess(a, b string) bool {
	re := regexp.MustCompile(`v?(\d+)\.(\d+)\.(\d+)`)

	matchA := re.FindStringSubmatch(a)
	matchB := re.FindStringSubmatch(b)

	if len(matchA) < 4 || len(matchB) < 4 {
		// If the version format doesn't match, sort alphabetically
		return a < b
	}

	// Convert version parts to integers for comparison
	av1, _ := strconv.Atoi(matchA[1])
	av2, _ := strconv.Atoi(matchA[2])
	av3, _ := strconv.Atoi(matchA[3])

	bv1, _ := strconv.Atoi(matchB[1])
	bv2, _ := strconv.Atoi(matchB[2])
	bv3, _ := strconv.Atoi(matchB[3])

	// Compare major version first
	if av1 != bv1 {
		return av1 < bv1
	}
	// If major version is the same, compare minor version
	if av2 != bv2 {
		return av2 < bv2
	}
	// If minor version is the same, compare patch version
	if av3 != bv3 {
		return av3 < bv3
	}
	// If all version parts are the same, compare the whole string
	return a < b
}

func main() {
	var fullUrl string

	if repo == "" {
		fmt.Println("Specify a docker hub repository using '-repo'. See --help.")
		os.Exit(0)
		// fmt.Printf("Specify target repository using -repo")
	} else {
		fullUrl = getUrl(repo)
		// fmt.Printf("The target repository is %s\n", repo)
		// fmt.Printf("URL to work with is %s\n", fullUrl)
	}

	body, err := fetchTags(fullUrl)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
	}

	tags, err := parseTags(body)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
	}

	if all {
		file, err := os.Create("tags.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// fmt.Println("all flag is set")
		for _, tag := range tags {
			_, err := io.WriteString(file, tag+"\n")
			if err != nil {
				panic(err)
			}
		}
	}
	// fmt.Println(tags)

	versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+`)
	var filteredTags []string
	for _, tag := range tags {
		if versionRegex.MatchString(tag) {
			filteredTags = append(filteredTags, tag)
		}
	}
	// fmt.Println("Filtered tags containing x.y.z versioning")

	// if len(filteredTags) > 0 {
	// for _, tag := range filteredTags {
	// fmt.Println("-", tag)
	// }
	// }

	// when additional filter is specified
	if filter != "" {
		filteredTags = slices.DeleteFunc(filteredTags, func(s string) bool {
			r := regexp.MustCompile(filter)
			return !r.MatchString(s)
		})
	}

	// when exclude pattern is specified
	if exclude != "" {
		filteredTags = slices.DeleteFunc(filteredTags, func(s string) bool {
			r := regexp.MustCompile(exclude)
			return r.MatchString(s)
		})
	}

	sort.Slice(filteredTags, func(i, j int) bool {
		return versionLess(filteredTags[i], filteredTags[j])
	})

	latestTag := filteredTags[len(filteredTags)-1]
	fmt.Println(latestTag)
}
