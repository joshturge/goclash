package goclash

import (
	"fmt"
	"net/http"
	"strings"
)

// Label holds information about a Clash of Clans label
type Label struct {
	Id      int32   `json:"id"`
	Name    string  `json:"name"`
	IconUrl IconUrl `json:"iconUrls"`
}

// IconUrl holds the image url to a label in various sizes
type IconUrl struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
}

// LabelService has methods that can retrieve information on labels
type LabelService service

func (l *LabelService) labelList(lType string, opt *Control) ([]*Label, error) {
	var path strings.Builder
	path.WriteString("labels/")
	path.WriteString(lType)

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
		Labels []*Label `json:"items"`
	}

	_, err = l.client.Do(req, &items)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	return items.Labels, nil
}

// ClanList will list all the labels for clans
func (l *LabelService) ClanList(opt *Control) ([]*Label, error) {
	return l.labelList("clans", opt)
}

// PlayerList will list all the labels for players
func (l *LabelService) PlayerList(opt *Control) ([]*Label, error) {
	return l.labelList("players", opt)
}
