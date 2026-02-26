# dwebble

## Rewrite

- [ ] Better way to log the API Errors

- [x] Auth Rewrite [AuthLab](https://github.com/Nanoteck137/authlab)
    - [ ] Better API Tokens
    - [ ] Quick Code integration
- [ ] Server Architecture/Structure
    - [ ] Services
        - [ ] Service should have a ping function to check if the service is available

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
    - [ ] Home Page
    - [ ] Add Quick Code login UI
    - [ ] Better handling of filters, like [Watchbook](https://github.com/Nanoteck137/watchbook)

- [ ] Some fields in AuthService can be private

- [ ] Better Search
    - Use [meilisearch](https://www.meilisearch.com)
    - [Bleve](https://blevesearch.com)

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
    - [ ] Multiple Directories
    - [ ] Have error count for albums when errors > 5 then stop the library syncing 
    - [ ] Metadata Validation
    - [ ] Metadata transformation (trim spaces, escape characters, more)
    - [ ] Artist handling
        - [ ] When syncing have a flag to error out on the track/album when a unknown artist is found, so that the user can add the infomation for that artist

- [ ] Add Compression for Static file routes
    - [ ] SPA Routes

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

- [ ] Import the old format

- [ ] Migration cleanup

- [ ] Docker
    - [ ] Use nix to build a docker image

- [ ] Rename Project
    - Musicbook (MB)
    - Tunebook (TB)
    - Need some ideas

- [ ] Use new logo from dwebble_app

- [x] Config Handling

- [x] Pyrin fix "Internal Server Errors"
    - [x] Print Errors to stderr using logger

- [x] Graceful Shutdown

- [x] Quick Playlist need fixing
    - [x] Backend
    - [x] Frontend

### How to handle the library

The problem I am having is how to multi-thread the import process because the 
albums and tracks need to have access to the artists. The problem is that they are created when 
found inside the album but I don't want to have duplicated artists.

#### Offline Caching / Database
One way:
For this I think the best way would to have a offline "caching / database" step 
and then the library sync read that to determine what to do. 
Have a program create a custom format that the library code can easily read
