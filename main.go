package main

import (
	"flag"
	"fmt"
	"github.com/hendrikcech/tent-scrobbler/config"
	"github.com/hendrikcech/tent-scrobbler/spotify"
	store "github.com/hendrikcech/tent-scrobbler/tent"
	"github.com/hendrikcech/tent-scrobbler/track"
	"github.com/tent/tent-client-go"
	"github.com/tent/hawk-go"
	"os"
	"os/user"
	"time"
)

const PostType string = "http://cech.im/types/song/v0#"

var entity *string
var player *string

var configFilePath string

func init() {
	// parse arguments
	entity = flag.String("entity", "", "set entity")
	player = flag.String("player", "", "specify player")
	flag.Parse()

	// get path to config file
	usr, err := user.Current()
	maybePanic(err)
	configFilePath = usr.HomeDir + "/.tentscrobbler"
}

func main() {
	var client *tent.Client
	var c config.Config
	var err error

	exists := config.Exists(configFilePath)

	if exists {
		c, err = config.Read(configFilePath)
		maybePanic(err)
	}

	if *entity != "" {
		client, err = store.AuthUser(*entity, PostType)
		maybePanic(err)
		c.ID = client.Credentials.ID
		c.Key = client.Credentials.Key
		c.App = client.Credentials.App
		c.Servers = client.Servers
	} else {
		client = &tent.Client{
			Servers: c.Servers,
			Credentials: &hawk.Credentials{
				ID: c.ID,
				Key: c.Key,
				App: c.App,
			},
		}
	}

	if *player != "" {
		c.Player = *player
	}

	if c.ID == "" || c.Key == "" || c.App == "" || len(c.Servers) == 0 {
		fmt.Println("invalid entity config. run again with -entity entity")
		os.Exit(1)
	}
	if c.Player == "" {
		fmt.Println("no player specified. running with default (spotify).")
		c.Player = "spotify"
	}

	err = config.Write(c, configFilePath)
	maybePanic(err)

	// TODO: select player by setting up map[string]interface{}(?) with string maped to static imported package

	// setup queue
	tracks := make(chan track.Track)
	scrobbles := make(chan track.Track)

	go watchPlayer(tracks)

	for {
		select {
		case track := <-tracks:
			go watchTrack(scrobbles, track)

		case track := <-scrobbles:
			go store.ScrobbleTrack(client, PostType, track)
		}
	}

	log("Done.")
}

func watchPlayer(tracks chan track.Track) {
	var track track.Track
	var currentTrack string
	var err error

	ticker := time.NewTicker(time.Millisecond * 1000)

	for _ = range ticker.C {
		track, err = spotify.CurrentTrack()
		maybePanic(err)

		if track.Name == "" {
			log("no track playing")
			continue
		}

		log(track.Name)

		if track.ID != currentTrack {
			currentTrack = track.ID
			tracks <- track
		}
	}

	return
}

var lastMessage string

func log(message string) {
	if lastMessage == message {
		fmt.Print(".")
	} else {
		fmt.Print("\n", message)
		lastMessage = message
	}
	return
}

func watchTrack(scrobbles chan track.Track, track track.Track) {
	offset := track.Remaining / 3
	sleepFor := track.Remaining - offset

	log(fmt.Sprintf("wait for %s in %v (offset: %v)", track.Name, sleepFor, offset))

	time.Sleep(sleepFor)

	log(fmt.Sprintf("come back for%s", track.Name))

	currentTrack, err := spotify.CurrentTrack()
	maybePanic(err)

	if currentTrack.ID == track.ID {
		log("should scrobble")
		scrobbles <- track
	} else {
		log("should not scrobble")
	}
}

func maybePanic(err error) {
	if err != nil {
		if resErr, ok := err.(*tent.ResponseError); ok && resErr.TentError != nil {
			fmt.Println(resErr.TentError)
		}
		panic(err)
	}
}
