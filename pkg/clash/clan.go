package goclash

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Clan holds almost all information about a clan
type Clan struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
	Type string `json:"type"`
	// Warning: Description of clan is not available when searching for clans
	Decription       string   `json:"description"`
	Location         Location `json:"location"`
	BadgeUrl         BadgeUrl `json:"badgeUrls"`
	Level            int      `json:"clanLevel"`
	Points           int32    `json:"clanPoints"`
	VersusPoints     int32    `json:"clanVersusPoints"`
	RequiredTrophies int      `json:"requiredTrophies"`
	WarFrequency     string   `json:"warFrequency"`
	WarWinStreak     int      `json:"warWinStreak"`
	WarWins          int      `json:"warWins"`
	WarTies          int      `json:"warTies"`
	WarLosses        int      `json:"warLosses"`
	IsWarLogPublic   bool     `json:"isWarLogPublic"`
	MemberCount      int      `json:"members"`
	Members          []Member `json:"memberList"`
	Labels           []Label  `json:"labels"`
}

// BadgeUrl holds the url to badge images in various sizes
type BadgeUrl struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// Member holds information about a clan member
type Member struct {
	Tag               string `json:"tag"`
	Name              string `json:"name"`
	Role              string `json:"role"`
	ExpLevel          int    `json:"expLevel"`
	League            League `json:"league"`
	Trophies          int    `json:"trophies"`
	VersusTrophies    int    `json:"versusTrophies"`
	Rank              int    `json:"clanRank"`
	PreviousRank      int    `json:"previousClanRank"`
	Donations         int    `json:"donations"`
	DonationsReceived int    `json:"donationsReceived"`
}

// ClashTime decodes the clash of clans timestamp string
type ClashTime struct {
	time.Time
}

func (ct *ClashTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if s == "null" {
		ct.Time = time.Time{}
		return nil
	}
	s = strings.Trim(s, "\"")
	ct.Time, err = time.Parse("20060102T150405.000Z", s)
	if err != nil {
		return err
	}

	return nil
}

// WarLog is a log of a previous war
type WarLog struct {
	Result       string    `json:"result"`
	EndTime      ClashTime `json:"endTime"`
	TeamSize     int       `json:"teamSize"`
	Clan         WarClan   `json:"clan"`
	OpponentClan WarClan   `json:"opponent"`
}

// WarClan is a clan that is/has participated in a war
type WarClan struct {
	Tag                   string   `json:"tag"`
	Name                  string   `json:"name"`
	BadgeUrl              BadgeUrl `json:"badgeUrls"`
	Level                 int      `json:"clanLevel"`
	Attacks               int      `json:"attacks"`
	Stars                 int      `json:"stars"`
	DestructionPercentage float32  `json:"destructionPercentage"`
	// Warning: Experience isn't earned when a clan loses a war
	ExpEarned int         `json:"expEarned"`
	Team      []WarMember `json:"members"`
}

// WarMember is a clan member that is currently in a war
type WarMember struct {
	Tag                string   `json:"tag"`
	Name               string   `json:"name"`
	MapPosition        int      `json:"mapPosition"`
	TownhallLevel      int      `json:"townhallLevel"`
	OpponentAttacks    int      `json:"opponentAttacks"`
	BestOpponentAttack Attack   `json:"bestOpponentAttack"`
	Attacks            []Attack `json:"attacks"`
}

// Attack is an attack made in a clan war
type Attack struct {
	Order                 int     `json:"order"`
	AttackerTag           string  `json:"attackerTag"`
	DefenderTag           string  `json:"defenderTag"`
	Stars                 int     `json:"stars"`
	DestructionPercentage float32 `json:"destructionPercentage"`
}

// War is a war that is currently active
type War struct {
	State                string    `json:"state"`
	Clan                 WarClan   `json:"clan"`
	OpponentClan         WarClan   `json:"opponent"`
	TeamSize             int       `json:"teamSize"`
	StartTime            ClashTime `json:"startTime"`
	PreparationStartTime ClashTime `json:"preparationStartTime"`
	EndTime              ClashTime `json:"endTime"`
}

// LeagueGroup is a clans current league group
type LeagueGroup struct {
	State    string            `json:"state"`
	Tag      string            `json:"tag"`
	Season   string            `json:"season"`
	Clans    []LeagueGroupClan `json:"clans"`
	BadgeUrl BadgeUrl          `json:"badgeUrls"`
	Rounds   []Round           `json:"rounds"`
}

// Round holds a slice of war tags
type Round struct {
	Tags []string `json:"warTags"`
}

// LeagueGroupClan is a clan that is in a league war
type LeagueGroupClan struct {
	Tag   string              `json:"tag"`
	Name  string              `json:"name"`
	Level int                 `json:"clanLevel"`
	Team  []LeagueGroupMember `json:"members"`
}

// LeagueGroupMember is a member that is in a league war
type LeagueGroupMember struct {
	Tag           string `json:"tag"`
	Name          string `json:"name"`
	TownhallLevel int    `json:"townhallLevel"`
}

// ClanService holds all the methods for receiving information on clans
type ClanService service

// Search will search for clans names that match the query. The query must
// be greater than 3 characters. The query is interpreted as a wild card search.
func (c *ClanService) Search(query string, opt Optional) ([]*Clan, error) {
	if len(query) < 3 {
		return nil, fmt.Errorf("search query is less than 3 characters long")
	}

	var (
		err error
		v   = url.Values{}
		req *http.Request
	)

	if !optionalNil(opt) {
		v, err = opt.Encode()
		if err != nil {
			return nil, errInvalidOptional
		}
	}
	v.Add("name", url.QueryEscape(query))

	req, err = c.client.NewRequest("clans", v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		Clans []*Clan `json:"items"`
	}

	_, err = c.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.Clans, nil
}

// Get will retrieve a single clan by its clan tag.
func (c *ClanService) Get(tag string) (*Clan, error) {
	if err := validateTag(tag); err != nil {
		return nil, err
	}

	req, err := c.client.NewRequest(buildURLPath("clans/", url.QueryEscape(tag)), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create a new request: %s", err.Error())
	}

	var clan Clan

	_, err = c.client.Do(req, &clan)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return &clan, nil
}

// GetMembers will retrieve the members of a clan.
func (c *ClanService) GetMembers(tag string, opt Optional) ([]*Member, error) {
	if err := validateTag(tag); err != nil {
		return nil, err
	}

	var (
		err error
		v   = url.Values{}
		req *http.Request
	)

	if !optionalNil(opt) {
		v, err = opt.Encode()
		if err != nil {
			return nil, errInvalidOptional
		}
	}

	req, err = c.client.NewRequest(buildURLPath("clans/", url.QueryEscape(tag), "/members"), v)
	if err != nil {
		return nil, fmt.Errorf("could not create a new request: %s", err.Error())
	}

	var items struct {
		Members []*Member `json:"items"`
	}

	_, err = c.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.Members, nil
}

// GetWarLogs will retrieve a clans war logs if it's made public
func (c *ClanService) GetWarLogs(tag string, opt Optional) ([]*WarLog, error) {
	if err := validateTag(tag); err != nil {
		return nil, err
	}

	var (
		err error
		v   = url.Values{}
		req *http.Request
	)

	if !optionalNil(opt) {
		v, err = opt.Encode()
		if err != nil {
			return nil, errInvalidOptional
		}
	}

	req, err = c.client.NewRequest(buildURLPath("clans/", url.QueryEscape(tag), "/warlog"), v)
	if err != nil {
		return nil, fmt.Errorf("could not create a new request: %s", err.Error())
	}

	var items struct {
		WarLogs []*WarLog `json:"items"`
	}

	_, err = c.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.WarLogs, nil
}

// GetCurrentWar will retrieve a clans current war if there is one
func (c *ClanService) GetCurrentWar(tag string) (*War, error) {
	if err := validateTag(tag); err != nil {
		return nil, err
	}

	req, err := c.client.NewRequest(buildURLPath("clans/", url.QueryEscape(tag), "/currentwar"), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create a new request: %s", err.Error())
	}

	var currentWar War

	_, err = c.client.Do(req, &currentWar)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return &currentWar, nil
}

func (c *ClanService) getLeagueGroup(pathStr string) (*LeagueGroup, error) {
	req, err := c.client.NewRequest(pathStr, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create a new request: %s", err.Error())
	}

	var leagueGroup LeagueGroup

	_, err = c.client.Do(req, &leagueGroup)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return &leagueGroup, nil
}

// GetLeagueGroup will get a clans league group
func (c *ClanService) GetLeagueGroup(tag string) (*LeagueGroup, error) {
	return c.getLeagueGroup(buildURLPath("clans/", url.QueryEscape(tag), "/currentwar/leaguegroup"))
}

// GetWarLeagueWar will get a clans league war
func (c *ClanService) GetWarLeagueWar(warTag string) (*LeagueGroup, error) {
	return c.getLeagueGroup(buildURLPath("clanwarleagues/wars/", url.QueryEscape(warTag)))
}
