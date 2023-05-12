# leader
Gym Leader facilitates modelling a distributed system and distrubuting the model to each service.

Create a Teamfile with the following syntax. Whitespace and semicolons are completely ignored when parsing.
```txt
steelix (
    dependencies []
    url "0.0.0.0:8081"
    endpoints (
        /jwtkeypub (
            methods [GET]
    )
    jwtInfo (
        issuerName "steelix"
        audienceName "steelix"
    )
)

klefki (
    dependencies [steelix]
    url "0.0.0.0:8083"
    endpoints (
        / (
            methods [GET, PATCH, DELETE]
        )
    )
    jwtInfo (
        audienceName "klefki"
    )
)

alakazam (
    dependencies [steelix]
    url "0.0.0.0:8082"
    endpoints (
        /ws (
            methods [GET]
        )
    )
    jwtInfo (
        audienceName "alakazam"
    )
)
```
