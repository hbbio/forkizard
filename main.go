package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/gocolly/colly"
	"github.com/umpc/go-sortedmap"
	"github.com/umpc/go-sortedmap/desc"
	"gopkg.in/cheggaaa/pb.v1"
)

// forkURL returns the URL for network members.
func forkURL(repo string) string {
	return fmt.Sprintf("https://github.com/%s/network/members", repo)
}

// repoURL returns the URL of a repository.
func repoURL(repo string) string {
	return fmt.Sprintf("https://github.com/%s", repo)
}

// listForks lists all forks of repo.
func listForks(repo string) []string {
	c := colly.NewCollector()
	res := []string{}
	c.OnHTML(".repo a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		// Only append repository links.
		if strings.Count(href, "/") == 2 {
			res = append(res, href)
		}
	})
	c.Visit(forkURL(repo))
	log.Printf("%d forks\n", len(res))
	return res
}

var re = regexp.MustCompile(`(?P<ahead>\d+) commit[s]? ahead(, (?P<behind>\d+) commit[s]? behind)?`)

// compareRepo compares a repository.
// FIXME: Fork chains...
func compareRepo(fork string) (int, int) {
	ahead := -1
	behind := -1
	c := colly.NewCollector()
	c.OnHTML(".flex-auto.d-flex", func(e *colly.HTMLElement) {
		// Only considering forks ahead.
		if strings.Contains(e.Text, "ahead") {
			match := re.FindStringSubmatch(e.Text)
			ahead, _ = strconv.Atoi(match[1])
			behind, _ = strconv.Atoi(match[3])
		}
	})
	c.Visit(repoURL(fork))
	// if ahead > 0 {
	// 	pp.Println(fork, ahead, behind)
	// }
	return ahead, behind
}

func main() {
	if (len(os.Args) < 2) {
		os.Stderr.WriteString("Usage: forkizard owner/repo\n")
		os.Exit(1)
	}
	forks := listForks(os.Args[1])
	mahead := make(map[string]int)
	mbehind := make(map[string]int)
	sm := sortedmap.New(len(forks), desc.Int)
	bar := pb.StartNew(len(forks))
	for _, fork := range forks {
		bar.Increment()
		ahead, behind := compareRepo(fork)
		if ahead > 0 {
			mahead[fork] = ahead
			mbehind[fork] = behind
			sm.Insert(fork, ahead-behind)
		}
	}
	bar.FinishPrint("done")
	iter, err := sm.IterCh()
	if (err != nil) {
		bar.FinishPrint(err.Error())
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		for rec := range iter.Records() {
			key := rec.Key.(string)
			fmt.Fprintln(w, fmt.Sprintf("%s\t+%d -%d", key, mahead[key], mbehind[key]))
		}
		w.Flush()
    }
}
