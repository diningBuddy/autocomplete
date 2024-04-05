package handler

import (
	"context"
	"net/http"
	"strings"
	"unicode"

	"google.golang.org/grpc/metadata"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/skku/autocomplete/app/model"
	"github.com/skku/autocomplete/hangul"
)

var (
	RestaurantKeyVersions = [...]string{"v1"}
)

const (
	RestaurantVersion = "v1" // parameterVersion "" or "v1"

	restaurantService = "restaurant"
)

func getItems(ctx context.Context, rd *redis.Client, key string) (model.Items, error) {
	res, err := rd.WithContext(ctx).Get(key).Result()
	if err != nil {
		return model.Items{}, err
	}
	items, err := model.GetItems([]byte(res))
	if err != nil {
		return model.Items{}, err
	}
	return items, nil
}

func getAutocomplete(ctx context.Context, rd *redis.Client, keyPrefix, query string) (model.Items, error) {
	key := strings.Join(strings.Fields(strings.TrimSpace(query)), " ")
	// 한글이면 검색어를 자음 모음으로 변환하여 검색
	if hangul.IsHangul(query) {
		log.Debug(keyPrefix + hangul.SplitJamoCharWithSplitDoubleJunJon(key))

		return getItems(ctx, rd, keyPrefix+hangul.SplitJamoCharWithSplitDoubleJunJon(key))
	}

	// 원 질의 소문자화 하여 한글이 포함된 경우는 자모음 분리 하여 검색
	if res, err := getItems(ctx, rd, keyPrefix+hangul.SplitJamoCharWithSplitDoubleJunJon(strings.ToLower((key)))); err == nil {
		log.Debug(keyPrefix + hangul.SplitJamoCharWithSplitDoubleJunJon(strings.ToLower((key))))
		return res, nil
	}

	// 전부다 대문자인 케이스면 전부다 소문자로 변경
	if isAllUpper(key) {
		key = strings.ToLower(key)
	}

	// 결과가 없는경우, 에러가 발생한 경우 한번더 한글로 전체 변환하여 검색
	log.Debug(keyPrefix + hangul.Eng2KorRaw(key))
	return getItems(ctx, rd, keyPrefix+hangul.Eng2KorRaw(key))
}

func isAllUpper(key string) bool {
	for _, k := range key {
		if !unicode.IsUpper(k) {
			return false
		}
	}
	return true
}

func autocomplete(
	ctx context.Context,
	rd *redis.Client,
	query, key, version string,
	versions map[string]string,
	rerankfunc func(string, model.Items) model.Items,
) (model.Items, error) {
	if query == "" {
		return model.Items{}, nil
	}
	if _, ok := versions[version]; !ok {
		return model.Items{}, nil
	}
	keyPrefix := key + ":" + version + ":" + versions[version] + ":"

	items, err := getAutocomplete(ctx, rd, keyPrefix, query)
	if err != nil {
		if err == redis.Nil {
			return model.Items{}, nil
		}
		return model.Items{}, err
	}
	if rerankfunc != nil {
		return rerankfunc(query, items), nil
	}

	return items, nil
}

func RestaurantAutocomplete(rd *model.AutocompleteRedis, v *model.Version, w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")

	ctx := convertHeaderToContext(r)
	version := RestaurantVersion
	res, err := autocomplete(ctx, rd.Search, query, restaurantService, version, v.Restaurant, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, model.Autocomplete{
		Query:   query,
		Results: res.WithMinimalInfo(),
	})
}

func convertHeaderToContext(r *http.Request) context.Context {
	headers := map[string]string{}

	for k, v := range r.Header {
		if len(v) == 0 {
			continue
		}
		headers[k] = v[0]
	}

	return metadata.NewIncomingContext(r.Context(), metadata.New(headers))
}
