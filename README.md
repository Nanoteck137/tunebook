# tunebook

[Logo](https://icon.kitchen/i/H4sIAAAAAAAAA0VSTW-DMAz9K5V3RVVaKG25bf06TZq03qpqCvmASIGgENZWiP8-J5SVA7Kf_Oxnv_TwS3UnWsh6kMX50QjIgGnVUOsggnyCCku5EvWInaYESa0zTWCTOdmusfJNyjwmBAvJPE23I7ImSRyQbZKMSEI2CxgioHWhcUCywjgvvks6KlCWIRyhpp3Rxo4c_wXsIKVgDsdCW1JubgH8opyruvBaUBNki2UEVhUlyvRhbpwz1RhrIQOKM10pKsEhk1S3AqsoL8Rr5FLEqyTs_Cna0rdujKr95EsPd8jIfLmK4DEFbCKu4nR_POB-z6r1VOWD_6r0Pd4fNjBcUQceX9xRFJzDMkcznlfSSmlkwodyzHS1m52s4rO96fJwn9zol3jlqFbsmYaeu6eT3ijhu3Mhaae9jYqZGoGqaxX7qY0T3g30G1V0Fj3oB8wrwzvtX8cFjeLWKO6ZpsX_TeRwHf4AdxWk6EACAAA)

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

- [ ] CLI: Fix the import command

- [ ] Backend: Fix: Maybe remove playlist_filters and replace with the global filters?

- [ ] Backend: Feature: Year over year

- [ ] Backend: Fix: API: Media api

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

- [ ] Backend: Media Service: Locking
- [ ] Backend: Media Service: Logging

- [ ] Backend: Job Service: Locking

- [ ] Backend: Database: Use database indexes

- [ ] Backend: Cleanup: Cleanup Jobs from base_app.go

- [ ] Backend: Cleanup: Database: Database code
- [ ] Backend: Cleanup: API: Go through all API structures and add all the fields (i.e created, updated, more)
- [ ] Backend: Cleanup: Media Service: Code Cleanup
- [ ] Backend: Cleanup: Job Service: Code Cleanup
- [ ] Backend: Cleanup: Library Service: Code Cleanup
- [ ] Backend: Cleanup: Auth Service: Code Cleanup

- [ ] Backend: Cleanup: API: Better way to log the API Errors

## Future

- [ ] Backend: Add back migrate command, maybe only for dev
- [ ] Pyrin: Generate Structures for SSE events

- [ ] Backend: API: Add Compression for Static file routes (SPA Routes)

- [ ] Backend: Library Service: Multi-threaded syncing

- [ ] Backend: Fix: Database: Figure out how to handle playlist items/tracks

```sql
CREATE TABLE user_listening_stats (
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id    INTEGER NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    period      TEXT NOT NULL,  -- '2024', '2024-Q1', '2024-03'

    play_count  INTEGER NOT NULL DEFAULT 0,
    skip_count  INTEGER NOT NULL DEFAULT 0,
    total_ms    INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (user_id, track_id, period)
);
```
