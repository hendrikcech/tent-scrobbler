# tent-scrobbler
Scrobble your played music to your Tent server. Just supports Spotify on OS X for now.

# install
Either build from source or download the compiled version from the [releases page](https://github.com/hendrikcech/tent-scrobbler/releases).

	git clone git@github.com:hendrikcech/tent-scrobbler.git
	go get
	go build

Then run the binary with `-entity https://yourentity.cupcake.is` to initialize the config and populate it. The file will be stored at `~/.tentscrobbler`.

# post scheme
Identifier: `http://cech.im/types/song/v0#`  
Keys:
- name: The name of the song
- album: Album
- artist: Artist

# license
MIT