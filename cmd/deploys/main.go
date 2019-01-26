package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/jtdoepke/go-flowdock/flowdock"
	"golang.org/x/oauth2"
	cli "gopkg.in/urfave/cli.v1"
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain Client ID and Secret, see the "OAuth 2 Credentials" section under
the "API Access" tab on this page: https://flowdock.com/account/authorized_applications

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache and you may run without the -id and -secret flags.
`

type Query struct {
	Tags []string
	Org  string
	Flow string
}

type AppDeployCount struct {
	Project     string
	Total       int
	DeployCount *map[string]int
}

const limit = 100

// example of counting number of deploys per month for specidied applications.
//
// This is utilizing the fact that a deploy command messages the flow inbox and is tagged appropriately.
func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "id", Value: "", Usage: "Client ID"},
		cli.StringFlag{Name: "secret", Value: "", Usage: "Client Secret"},
		cli.StringFlag{Name: "redirect_url", Value: "urn:ietf:wg:oauth:2.0:oob", Usage: "Redirect URL"},
		cli.StringFlag{Name: "auth_url", Value: "https://api.flowdock.com/oauth/authorize", Usage: "Authentication URL"},
		cli.StringFlag{Name: "token_url", Value: "https://api.flowdock.com/oauth/token", Usage: "Token URL"},
		cli.StringFlag{Name: "code", Value: "", Usage: "Authorization Code"},
		cli.StringFlag{Name: "cache", Value: "cache.json", Usage: "Token cache file"},
		cli.StringFlag{Name: "environment, e", Value: "production", Usage: "the deploy target"},
		cli.StringFlag{Name: "organization, o", Value: "iora", Usage: "the organization of the flow"},
		cli.StringFlag{Name: "flow, f", Value: "tech-stuff", Usage: "the name of the flow to query"},
	}

	app.Name = "deploys"
	app.Usage = "Counts the deploys for the listed applications by month"

	app.Action = func(c *cli.Context) {
		client := flowdock.NewClient(AuthenticationRequest(c))
		// args := []string{"bouncah", "icis", "cronos", "snowflake"} //c.Args()
		args := c.Args()

		if len(args) == 0 {
			_ = cli.ShowAppHelp(c)
			os.Exit(1)
		}

		channel := make(chan AppDeployCount, len([]string(args)))

		for _, app := range []string(args) {
			tags := []string{
				"deployment",
				"deploy_end",
				c.String("environment"),
				app,
			}
			q := Query{tags, c.String("organization"), c.String("flow")}
			getAppDeployCount(q, client, channel)
		}

		displayAppDeployCount(channel)
	}

	_ = app.Run(os.Args)
}

func getAppDeployCount(q Query, client *flowdock.Client, channel chan AppDeployCount) {
	var deployCount = map[string]int{}

	go func() {

		app := q.Tags[len(q.Tags)-1]
		opt := flowdock.MessagesListOptions{Limit: 100, TagMode: "and"}

		opt.Tags = q.Tags
		opt.Search = "production to production"
		opt.Event = "mail"

		messages, _, err := client.Messages.List(q.Org, q.Flow, &opt)

		if err != nil {
			log.Fatal("Get:", err)
		}

		total := 0
		for _, msg := range messages {
			if !stringInSlice("preproduction", *msg.Tags) {
				total++
				month := msg.Sent.Format("2006-Jan")
				deployCount[month]++
				// fmt.Println("MSG:", month, *msg.ID, *msg.Event, *msg.Tags)
			}
		}

		if len(messages) == limit {
			removeEarliestMonth(&deployCount)
		}

		channel <- AppDeployCount{app, total, &deployCount}
	}()
}

func removeEarliestMonth(displayCount *map[string]int) {
	sortedKeys := sortedKeys(displayCount)
	firstKey := (*sortedKeys)[0]
	delete(*displayCount, firstKey)
}

func displayAppDeployCount(adcChan <-chan AppDeployCount) {
	for i := 0; i < cap(adcChan); i++ {
		adc := <-adcChan

		fmt.Println()
		fmt.Println("Application:", adc.Project)
		fmt.Println()

		for k, v := range *adc.DeployCount {
			fmt.Println(k, v)
		}

		fmt.Println()
		fmt.Println("  Total:", adc.Total)
		fmt.Println()
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func sortedKeys(m *map[string]int) *[]string {
	mk := make([]string, len(*m))
	i := 0
	for k := range *m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)
	return &mk
}

func AuthenticationRequest(c *cli.Context) *http.Client {
	var err error
	ctx := context.TODO()

	// Try to pull the token from the cache; if this fails, we need to get one.
	tok, _ := cachedAuthToken(c.String("cache")) // Read auth token from cache, ignoring errors.
	if !tok.Valid() {
		if c.String("id") == "" || c.String("secret") == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
	}

	// Set up a configuration.
	config := &oauth2.Config{
		ClientID:     c.String("id"),
		ClientSecret: c.String("secret"),
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.String("auth_url"),
			TokenURL: c.String("token_url"),
		},
	}

	if !tok.Valid() {
		if c.String("code") == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE")
			fmt.Println(url)
			os.Exit(0)
		}

		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		tok, err = config.Exchange(ctx, c.String("code"))
		if err != nil {
			log.Fatal("Exchange:", err)
		}
	}

	// Set up the token cache
	cache := &cachingTokenSource{
		Path:   c.String("cache"),
		Tok:    tok,
		Source: config.TokenSource(ctx, tok),
	}
	err = cache.Write() // Cache the new token.
	if err != nil {
		log.Fatal("Cache token:", err)
	}
	fmt.Printf("Token is cached in %v\n", cache.Path)

	client := oauth2.NewClient(ctx, cache)
	return client
}

// cachedAuthToken reads a JSON-serialized oauth2.Token from a file.
func cachedAuthToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// cachingTokenSource wraps an `oauth2.TokenSource` to serialize new
// tokens to a JSON file.
type cachingTokenSource struct {
	Path   string             // Path to cache file
	Tok    *oauth2.Token      // Current token
	Source oauth2.TokenSource // Wrapped TokenSource

	lock sync.Mutex // protects Tok
}

// Token returns an `oauth2.Token`, serializing it to the cache file
// if it's new.
func (c *cachingTokenSource) Token() (*oauth2.Token, error) {
	tok, err := c.Source.Token()
	if err != nil {
		return tok, err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if tok != c.Tok {
		c.Tok = tok
		err = c.Write()
	}
	return tok, err
}

// Write writes the current token to the cache file.
func (c *cachingTokenSource) Write() error {
	if c.Tok.Valid() {
		f, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		err = json.NewEncoder(f).Encode(c.Tok)
		return err
	}
	return nil
}
