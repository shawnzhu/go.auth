package auth

import (
	"net/http"
)

type TwitterUser struct {
	UserId string `json:"screen_name"`
}

func (u *TwitterUser) Id()       string { return u.UserId }
func (u *TwitterUser) Provider() string { return "twitter.com" }
func (u *TwitterUser) Name()     string { return u.UserId }
func (u *TwitterUser) Email()    string { return "" }
func (u *TwitterUser) Link()     string { return "https://www.twitter.com/" + u.UserId }
func (u *TwitterUser) Picture()  string { return "" }
func (u *TwitterUser) Org()      string { return "" }


// TwitterProvider is an implementation of Twitters's Oauth1.0a protocol.
// See https://dev.twitter.com/docs/auth/implementing-sign-twitter
type TwitterProvider struct {
	OAuth1Mixin
	UserResourceUrl string
}

// NewTwitterProvider allocates and returns a new BitbucketProvider.
func NewTwitterProvider(key, secret, callback string) *TwitterProvider {
	twitter := TwitterProvider{}
	twitter.AuthorizeUrl = "https://api.twitter.com/oauth/authorize"
	twitter.RequestToken = "https://api.twitter.com/oauth/request_token"
	twitter.AccessToken =  "https://api.twitter.com/oauth/access_token"
	twitter.UserResourceUrl = "https://api.twitter.com/1.1/account/settings.json" //"http://api.twitter.com/1.1/users/show.json"

	twitter.CallbackUrl = callback
	twitter.ConsumerKey = key
	twitter.ConsumerSecret = secret
	return &twitter
}

// Redirect will do an http.Redirect, sending the user to the Twitter login
// screen.
func (self *TwitterProvider) Redirect(w http.ResponseWriter, r *http.Request) {
	err := self.OAuth1Mixin.AuthorizeRedirect(w, r, self.AuthorizeUrl)
	if err != nil {
		println(err.Error())
	}
}

// GetAuthenticatedUser will upgrade the oauth_token to an access token, and
// invoke the appropriate Twitter REST API call to get the User's information.
func (self *TwitterProvider) GetAuthenticatedUser(r *http.Request) (User, error) {

	// upgrade the oauth_token to an access token
	token, secret, err := self.OAuth1Mixin.AuthorizeToken(r)
	if err != nil {
		return nil, err
	}

	// get the Bitbucket User details
	user := TwitterUser{}
	if err := self.OAuth1Mixin.GetAuthenticatedUser(self.UserResourceUrl, token, secret, &user); err != nil {
		return nil, err
	}
	return &user, err
}
