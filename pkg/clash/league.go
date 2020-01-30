package goclash

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// League holds information about a player league
type League struct {
	Id      int32         `json:"id"`
	Name    string        `json:"name"`
	IconUrl LeagueIconUrl `json:"iconUrls"`
}

// LeagueIconUrl holds the image url for leagues in various sizes
type LeagueIconUrl struct {
	Tiny   string `json:"tiny"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
}

// LegendSeason holds the id of a legend season
type LegendSeason struct {
	Id string `json:"id"`
}

// LegendSeasonPlayer holds information on a legend league player
type LegendSeasonPlayer struct {
	Tag          string `json:"tag"`
	Name         string `json:"name"`
	ExpLevel     int    `json:"expLevel"`
	Rank         int    `json:"rank"`
	PreviousRank int    `json:"previousRank"`
	Trophies     int    `json:"trophies"`
	League       League `json:"league"`
	Clan         Clan   `json:"clan"`
	AttackWins   int    `json:"attackWins"`
	DefenceWins  int    `json:"defenseWins"`
}

// LeagueClan holds the current league of a clan
type LeagueClan struct {
	Tag      string   `json:"tag"`
	Name     string   `json:"name"`
	BadgeUrl BadgeUrl `json:"badgeUrls"`
}

// LeagueService holds methods that can get information about leagues
type LeagueService service

// List will list all leagues
func (l *LeagueService) List(opt *Control) ([]*League, error) {
	if opt != nil {
		if opt.Before != "" && opt.After != "" {
			return nil, fmt.Errorf("both Before and After have been set, this is not allowed")
		}
	}

	v, err := encodeOptional(opt)
	if err != nil {
		return nil, fmt.Errorf("could not encode optional arguments for request: %s", err.Error())
	}

	var req *http.Request
	req, err = l.client.NewRequest("leagues", v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		Leagues []*League `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.Leagues, nil
}

func (l *LeagueService) Get(leagueId int32) (*League, error) {
	var path strings.Builder
	path.WriteString("leagues/")
	path.WriteString(strconv.FormatInt(int64(leagueId), 10))

	req, err := l.client.NewRequest(path.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var league League

	_, err = l.client.Do(req, &league)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return &league, nil
}

// GetSeasons will list legend seasons that are in a specified league
func (l *LeagueService) GetSeasons(leagueId int32, opt *Control) ([]*LegendSeason, error) {
	if opt != nil {
		if opt.Before != "" && opt.After != "" {
			return nil, fmt.Errorf("both Before and After have been set, this is not allowed")
		}
	}

	v, err := encodeOptional(opt)
	if err != nil {
		return nil, fmt.Errorf("could not encode optional arguments for request: %s", err.Error())
	}

	var path strings.Builder
	path.WriteString("leagues/")
	path.WriteString(strconv.FormatInt(int64(leagueId), 10))
	path.WriteString("/seasons")

	req, err := l.client.NewRequest(path.String(), v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		LegendSeasons []*LegendSeason `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.LegendSeasons, nil
}

// GetSeasonRankings will get the rankings for a legend season
func (l *LeagueService) GetSeasonRankings(leagueId int32, seasonId string,
	opt *Control) ([]*LegendSeasonPlayer, error) {
	if opt != nil {
		if opt.Before != "" && opt.After != "" {
			return nil, fmt.Errorf("both Before and After have been set, this is not allowed")
		}
	}

	v, err := encodeOptional(opt)
	if err != nil {
		return nil, fmt.Errorf("could not encode optional arguments for request: %s", err.Error())
	}

	var path strings.Builder
	path.WriteString("leagues/")
	path.WriteString(strconv.FormatInt(int64(leagueId), 10))
	path.WriteString("/seasons/")
	path.WriteString(seasonId)

	req, err := l.client.NewRequest(path.String(), v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		SeasonPlayers []*LegendSeasonPlayer `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.SeasonPlayers, nil
}
