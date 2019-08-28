package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	lines, err := parse("./sampledata.csv")
	if err != nil {
		panic(err)
	}
	q1(lines)
	fmt.Println()
	q2(lines)
	fmt.Println()
	q3(lines)
	fmt.Println()
	q4(lines)
}

/*
For each user, determine the first site he/she visited and the last site he/she visited based on the timestamp data.
Compute the number of users whose first/last visits are to the same website. What is the number?
*/
func q4(lines []line) {
	// This problem specifically mentions timestamp data. Unlike Q3, I am not going to assume the input set is sorted.

	// Iterate over user actions and save their earliest and latest actions.
	// Keep a running counter of how many users have the same first and last actions and increment/decrement as appropriate.

	actionmap := make(map[string]firstLastActions)
	var same int
	for _, l := range lines {
		// does the user exist?
		if a, uexists := actionmap[l.UserID]; !uexists {
			// user does not exist
			actionmap[l.UserID] = firstLastActions{l.TS, l.SiteID, l.TS, l.SiteID}
			same++
		} else {
			// user exists. Is it outside known action boundaries?
			if a.First.Sub(l.TS) > 0 {
				// before first known action.
				// Correct "same"
				if a.FSite == a.LSite { //was the same
					same--
				}
				if l.SiteID == a.LSite { //is now the same
					same++
				}
				a.First = l.TS
				a.FSite = l.SiteID
				actionmap[l.UserID] = a
			} else if a.Last.Sub(l.TS) < 0 {
				// after last known action
				// Correct "same"
				if a.FSite == a.LSite { //was the same
					same--
				}
				if l.SiteID == a.FSite { //is now the same
					same++
				}
				a.Last = l.TS
				a.LSite = l.SiteID
				actionmap[l.UserID] = a
			} // else action is neither first nor last and can be ignored.
		}
	}
	fmt.Printf("Q4, Number of visitors with a favorite site: %d\n", same)
}

type firstLastActions struct {
	First time.Time
	FSite string
	Last  time.Time
	LSite string
}

/*
For each site, compute the unique number of users whose last visit (found in the original data set) was to that site.
For instance, user "LC3561"'s last visit is to "N0OTG" based on timestamp data.
Based on this measure, what are top three sites?
(hint: site "3POLC" is ranked at 5th with 28 users whose last visit in the data set was to 3POLC).
Provide three pairs in the form (site_id, number of users).
*/
func q3(lines []line) {
	// My first impression was to do this by iterating through all visits,
	//   tracking each visit as though it is the user's last action and undoing any previous actions.
	// But that involves another map of maps, which is boring by now.

	// So I'm just going to map user actions from the end, ignoring anything else they find.
	// Note that this trivial case only works because the input set is sorted by time.
	// If it weren't, I would need a second map of user actions onto timestamp and the first to be map[userid]siteid for backout.

	lastactionmap := make(map[string]bool) // map of user onto nothing
	visitmap := make(map[string]int)       // map of site onto number of visits
	for i := len(lines) - 1; i >= 0; i-- {
		if _, exists := lastactionmap[lines[i].UserID]; !exists {
			// user does not exist
			lastactionmap[lines[i].UserID] = true
			if n, sfound := visitmap[lines[i].SiteID]; sfound {
				visitmap[lines[i].SiteID] = n + 1
			} else {
				visitmap[lines[i].SiteID] = 1
			}
		} // else user exists. Counting backwards in a sorted input set means this is not the user's last action, so ignore.
	}

	// Find the top 3 (5) most last-visited sites.
	// An item-by-item insertion sort would be O(s^2), while both a heap/pqueue and simple sort would be O(s lg s).
	mysortable := make([]mySortable, 0, len(visitmap))
	for s, n := range visitmap {
		mysortable = append(mysortable, mySortable{Site: s, Visits: n})
	}
	sort.Sort(ByVisits(mysortable))
	// output
	fmt.Print("Q3, Last Visits in format (site_id, number_of_users)\n\t")
	for _, site := range mysortable {
		fmt.Printf("(%s,%d) ", site.Site, site.Visits)
	}
	fmt.Println()
}

type mySortable struct {
	Site   string
	Visits int
}

// ByVisits implements sort.Interface for []mySortable based on the Visits field.
// See https://golang.org/pkg/sort/
type ByVisits []mySortable

func (a ByVisits) Len() int           { return len(a) }
func (a ByVisits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVisits) Less(i, j int) bool { return a[i].Visits > a[j].Visits }

/*
Between 2019-02-03 00:00:00 and 2019-02-04 23:59:59, there are four users who visited a certain site more than 10 times.
Find these four users & which sites they (each) visited more than 10 times.
Provide four triples in the form (user_id, site_id, number of visits).
*/
func q2(lines []line) {
	ts1, _ := time.Parse("2006-01-02 15:04:05", "2019-02-03 00:00:00")
	ts2, _ := time.Parse("2006-01-02 15:04:05", "2019-02-04 23:59:59")

	visitmap := make(map[string]map[string]int) //map of UserID onto map of SiteID onto visits

	// Since visits can be greater than 10, we simply mark when it passes 10 then retrieve the results later
	resultuser := make([]string, 0, 4)
	resultsite := make([]string, 0, 4)

	for _, l := range lines {
		if ts2.Sub(l.TS) <= 0 || l.TS.Sub(ts1) <= 0 {
			// either before ts1 or after ts2, ignore
			continue
		}

		if _, uexists := visitmap[l.UserID]; !uexists {
			//user does not exist, create
			smap := make(map[string]int)
			smap[l.SiteID] = 1
			visitmap[l.UserID] = smap
		} else if n, sfound := visitmap[l.UserID][l.SiteID]; !sfound { //dodged
			// user exists, site does not
			visitmap[l.UserID][l.SiteID] = 1
		} else {
			// user and site both exist
			visitmap[l.UserID][l.SiteID] = n + 1
			if n == 9 { //so new N = 10, only triggers once
				resultuser = append(resultuser, l.UserID)
				resultsite = append(resultsite, l.SiteID)
			}
		}
	}

	fmt.Print("Q2, More than 10 visits between timestamps in format (user_id, site_id, visits)\n\t")
	for i := range resultuser { //and resultsite
		fmt.Printf("(%s,%s,%d) ", resultuser[i], resultsite[i], visitmap[resultuser[i]][resultsite[i]])
	}
	fmt.Println()
}

/*
Consider only the rows with country_id = "BDV" (there are 844 such rows).
For each site_id, we can compute the number of unique user_id's found in these 844 rows.
Which site_id has the largest number of unique users? And what's the number?
*/
func q1(lines []line) { // ([]string, int) {
	var max int
	uumap := make(map[string]map[string]bool)
	var sites []string

	for _, l := range lines {
		if l.CountryID != "BDV" {
			continue
		}
		var tmax int
		if m, ok := uumap[l.SiteID]; !ok {
			// site does not exist, create
			uumap[l.SiteID] = make(map[string]bool)
			uumap[l.SiteID][l.UserID] = true
			tmax = 1
		} else if _, exists := m[l.UserID]; !exists {
			// site exists, user does not
			uumap[l.SiteID][l.UserID] = true
			tmax = len(uumap[l.SiteID])
		} // else both site and user exist, do not set tmax
		if tmax == max {
			sites = append(sites, l.SiteID)
		} else if tmax > max {
			sites = []string{l.SiteID}
			max = tmax
		}
	}
	// return sites, max
	fmt.Printf("BDV unique visitors\n\tSites: %v\n\tTotal: %d\n", sites, max)
}

type line struct {
	TS        time.Time
	UserID    string
	CountryID string
	SiteID    string
}

func parse(filename string) ([]line, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan() //ignore the first line

	lines := make([]line, 0)

	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), ",")
		t, err := time.Parse("2006-01-02 15:04:05", tokens[0]) // no timezone
		if err != nil {
			return nil, err
		}
		l := line{
			TS:        t,
			UserID:    tokens[1],
			CountryID: tokens[2],
			SiteID:    tokens[3]}
		lines = append(lines, l)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
