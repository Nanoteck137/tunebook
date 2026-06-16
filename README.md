# tunebook

[Logo](https://icon.kitchen/i/H4sIAAAAAAAAA0VSTW-DMAz9K5V3RVVaKG25bf06TZq03qpqCvmASIGgENZWiP8-J5SVA7Kf_Oxnv_TwS3UnWsh6kMX50QjIgGnVUOsggnyCCku5EvWInaYESa0zTWCTOdmusfJNyjwmBAvJPE23I7ImSRyQbZKMSEI2CxgioHWhcUCywjgvvks6KlCWIRyhpp3Rxo4c_wXsIKVgDsdCW1JubgH8opyruvBaUBNki2UEVhUlyvRhbpwz1RhrIQOKM10pKsEhk1S3AqsoL8Rr5FLEqyTs_Cna0rdujKr95EsPd8jIfLmK4DEFbCKu4nR_POB-z6r1VOWD_6r0Pd4fNjBcUQceX9xRFJzDMkcznlfSSmlkwodyzHS1m52s4rO96fJwn9zol3jlqFbsmYaeu6eT3ijhu3Mhaae9jYqZGoGqaxX7qY0T3g30G1V0Fj3oB8wrwzvtX8cFjeLWKO6ZpsX_TeRwHf4AdxWk6EACAAA)

```bash
# Useful command to find every TODO in the project, and then use gF in 
# nvim to goto them
rg --no-heading -in "TODO"

# No preview text
rg --no-heading -in "TODO" | cut -d: -f1,2
```
