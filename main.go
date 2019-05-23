package main

import (
	"compress/flate"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atrox/env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var client *github.Client

func init() {
	githubToken := env.Get("GITHUB_TOKEN")
	if githubToken == "" {
		client = github.NewClient(nil)
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)
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
		// TODO: create playground to dynamically create the badges
		http.Redirect(w, r, "https://atrox.dev", http.StatusFound)
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

		checks, _, err := client.Checks.ListCheckSuitesForRef(ctx, owner, repo, ref, &github.ListCheckSuiteOptions{
			AppID: github.Int(15368),
		})
		if err != nil {
			sendJSONResponse(w, err)
			return
		}

		check := getRelevantCheckSuite(checks.CheckSuites)
		if check == nil {
			sendJSONResponse(w, errors.New("no check found"))
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
		sendEndpointResponse(w, endpoint)
		return
	case "completed":
		// continue
	default:
		endpoint.ServerError()
		sendEndpointResponse(w, endpoint)
		return
	}

	conclusion := check.GetConclusion()
	if conclusion == "" {
		endpoint.ServerError()
		sendEndpointResponse(w, endpoint)
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
	sendEndpointResponse(w, endpoint)
}

func gotoRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	check := ctx.Value("check").(*github.CheckSuite)

	owner := chi.URLParamFromCtx(ctx, "owner")
	repo := chi.URLParamFromCtx(ctx, "repo")

	runs, _, err := client.Checks.ListCheckRunsCheckSuite(ctx, owner, repo, check.GetID(), &github.ListCheckRunsOptions{})
	if err != nil {
		sendJSONResponse(w, err)
		return
	}

	if len(runs.CheckRuns) <= 0 {
		sendJSONResponse(w, errors.New("no check runs found"))
		return
	}

	http.Redirect(w, r, runs.CheckRuns[runs.GetTotal()-1].GetHTMLURL(), http.StatusFound)
}

// getRelevantCheckSuite returns the most relevant check suite
func getRelevantCheckSuite(checks []*github.CheckSuite) (finalCheck *github.CheckSuite) {
	// var endCheck *github.CheckSuite

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
		if conclusion == "" {
			return check
		}

		switch conclusion {
		case "success":
			finalCheck = check
		case "failure":
			return check
		case "neutral":
			if finalCheck == nil || finalCheck.GetConclusion() != "success" {
				finalCheck = check
			}
		case "cancelled":
			return check
		case "timed_out":
			return check
		case "action_required":
			return check
		default:
			return check
		}
	}

	return
}
