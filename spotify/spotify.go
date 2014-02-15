package spotify

import (
	"encoding/json"
	_ "fmt"
	"github.com/hendrikcech/tent-scrobbler/track"
	"os/exec"
	"time"
)

func CurrentTrack() (song track.Track, err error) {
	// check if spotify is running
	cmd := exec.Command("pgrep", "-x", "Spotify")
	isRunning := cmd.Run()
	if isRunning != nil {
		return track.Track{}, nil
	}

	cmd = exec.Command("osascript", "-e", script)
	JSONTrack, err := cmd.Output()
	if err != nil {
		return track.Track{}, nil
	}

	var res map[string]interface{}

	err = json.Unmarshal(JSONTrack, &res)
	if err != nil {
		return
	}

	if res["state"] != "playing" {
		return
	}

	song = track.Track{
		ID:        res["id"].(string),
		Name:      res["name"].(string),
		Album:     res["album"].(string),
		Artist:    res["artist"].(string),
		Remaining: time.Duration(res["remaining"].(float64)) * time.Second,
		Duration:  time.Duration(res["duration"].(float64)) * time.Second,
	}

	return
}

const script string = `
on escape_quotes(string_to_escape)
	set AppleScript's text item delimiters to the "\""
	set the item_list to every text item of string_to_escape
	set AppleScript's text item delimiters to the "\\\""
	set string_to_escape to the item_list as string
	set AppleScript's text item delimiters to ""
	return string_to_escape
end escape_quotes

tell application "Spotify"
	set ctrack to "{"
	set ctrack to ctrack & "\"id\": \"" & current track's id & "\""
	set ctrack to ctrack & ",\"name\": \"" & my escape_quotes(current track's name) & "\""
	set ctrack to ctrack & ",\"album\": \"" & my escape_quotes(current track's album) & "\""
	set ctrack to ctrack & ",\"artist\": \"" & my escape_quotes(current track's artist) & "\""
	set ctrack to ctrack & ",\"remaining\": " & ((current track's duration) - (round (player position as real)))
	set ctrack to ctrack & ",\"duration\": " & current track's duration
	set ctrack to ctrack & ",\"state\": \"" & player state & "\""
	set ctrack to ctrack & "}"
end tell`
