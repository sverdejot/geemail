package auth

import (
    _ "embed"
    "log"
    "net/http"
)

//go:embed static/callback.html
var callbackPage []byte

func callback() string {
    var code string

    srv := &http.Server{
        Addr: ":80",
    }

    mux := http.NewServeMux()

    fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            go srv.Shutdown(r.Context())
        }()

        code = r.URL.Query().Get("code")
        w.WriteHeader(http.StatusOK)
        w.Write(callbackPage)
    })

    mux.Handle("/callback", fn)
    srv.Handler = mux

    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("failed listening for callback requests: %v", err)
    }

    return code
}
