// Basic OAuth2 command line helper.

package auth

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/oauth2"
)

var (
	clientID     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	scope        = flag.String("scope", "flow private manage profile offline_access", "OAuth scope")
	redirectURL  = flag.String("redirect_url", "urn:ietf:wg:oauth:2.0:oob", "Redirect URL")
	authURL      = flag.String("auth_url", "https://api.flowdock.com/oauth/authorize", "Authentication URL")
	tokenURL     = flag.String("token_url", "https://api.flowdock.com/oauth/token", "Token URL")
	code         = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain Client ID and Secret, see the "OAuth 2 Credentials" section under
the "API Access" tab on this page: https://flowdock.com/account/authorized_applications

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache and you may run without the -id and -secret flags.
`

func AuthenticationRequest() *http.Client {
	flag.Parse()
	var err error
	ctx := context.TODO()

	// Try to pull the token from the cache; if this fails, we need to get one.
	tok, _ := cachedAuthToken(*cachefile) // Read auth token from cache, ignoring errors.
	if !tok.Valid() {
		if *clientID == "" || *clientSecret == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
	}

	// Set up a configuration.
	config := &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURL,
		Scopes:       strings.Split(*scope, " "),
		Endpoint: oauth2.Endpoint{
			AuthURL:  *authURL,
			TokenURL: *tokenURL,
		},
	}

	if !tok.Valid() {
		if *code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE")
			fmt.Println(url)
			os.Exit(0)
		}

		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		tok, err = config.Exchange(ctx, *code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
	}

	// Set up the token cache
	cache := &cachingTokenSource{
		Path:   *cachefile,
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
