# dwebble

## Rewrite

- [ ] Add color to playlist

- [ ] User Tracking
    - [ ] Number of plays 
    - [ ] Favorites
    - [ ] Year over Year 

- [ ] Jobs
    - [ ] Search Indexing
    - [ ] Cleanup Library

- [ ] Better API Tokens

- [ ] Cleanup modal forms 
    - [ ] Errors is hard to see
    - [ ] Handle api errors in a clean way

- [ ] Figure out how to handle logging through out the project
- [ ] Figure out how to handle service errors
- [ ] Figure out how to handle service logging

- [ ] Better way to log the API Errors

- [ ] Edit Playlist CoverURL implementation

- [ ] Web Frontend
    - [ ] Add
        - [ ] Playlist Filtering/Sorting
    - [ ] Finish the design of the pages
    - [ ] Home Page
    - [ ] Better handling of filters, like [Watchbook](https://github.com/Nanoteck137/watchbook)
    - [ ] Update Display Name need fixing, both frontend and backend on /users

- [ ] Cleanup
    - [ ] Migration cleanup
    - [ ] Cleanup database code
    - [ ] Go through all API structures and add all the fields (i.e created, updated, more)

- [ ] Make Search Service init lazily, same as the auth service

- [ ] Media Service
    - [ ] API Route for getTrackStream need to handle errors
    - [ ] Cleanup
    - [ ] Add lock
    - [ ] Add more logging

- [ ] Search Service
    - [ ] Add RWLock lock
    - [ ] Code Cleanup
    - [ ] Add Playlists
    - [ ] Add Users
    - [ ] Batch indexing
    - [ ] Add more logging

- [ ] Job Service
    - [ ] Add lock

- [ ] Library Service
    - [ ] Multi-threaded syncing
    - [ ] Clear the cache after sync
    - [ ] Cleanup after sync (artists, albums and tracks not existing anymore)

- [ ] Library Handling
    - [ ] Update Cmd
        - [ ] Metadata Validation
        - [ ] Metadata transformation (trim spaces, escape characters, more)

- [ ] Media
    - [ ] Fix "packMediaResult"

- [ ] Use database indexes

- [ ] Database seperate track metadata data and track media stuff

- [ ] Docker
    - [ ] Use nix to build a docker image

- [ ] Rename Project
    - Musicbook (MB)
    - Tunebook (TB)
    - Need some ideas

- [ ] Use new logo from dwebble_app

- [ ] Add Compression for Static file routes
    - [ ] SPA Routes
