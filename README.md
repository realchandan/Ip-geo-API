# Intro

Hi, there!

Here is a ready to deploy ip geo location server, It works for both ip v4 and ip v6,
However the underlying database is not that huge,
It makes use of the [OpenGeoFeed database](https://github.com/sapics/ip-location-db/blob/master/geo-whois-asn-country/README.md#geofeed-database-update-daily) found [here](https://github.com/sapics/ip-location-db),

# How to use

Here's a command to run using docker ->

```bash
docker run -p 8080:8080 ghcr.io/realchandan/ip-geo-api
```

## Request

```bash
curl localhost:8080/getIpInfo?addr=140.82.114.3
```

## Response

```json
{ "country": "US", "ok": true }
```

# License

The database is licensed under [CC0](https://creativecommons.org/share-your-work/public-domain/cc0/),

It means -

> CC0 doesn't legally require users of the data to cite the source!

But feel free to attribute the ip database [provider](https://opengeofeed.org/)! <3

The code in this repository is licensed under the MIT License!

There are some work yet to be done, like setting up github actions, but the current docker image which is hosted on docker hub is usable!

# Thanks!
