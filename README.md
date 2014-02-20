# tent-scrobbler
Scrobble your played music to your Tent server. Just supports Spotify on OS X for now.

# install
Either build from source or download the compiled version from the [releases page](https://github.com/hendrikcech/tent-scrobbler/releases).

	git clone git@github.com:hendrikcech/tent-scrobbler.git
	go get
	go build

Run the binary with `-entity https://yourentity.cupcake.is` to initialize the config and populate it. The player adapter defaults to Spotify for OS X. To select another one, start the binary with `-player playername`.  
The file will be stored at `~/.tentscrobbler`.

# player adapter
Currently only Spotify is supported. To make this project work with your favorite music player, start by fork this repository. Duplicate the `example` folder afterwards and add the needed logic to the `CurrentTrack` function.  
Import your new package in `main.go` and add it to [`playerMap`](https://github.com/hendrikcech/tent-scrobbler/blob/master/main.go#L38).  
When you're convinced your adapter is working properly, submit a pull request and I will be happy to merge it :)

# post scheme
Identifier: `http://cech.im/types/song/v0#`  
Keys:
- name: The name of the song
- album: Album
- artist: Artist

# license
MIT