package goclash_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/joshturge/goclash/pkg/clash"
)

var (
	client    *goclash.Client
	clanTag   string
	playerTag string
	token     = os.Getenv("CLASH_TOKEN")
)

func checkClient(t *testing.T) {
	if client == nil {
		var err error
		client, err = goclash.NewClient(token)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

func TestClient(t *testing.T) {
	if token == "" {
		t.Error("CLASH_TOKEN environment variable not set")
		t.FailNow()
	}

	var err error
	client, err = goclash.NewClient(token)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestClanSearch(t *testing.T) {
	/*opt := &goclash.ClanSearchOptions{
		Control: goclash.Control{
			Limit: 5,
		},
		WarFrequency: "always",
		MinMembers:   15,
	}*/
	checkClient(t)

	clans, err := client.Clan.Search("ClashOfClans", nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(clans) == 0 {
		t.Error("No clans found")
		t.FailNow()
	}

	fmt.Println("Successfully found clans...")

	for _, clan := range clans {
		fmt.Printf("Clan Tag: %s\nClan Name: %s\n\n", clan.Tag, clan.Name)
	}

	clanTag = clans[2].Tag
}

func TestClanGet(t *testing.T) {
	checkClient(t)

	clan, err := client.Clan.Get(clanTag)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if clan.Tag != clanTag {
		t.Errorf("Wanted: %s\tGot: %s\n", clanTag, clan.Tag)
		t.FailNow()
	}

	fmt.Println("Found clan...")
	fmt.Printf("Clan Name: %s\nMember Count: %d\n\n", clan.Name, clan.MemberCount)
}

func TestClanGetMembers(t *testing.T) {
	checkClient(t)

	members, err := client.Clan.GetMembers(clanTag, &goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(members) == 0 {
		t.Error("no members found")
		t.FailNow()
	}

	for _, member := range members {
		fmt.Printf("Member Tag: %s\nMember Name: %s\n\n", member.Tag, member.Name)
	}
}

func TestClanGetWarLogs(t *testing.T) {
	checkClient(t)

	warlogs, err := client.Clan.GetWarLogs(clanTag, &goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(warlogs) == 0 {
		t.Error("no war logs found")
		t.FailNow()
	}

	fmt.Println("found war logs")

	for _, warlog := range warlogs {
		fmt.Printf("War result: %s\nWar end: %s\n\n", warlog.Result, warlog.EndTime.String())
	}
}

func TestClanGetCurrentWar(t *testing.T) {
	checkClient(t)

	war, err := client.Clan.GetCurrentWar(clanTag)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if war.Clan.Tag == "" {
		t.Error("clan tag is empty")
		t.FailNow()
	}

	fmt.Println("found current war")

	fmt.Printf("War State: %s\nWar Start Time: %s\n\n", war.State, war.StartTime.String())
}

func TestClanGetLeagueGroup(t *testing.T) {
	checkClient(t)

	leagueGroup, err := client.Clan.GetLeagueGroup(clanTag)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if leagueGroup.State == "" {
		t.Error("leagueGroup state not found")
		t.FailNow()
	}

	fmt.Println("Found league group")

	fmt.Printf("LeagueGroup State: %s\nWarTags:\n", leagueGroup.State)
	for _, warTag := range leagueGroup.Rounds[0].Tags {
		fmt.Println(warTag)
	}
}

func TestPlayerGet(t *testing.T) {
	checkClient(t)

	player, err := client.Player.Get("#RQ8JLVQ")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if player.Tag != "#RQ8JLVQ" {
		t.Error("player tag does not match")
		t.FailNow()
	}

	fmt.Println("player found...")

	fmt.Printf("Player Name: %s\nPlayerThLvl: %d\n\n", player.Name, player.TownhallLevel)
}

func TestLeagueList(t *testing.T) {
	checkClient(t)

	leagues, err := client.League.List(&goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(leagues) == 0 {
		t.Error("no leagues found")
		t.FailNow()
	}

	fmt.Println("leagues found...")

	for _, league := range leagues {
		fmt.Printf("League ID: %d\nLeague Name: %s\n\n", league.Id, league.Name)
	}
}

func TestLeagueGet(t *testing.T) {
	checkClient(t)

	league, err := client.League.Get(29000014)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if league.Id != 29000014 {
		t.Error("league id does not match")
		t.FailNow()
	}

	fmt.Println("found league")

	fmt.Printf("League Name: %s\nURL: %s\n", league.Name, league.IconUrl.Tiny)
}

func TestLeagueGetSeasons(t *testing.T) {
	checkClient(t)

	seasons, err := client.League.GetSeasons(29000022, &goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(seasons) == 0 {
		t.Error("no seasons found")
		t.FailNow()
	}

	fmt.Println("found seasons...")

	for _, season := range seasons {
		fmt.Printf("Season ID: %s\n", season.Id)
	}
}

func TestLeagueGetSeasonRankings(t *testing.T) {
	checkClient(t)

	seasonPlayers, err := client.League.GetSeasonRankings(29000022, "2015-07", &goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(seasonPlayers) == 0 {
		t.Error("season players not found")
		t.FailNow()
	}

	fmt.Println("season players found...")

	for _, seasonPlayer := range seasonPlayers {
		fmt.Printf("Season Player Name: %s\nClan Name: %s\n", seasonPlayer.Name, seasonPlayer.Clan.Name)
	}
}

func TestLocationList(t *testing.T) {
	checkClient(t)

	locations, err := client.Location.List(&goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(locations) == 0 {
		t.Error("locations not found")
		t.FailNow()
	}

	fmt.Println("locations found...")

	for _, location := range locations {
		fmt.Printf("Location ID: %d\nLocation Name: %s\n\n", location.Id, location.Name)
	}
}

func TestLocationGet(t *testing.T) {
	checkClient(t)

	location, err := client.Location.Get(32000004)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if location.Id != 32000004 {
		t.Error("Id does not match")
		t.FailNow()
	}

	fmt.Println("found location...")

	fmt.Printf("Location ID: %d\nLocation Name: %s\n\n", location.Id, location.Name)
}

func TestLocationGetClanRankings(t *testing.T) {
	checkClient(t)

	clanRankings, err := client.Location.GetClanRankings(32000003, &goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(clanRankings) == 0 {
		t.Error("clan rankings not found")
		t.FailNow()
	}

	fmt.Println("clan rankings found...")

	for _, clanRanking := range clanRankings {
		fmt.Printf("Clan Tag: %s\nClan Name: %s\n\n", clanRanking.Tag, clanRanking.Name)
	}
}

func TestLabelList(t *testing.T) {
	checkClient(t)

	labels, err := client.Label.ClanList(&goclash.Control{Limit: 5})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(labels) == 0 {
		t.Error("no labels found...")
	}

	fmt.Println("labels found...")

	for _, label := range labels {
		fmt.Printf("Label ID: %d\nLabel Name: %s\n\n", label.Id, label.Name)
	}
}
