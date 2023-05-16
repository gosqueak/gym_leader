# leader
Gym Leader facilitates modelling a distributed system and distrubuting the model to each service.

Create a Teamfile with the following syntax. Whitespace and semicolons are completely ignored when parsing.
```txt
steelix (
    dependencies []
    url (scheme "http" domain "localhost" port "8081")
    listenAddress ":8081"
    endpoints (
        /jwtkeypub (
            methods [GET]
        )
        /register (
            methods [POST]
        )
    )
    jwtInfo (
        issuerName "steelix"
        audienceName "steelix"
    )
)

klefki (
    dependencies [steelix]
    url (scheme "http" domain "localhost" port "8082")
    listenAddress ":8082"
    endpoints (
        / (
            methods [GET, PATCH, DELETE]
        )
    )
    jwtInfo (
        audienceName "klefki"
    )
)
```

# Export a Teamfile to JSON
```go
tm := team.Load("./Teamfile")
tm.SaveJSON("./Teamfile.json")
```

# Download team config from JSON
```go
tm := team.Download(jsonTeamfileURL)
klefkiConfig := tm["klefi"]
fmt.Println(klefkiConfig.Endpoints)
```
