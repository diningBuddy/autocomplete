package model

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

type AutocompleteRedis struct {
	Search *redis.Client
}

type Items []Item
type MinimalItems []MinimalItem
type Autocomplete struct {
	Query   string      `json:"query"`
	Results interface{} `json:"results"`
}

func (i Items) WithMinimalInfo() MinimalItems {
	r := make(MinimalItems, len(i))
	for k, v := range i {
		r[k] = v.WithMinimalInfo()
	}
	return r
}

type Item struct {
	DataType           string    `json:"data_type"`
	Info               Info      `json:"info"`
	Score              int64     `json:"score"`
	OrgDisplay         string    `json:"org_display"`
	KeyList            []KeyList `json:"key_list"`
	MatchKey           string    `json:"match_key"`
	BeginPos           int64     `json:"begin_pos"`
	PositionList       []int64   `json:"position_list"`
	IsDuplicated       bool      `json:"is_duplicated"`
	HighlightedDisplay string    `json:"highlighted_display"`
	MatchBoosting      int64     `json:"match_boosting"`
	Category           string    `json:"category,omitempty"`
}

type MinimalItem struct {
	OrgDisplay         string `json:"org_display"`
	HighlightedDisplay string `json:"highlighted_display"`
	Category           string `json:"category,omitempty"`
}

type MinimalInfo struct {
	SearchKeyword *string `json:"search_keyword,omitempty"`

	// restaurant category
	OldCode           *string `json:"old_code,omitempty"`
	DisplayCategoryId *int64  `json:"display_category_id,omitempty"`

	// expert consulting
	ContentUrl  *string `json:"content_url,omitempty"`
	IconFileUrl *string `json:"icon_file_url,omitempty"`
	CTA         *string `json:"cta,omitempty"`

	// apartment
	OID *string `json:"oid,omitempty"`
}

func (i Item) WithMinimalInfo() MinimalItem {
	return MinimalItem{
		OrgDisplay:         i.OrgDisplay,
		HighlightedDisplay: i.HighlightedDisplay,
		Category:           i.Category,
	}
}

type Info struct {
	// expert consulting
	ContentType *int64  `json:"content_type,omitempty"`
	ContentUrl  *string `json:"content_url,omitempty"`
	IconFileUrl *string `json:"icon_file_url,omitempty"`
	CTA         *string `json:"cta,omitempty"`
	Keywords    *string `json:"keywords,omitempty"`
	Keyword     *string `json:"keyword,omitempty"`

	// restaurant category
	CategoryID      *string `json:"category_id,omitempty"`
	UserCntLong     *int64  `json:"user_cnt_long,omitempty"`
	UserCntShort    *int64  `json:"user_cnt_short,omitempty"`
	CntLong         *int64  `json:"cnt_long,omitempty"`
	CntShort        *int64  `json:"cnt_short,omitempty"`
	ID              *int64  `json:"id,omitempty"`
	DisplayName     *string `json:"display_name,omitempty"`
	FullDisplayName *string `json:"full_display_name,omitempty"`
	OldCode         *string `json:"old_code,omitempty"`
	ImageURL        *string `json:"image_url,omitempty"`
	Date            string  `json:"date"`
	SearchKeyword   *string `json:"search_keyword,omitempty"`
	QcLong          *int64  `json:"qc_long,omitempty"`
	QcShort         *int64  `json:"qc_short,omitempty"`
	CcLong          *int64  `json:"cc_long,omitempty"`
	CcShort         *int64  `json:"cc_short,omitempty"`

	// address
	Sidoname    *string `json:"sidoname,omitempty"`
	Sigunguname *string `json:"sigunguname,omitempty"`
	Dongname    *string `json:"dongname,omitempty"`
	Roadname    *string `json:"roadname,omitempty"`
	Cnt         *int64  `json:"cnt,omitempty"`

	// apartment
	OID *string `json:"oid,omitempty"`
}

type KeyList struct {
	Key            string `json:"key"`
	BeginPos       int64  `json:"begin_pos"`
	IsChosungMatch bool   `json:"is_chosung_match"`
}

func GetItems(raw []byte) (res Items, err error) {
	err = json.Unmarshal(raw, &res)
	return res, err
}

type Version struct {
	Restaurant map[string]string
}
