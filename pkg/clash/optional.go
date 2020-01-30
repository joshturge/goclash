package goclash

import (
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// Optional allows users of this wrapper to create the own optional parameters
type Optional interface {
	Encode() (url.Values, error)
}

func optionalNil(opt Optional) bool {
	if opt == nil {
		return true
	}
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return true
	}
	return false
}

func encodeOptional(opt interface{}) (url.Values, error) {
	qs, err := query.Values(opt)
	if err != nil {
		return nil, err
	}

	return qs, nil
}

// ClanSearchOptions are optional parameters for a clan search
type ClanSearchOptions struct {
	WarFrequency  string  `url:"warFrequency,omitempty"`
	LocationId    int32   `url:"locationId,omitempty"`
	MinMembers    int     `url:"minMembers,omitempty"`
	MaxMembers    int     `url:"maxMembers,omitempty"`
	MinClanPoints int     `url:"minClanPoints,omitempty"`
	MinClanLevel  int     `url:"minClanLevel,omitempty"`
	Labels        []int32 `url:"labels,comma,omitempty"`
	Control
}

// Encode will encode clan search options into url values so they can be passed
// to the Clash of Clans API
func (cso *ClanSearchOptions) Encode() (url.Values, error) {
	if cso.Before != "" && cso.After != "" {
		return nil, errBeforeAfterSet
	}
	return encodeOptional(cso)
}

// Control are optional parameters to control how much data you get back from the
// Clash of Clans API
type Control struct {
	Limit  int    `url:"limit,omitempty"`
	Before string `url:"before,omitempty"`
	After  string `url:"afer,omitempty"`
}

// Encode will encode clan search options into url values so they can be passed
// to the Clash of Clans API
func (ctrl *Control) Encode() (url.Values, error) {
	if ctrl.Before != "" && ctrl.After != "" {
		return nil, errBeforeAfterSet
	}
	return encodeOptional(ctrl)
}
