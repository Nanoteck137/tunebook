# dwebble

## Rewrite

- [x] Pyrin fix "Internal Server Errors"
    - [x] Print Errors to stderr using logger

- [ ] Better way to log the API Errors

- [ ] Config Handling

- [ ] Migration cleanup

- [ ] Auth Rewrite [AuthLab](https://github.com/Nanoteck137/authlab)
    - [ ] Better API Tokens
- [ ] Server Architecture/Structure
    - [ ] Services

- [ ] Web Frontend
    - [ ] Move away from SvelteKit server and use SvelteKit SPA Mode
        - [x] /login
        - [x] /tracks
        - [x] /taglists
        - [x] /server
        - [ ] /search
        - [x] /playlists
        - [x] /artists
        - [x] /albums
        - [ ] /account
    - [ ] Add
        - [ ] Playlist Editing
        - [ ] Playlist Filtering/Sorting
        - [ ] Taglist Editing
        - [ ] Taglist Filtering/Sorting
        - [ ] On /albums/:id use diffrent colors on odd tracks
        - [ ] On /albums/:id hover highlight tracks
    - [ ] Finish the design of the pages
    - [ ] Add Quick Code login UI
    - [ ] Better handling of filters, like [Watchbook](https://github.com/Nanoteck137/watchbook)

- [x] Quick Playlist need fixing
    - [x] Backend
    - [x] Frontend

- [ ] Media
    - [ ] Better management of media and transcoding
    - [ ] More options for transcoding
    - [ ] Caching settings

- [ ] Playlists/Taglists
    - [ ] Taglists renamed to SmartPlaylists
    - [ ] Visibility
    - [ ] Playlist Tracks should be ordered and then be re-ordered by the user
    - [ ] Playlist should have cover image
        - [ ] Generate the cover image (like Youtube Music)
        - [ ] Custom Covers

- [ ] Library Handling
    - [ ] Faster syncing
    - [ ] Metadata Validation
    - [ ] Metadata transformation (trim spaces, escape characters, more)
    - [ ] Artist handling
        - [ ] When syncing have a flag to error out on the track/album when a unknown artist is found, so that the user can add the infomation for that artist

- [ ] Graceful Shutdown

- [ ] Jobs
    - [ ] Auth Cleanup
    - [ ] Clear out cache
    - [ ] Create new search index
    - [ ] Library Syncing

- [ ] User Tracking
    - [ ] Number of plays 
    - [ ] Favorites
    - [ ] Year over Year 

- [ ] Server Handling for Admins
    - [ ] Notifications
    - [ ] SSE Events

- [ ] Setup Process

- [ ] Docker
    - [ ] Use nix to build a docker image

- [ ] Better Search
    - Use [Bleve](https://blevesearch.com)

- [ ] Rename Project
    - Musicbook (MB)
    - Tunebook (TB)
    - Need some ideas

- [ ] Use new logo from dwebble_app
