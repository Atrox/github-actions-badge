package main

import (
	"compress/flate"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atrox/env"
	raven "github.com/getsentry/raven-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var defaultGithubClient *github.Client

func init() {
	githubToken := env.Get("GITHUB_TOKEN")
	if githubToken == "" {
		defaultGithubClient = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
		tc := oauth2.NewClient(context.Background(), ts)
		defaultGithubClient = github.NewClient(tc)
	}

	// capture errors in production
	if env.IsProduction() {
		err := raven.SetDSN(env.Get("SENTRY_DSN"))
		if err != nil {
			log.Fatalf("raven could not be initialized: %s", err.Error())
		}
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.NewCompressor(flate.DefaultCompression).Handler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, playgroundHTML)
		if err != nil {
			sendJSONResponse(w, r, err)
		}
	})

	r.Route("/{owner}/{repo}", func(r chi.Router) {
		r.Use(getCheck)

		r.Get("/badge", badgeRoute)
		r.Get("/goto", gotoRoute)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", env.GetDefault("PORT", "3000")), r); err != nil {
		log.Fatal(err)
	}
}

func getCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		owner := chi.URLParamFromCtx(ctx, "owner")
		repo := chi.URLParamFromCtx(ctx, "repo")

		ref := r.URL.Query().Get("ref")
		if ref == "" {
			ref = "master"
		}

		var client *github.Client
		token := r.URL.Query().Get("token")
		if token != "" {
			ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
			tc := oauth2.NewClient(ctx, ts)
			client = github.NewClient(tc)
		} else {
			client = defaultGithubClient
		}
		ctx = context.WithValue(ctx, "client", client)

		checks, _, err := client.Checks.ListCheckSuitesForRef(ctx, owner, repo, ref, &github.ListCheckSuiteOptions{
			AppID: github.Int(15368),
		})
		if err != nil {
			if githubError, ok := err.(*github.ErrorResponse); ok {
				if githubError.Response.StatusCode == http.StatusNotFound {
					endpoint := NewEndpoint()
					endpoint.RepositoryNotFound()
					sendEndpointResponse(w, r, endpoint)
					return
				}
			}

			sendJSONResponse(w, r, err)
			return
		}

		check := getRelevantCheckSuite(checks.CheckSuites)
		if check == nil {
			endpoint := NewEndpoint()
			endpoint.NoRuns()
			sendEndpointResponse(w, r, endpoint)
			return
		}

		ctx = context.WithValue(ctx, "check", check)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func badgeRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	check := ctx.Value("check").(*github.CheckSuite)
	endpoint := NewEndpoint()

	status := check.GetStatus()
	switch status {
	case "queued", "in_progress":
		endpoint.Pending()
		sendEndpointResponse(w, r, endpoint)
		return
	case "completed":
		// continue
	default:
		endpoint.ServerError()
		sendEndpointResponse(w, r, endpoint)
		return
	}

	conclusion := check.GetConclusion()
	if conclusion == "" {
		endpoint.ServerError()
		sendEndpointResponse(w, r, endpoint)
		return
	}

	switch conclusion {
	case "success":
		endpoint.Success()
	case "failure":
		endpoint.Failure()
	case "neutral":
		endpoint.Neutral()
	case "cancelled":
		endpoint.Cancelled()
	case "timed_out":
		endpoint.TimedOut()
	case "action_required":
		endpoint.ActionRequired()
	default:
		endpoint.ServerError()
	}
	sendEndpointResponse(w, r, endpoint)
}

func gotoRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	check := ctx.Value("check").(*github.CheckSuite)
	client := ctx.Value("client").(*github.Client)

	owner := chi.URLParamFromCtx(ctx, "owner")
	repo := chi.URLParamFromCtx(ctx, "repo")

	runs, _, err := client.Checks.ListCheckRunsCheckSuite(ctx, owner, repo, check.GetID(), &github.ListCheckRunsOptions{})
	if err != nil {
		sendJSONResponse(w, r, err)
		return
	}

	if len(runs.CheckRuns) <= 0 {
		sendJSONResponse(w, r, errors.New("no check runs found"))
		return
	}

	http.Redirect(w, r, runs.CheckRuns[runs.GetTotal()-1].GetHTMLURL(), http.StatusFound)
}

// getRelevantCheckSuite returns the most relevant check suite
func getRelevantCheckSuite(checks []*github.CheckSuite) (finalCheck *github.CheckSuite) {
	for _, check := range checks {
		status := check.GetStatus()
		switch status {
		case "queued", "in_progress":
			return check
		case "completed":
			// continue
		default:
			return check
		}

		conclusion := check.GetConclusion()
		switch conclusion {
		case "success":
			finalCheck = check
		case "neutral":
			if finalCheck == nil || finalCheck.GetConclusion() != "success" {
				finalCheck = check
			}
		case "failure", "cancelled", "timed_out", "action_required":
			return check
		default:
			return check
		}
	}
	return
}
