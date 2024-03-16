package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func BuildInfoHandler(w http.ResponseWriter, req *http.Request) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		http.Error(w, "no build info", http.StatusNotImplemented)
	}

	res := map[string]string{}
	for _, s := range bi.Settings {
		slog.Info("build-info", "key", s.Key, "value", s.Value)
		if s.Key == "vcs.revision" || s.Key == "vcs.time" || s.Key == "vcs.modified" {
			res[s.Key] = s.Value
		}
	}

	enc := json.NewEncoder(w)
	enc.Encode(res)
}
