package tent

import (
	"fmt"
	"github.com/hendrikcech/tent-scrobbler/track"
	"github.com/skratchdot/open-golang/open"
	"github.com/tent/tent-client-go"
)

func AuthUser(entity string, postType string) (client *tent.Client, err error) {
	var meta *tent.MetaPost

	// meta post
	meta, err = tent.Discover(entity)
	if err != nil {
		return
	}

	client = &tent.Client{Servers: meta.Servers}

	// app post
	post := tent.NewAppPost(&tent.App{
		Name: "Tent Scrobbler",
		URL:  "https://app.example.com",
		Types: tent.AppTypes{
			Write: []string{postType},
		},
		RedirectURI: "https://app.example.com/oauth",
		Scopes:      []string{"permissions"},
	})
	err = client.CreatePost(post)
	if err != nil {
		return
	}

	// credentials
	client.Credentials, _, err = post.LinkedCredentials()
	if err != nil {
		return
	}

	// redirect url
	oauthURL := meta.Servers[0].URLs.OAuthURL(post.ID, "randomState")

	fmt.Println("accept the authorisation request and paste the code back in")

	err = open.Run(oauthURL)
	if err != nil {
		return
	}

	// input code
	var code string
	_, err = fmt.Scanf("%s", &code)
	if err != nil {
		return
	}

	// request access token
	client.Credentials, err = client.RequestAccessToken(code)
	if err != nil {
		return
	}

	return client, nil
}

func ScrobbleTrack(client *tent.Client, postType string, track track.Track) (err error) {
	content := fmt.Sprintf(`{"name": "%s", "album": "%s", "artist": "%s"}`,
		track.Name, track.Album, track.Artist)

	post := &tent.Post{
		Type:        postType,
		Content:     []byte(content),
		Permissions: &tent.PostPermissions{},
	}

	err = client.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}
