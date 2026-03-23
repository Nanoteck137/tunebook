# dwebble

```bash
# Useful command to find every TODO in the project, and then use gF in 
# nvim to goto them
rg --no-heading -in "TODO"

# No preview text
rg --no-heading -in "TODO" | cut -d: -f1,2
```

## Rewrite

```go
api.Error(ErrPlaylistNotFound, 400).Message("Hello World")
```

- [ ] Backend: API: Fix: Media api

- [ ] Backend: Fix: Playlist Service: When deleting tracks we need to reorder the tracks again
- [ ] Backend: Feature: Ability to get all playlist + filters
- [ ] Backend: Fix: Database: Figure out how to handle playlist items/tracks

- [ ] CLI: Feature: Update Cmd: Use library.json
- [ ] CLI: Feature: Init Cmd: Add init library cmd
- [ ] CLI: Feature: Init Cmd: Better logging
- [ ] CLI: Cleanup: Update Cmd: Cleanup Code 
- [ ] CLI: Cleanup: Code Cleanup
- [ ] CLI: Cleanup: Cleanup init commands

- [ ] CLI: Fix: Fix the dwebble migrate commands
- [ ] CLI: Fix: Add some reminders to migrate commands

- [ ] Frontend: Design: Re-design the header, add profile picture + drop down
- [ ] Frontend: Design: Re-design the player UI
- [ ] Frontend: Design: Re-design the album items on /albums
- [ ] Frontend: Design: Re-design /server
- [ ] Frontend: Design: Re-design /users/:id
- [ ] Frontend: Feature: Lazy load the quick playlist menu
- [ ] Frontend: Feature: Add showPlaylistModal
- [ ] Frontend: Design: Home Page
- [ ] Frontend: Design: Re-design the track list component
- [ ] Frontend: Feature: Design: Album filtering
- [ ] Frontend: Feature: Design: Artist filtering
- [ ] Frontend: Feature: Design: Track filtering
- [ ] Frontend: Feature: Design: Playlist Filtering/Sorting
- [ ] Frontend: Fix: Album track count
- [ ] Frontend: Fix: Form errors is hard to see
- [ ] Frontend: Fix: Handle api errors in a clean way
- [ ] Frontend: Fix: Music Manager
- [ ] Frontend: Fix: Small Audio player 

- [ ] Frontend: Cleanup: Modal forms 
- [ ] Frontend: Fix: cancel() on forms on api errors, this makes sure that the form is not reset when error occurs

- [ ] Backend: Feature: User Tracking
- [ ] Backend: Feature: Favorites
- [ ] Backend: Feature: Year over year

- [ ] Backend: Feature: Search indexing job
- [ ] Backend: Feature: Library cleanup job

- [ ] Backend: Cleanup: Figure out how to handle logging through out the project
- [ ] Backend: Cleanup: Figure out how to handle service errors
- [ ] Backend: Cleanup: Figure out how to handle service logging
- [ ] Backend: Cleanup: API: Better way to log the API Errors
- [ ] Backend: Cleanup: Cleanup user picture code
- [ ] Backend: Cleanup: Database: Database code
- [ ] Backend: Cleanup: API: Go through all API structures and add all the fields (i.e created, updated, more)
- [ ] Backend: Cleanup: Search Service: Code Cleanup
- [ ] Backend: Cleanup: Media Service: Code Cleanup
- [ ] Backend: Cleanup: Job Service: Code Cleanup
- [ ] Backend: Cleanup: Library Service: Code Cleanup
- [ ] Backend: Cleanup: Auth Service: Code Cleanup
- [ ] Backend: Cleanup: CLI: Code Cleanup

- [ ] Backend: Search Service: init lazily, same as the auth service
- [ ] Backend: Search Service: RWLock lock
- [ ] Backend: Search Service: Add Playlists
- [ ] Backend: Search Service: Add Users
- [ ] Backend: Search Service: Batch indexing
- [ ] Backend: Search Service: Add more logging
- [ ] Backend: Fix: API: Move every search handler to a /search endpoint

- [ ] Backend: Media Service: API Route for getTrackStream need to handle errors
- [ ] Backend: Media Service: Locking
- [ ] Backend: Media Service: Logging

- [ ] Backend: Job Service: Locking

- [ ] Backend: Library Service: Multi-threaded syncing
- [ ] Backend: Library Service: Clear the cache after sync
- [ ] Backend: Library Service: Cleanup after sync (artists, albums and tracks not existing anymore)

- [ ] Backend: Database: Use database indexes

- [ ] Backend: Fix: SPA Web handling

- [ ] Backend: Cleanup: Cleanup Jobs from base_app.go
- [ ] Backend: Fix: When shutting down, we need to wait for all the jobs to finish

- [ ] General: Use nix to build a docker image

- [ ] General: Rename Project
    - Musicbook (MB)
    - Tunebook (TB)
    - Need some ideas

- [ ] General: Use new logo from dwebble_app

- [ ] Backend: API: Add Compression for Static file routes (SPA Routes)
