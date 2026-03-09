# dwebble

## Rewrite

- [ ] Better way to log the API Errors
- [ ] Better API Tokens
- [ ] Quick Code integration
- [ ] Server Architecture/Structure
    - [ ] Services
        - [ ] Service should have a ping function to check if the service is available

- [ ] Web Frontend
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
    - [ ] Update Display Name need fixing, both frontend and backend on /users

- [ ] Cleanup
    - [ ] Migration cleanup
    - [ ] Go through all API structures and add all the fields (i.e created, updated, more)

- [ ] Make Search Service init lazily, same as the auth service

- [ ] Search Service
    - [ ] Add RWLock lock
    - [ ] Code Cleanup
    - [ ] Add Playlists
    - [ ] Add Users
    - [ ] Batch indexing

- [ ] Utils Media Probe: Needs some work to support more media types

- [ ] Media
    - [ ] Better management of media and transcoding
    - [ ] More options for transcoding
    - [ ] Caching settings
    - [ ] Display the transcoding settings on the frontend

- [ ] Add order/orderNum to API tracks

- [ ] Playlists/Taglists
    - [ ] Taglists renamed to VirtualPlaylist
    - [ ] User ability to re-order playlist items
    - [ ] Custom Covers
    - [ ] Virtual playlists for filtering through a "playlist/all tracks in database"
    - [ ] If a playlist is set on a "virtual playlist" then add a option to maybe disable that and use all the available tracks or make it possible to apply a virtual playlist to other playlists

- [ ] Library Handling
    - [ ] Multi-threaded syncing
    - [ ] Clear the cache after sync
    - [ ] Update Cmd
        - [ ] Metadata Validation
        - [ ] Metadata transformation (trim spaces, escape characters, more)


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

- [ ] Import the old format

- [ ] Use database indexes

- [ ] Docker
    - [ ] Use nix to build a docker image

- [ ] Rename Project
    - Musicbook (MB)
    - Tunebook (TB)
    - Need some ideas

- [ ] Use new logo from dwebble_app

- [ ] Add Compression for Static file routes
    - [ ] SPA Routes
