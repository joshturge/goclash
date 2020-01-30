package goclash

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Location holds information about a location
type Location struct {
	Id            int32  `json:"id"`
	Name          string `json:"name"`
	LocalizedName string `json:"localizedName"`
	IsCountry     bool   `json:"isCountry"`
	Code          string `json:"countryCode"`
}

// ClanRanking for a clan in a location
type ClanRanking struct {
	Tag          string `json:"tag"`
	Name         string `json:"name"`
	MemberCount  int    `json:"members"`
	Level        int    `json:"clanLevel"`
	Rank         int    `json:"rank"`
	PreviousRank int    `json:"previousRank"`
	Points       int    `json:"clanPoints"`
	BadgeUrl     BadgeUrl
	Location     Location `json:"location"`
}

// PlayerRanking for a player in a location
type PlayerRanking struct {
	Tag          string            `json:"tag"`
	Name         string            `json:"name"`
	League       League            `json:"league"`
	Clan         PlayerRankingClan `json:"clan"`
	AttackWins   int               `json:"attackWins"`
	DefenceWins  int               `json:"defenseWins"`
	ExpLevel     int               `json:"expLevel"`
	Rank         int               `json:"rank"`
	PreviousRank int               `json:"previousRank"`
	Trophies     int               `json:"trophies"`
}

// PlayerRankingClan holds information about a players clan
type PlayerRankingClan struct {
	Tag      string   `json:"tag"`
	Name     string   `json:"name"`
	BadgeUrl BadgeUrl `json:"badgeUrls"`
}

// ClanVersusRanking holds information about a clans points in a location
type ClanVersusRanking struct {
	Points       int `json:"clanPoints"`
	VersusPoints int `json:"clanVersusPoints"`
}

// PlayerVersusRanking holds information about a players versus ranking in a location
type PlayerVersusRanking struct {
	Tag          string            `json:"tag"`
	Name         string            `json:"name"`
	Clan         PlayerRankingClan `json:"clan"`
	BattleWins   int               `json:"versusBattleWins"`
	ExpLevel     int               `json:"expLevel"`
	Rank         int               `json:"rank"`
	PreviousRank int               `json:"previousRank"`
	Trophies     int               `json:"versusTrophies"`
}

// LocationService holds methods for getting information on a location
type LocationService service

// List will list all locations available
func (l *LocationService) List(opt *Control) ([]*Location, error) {
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
	req, err = l.client.NewRequest("locations", v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		Locations []*Location `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.Locations, nil
}

// Get will retrieve a location by Id
func (l *LocationService) Get(locationId int32) (*Location, error) {
	var path strings.Builder
	path.WriteString("locations/")
	path.WriteString(strconv.FormatInt(int64(locationId), 10))

	req, err := l.client.NewRequest(path.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var location Location

	_, err = l.client.Do(req, &location)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return &location, nil
}

// GetClanRankings will get clan rankings in a specific location
func (l *LocationService) GetClanRankings(locationId int32, opt *Control) ([]*ClanRanking, error) {
	var path strings.Builder
	path.WriteString("locations/")
	path.WriteString(strconv.FormatInt(int64(locationId), 10))
	path.WriteString("/rankings/clans")

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
	req, err = l.client.NewRequest(path.String(), v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		ClanRankings []*ClanRanking `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.ClanRankings, nil
}

// GetPlayerRankings will get player rankings in a specific location
func (l *LocationService) GetPlayerRankings(locationId int32, opt *Control) ([]*PlayerRanking, error) {
	var path strings.Builder
	path.WriteString("locations/")
	path.WriteString(strconv.FormatInt(int64(locationId), 10))
	path.WriteString("/rankings/players")

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
	req, err = l.client.NewRequest(path.String(), v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		PlayerRankings []*PlayerRanking `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.PlayerRankings, nil
}

// GetClanVersusRankings will get clan versus rankings in a specific location
func (l *LocationService) GetClanVersusRankings(locationId int32, opt *Control) ([]*ClanVersusRanking,
	error) {
	var path strings.Builder
	path.WriteString("locations/")
	path.WriteString(strconv.FormatInt(int64(locationId), 10))
	path.WriteString("/rankings/clans-versus")

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
	req, err = l.client.NewRequest(path.String(), v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		ClanVersusRankings []*ClanVersusRanking `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.ClanVersusRankings, nil
}

// GetPlayerVersusRankings will get player versus rankings in a specific location
func (l *LocationService) GetPlayerVersusRankings(locationId int32, opt *Control) ([]*PlayerVersusRanking,
	error) {
	var path strings.Builder
	path.WriteString("locations/")
	path.WriteString(strconv.FormatInt(int64(locationId), 10))
	path.WriteString("/rankings/players-versus")

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
	req, err = l.client.NewRequest(path.String(), v)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %s", err.Error())
	}

	var items struct {
		PlayerVersusRankings []*PlayerVersusRanking `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.PlayerVersusRankings, nil
}
