package goclash

import (
	"fmt"
	"net/url"
	"strings"
)

// Player holds information about a player account
type Player struct {
	Tag                  string `json:"tag"`
	Name                 string `json:"name"`
	ExpLevel             int    `json:"expLevel"`
	WarStars             int    `json:"warStars"`
	Trophies             int    `json:"trophies"`
	VersusTrophies       int    `json:"versusTrophies"`
	BestTrophies         int    `json:"bestTrophies"`
	Donations            int    `json:"donations"`
	DonationsReceived    int    `json:"donationsReceived"`
	League               League
	Clan                 PlayerClan       `json:"clan"`
	Role                 string           `json:"role"`
	AttackWins           int              `json:"attackWins"`
	DefenceWins          int              `json:"defenseWins"`
	TownhallLevel        int              `json:"townhallLevel"`
	BuilderHallLevel     int              `json:"builderHallLevel"`
	TownhallWeaponLevel  int              `json:"townhallWeaponLevel"`
	VersusBattleWins     int              `json:"versusBattleWins"`
	VersusBattleWinCount int              `json:"versusBattleWinCount"`
	LegendStatistics     LegendStatistics `json:"legendStatistics"`
	Troops               []Troop          `json:"troops"`
	Heros                []Troop          `json:"heros"`
	Spells               []Troop          `json:"spells"`
	Labels               []Label          `json:"labels"`
	Achievements         []Achievement    `json:"achievements"`
}

// PlayerClan holds information about a players clan
type PlayerClan struct {
	Tag      string   `json:"tag"`
	Name     string   `json:"name"`
	Level    int      `json:"clanLevel"`
	BadgeUrl BadgeUrl `json:"badgeUrls"`
}

// LegendStatistics holds information about a players legend statistics
type LegendStatistics struct {
	CurrentSeason        Season `json:"currentSeason"`
	PreviousSeason       Season `json:"previousSeason"`
	BestSeason           Season `json:"bestSeason"`
	PreviousVersusSeason Season `json:"previousVersusSeason"`
	BestVersusSeason     Season `json:"bestVersusSeason"`
	LegendTrophies       int    `json:"legendTrophies"`
}

// Season holds the players statistics for the current season
type Season struct {
	Id       string `json:"id"`
	Trophies int    `json:"trophies"`
	Rank     int    `json:"rank"`
}

// Troop holds information about a players troop
type Troop struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	MaxLevel int    `json:"maxLevel"`
	Village  string `json:"village"`
}

// Achievement is a players completed achievement
type Achievement struct {
	Name           string `json:"name"`
	Value          int    `json:"value"`
	Stars          int    `json:"stars"`
	Target         int    `json:"target"`
	Info           string `json:"info"`
	CompletionInfo string `json:"completionInfo"`
	Village        string `json:"village"`
}

// PlayerService holds methods that retrieves information about a player
type PlayerService service

// Get will get a player via a player tag
func (c *PlayerService) Get(tag string) (*Player, error) {
	if tag[:1] != "#" {
		return nil, fmt.Errorf("player tag was not valid: tag length: %d, first character: %s",
			len(tag), tag[:1])
	}

	var path strings.Builder
	path.WriteString("players/")
	path.WriteString(url.QueryEscape(tag))

	req, err := c.client.NewRequest(path.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create a new request: %s", err.Error())
	}

	var player Player

	_, err = c.client.Do(req, &player)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return &player, nil
}
