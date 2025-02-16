package urlshort

import (
	"gopkg.in/yaml.v3"
	"net/http"
)

type YMLData struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if pathsToUrls[url] != "" {
			http.Redirect(w, r, pathsToUrls[url], http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var data []YMLData
	err := yaml.Unmarshal(yml, &data)
	if err != nil {
		panic(err)
		return nil, err
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			url := r.URL.Path
			for _, datum := range data {
				if datum.Path == url {
					http.Redirect(w, r, datum.Path, http.StatusSeeOther)
				}
			}
			fallback.ServeHTTP(w, r)
		}, nil
	}
}
