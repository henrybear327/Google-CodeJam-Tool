# API list

There are 3 main routes: `/`, `dashboard`, and `scoreboard`.

* Get all contest information `https://codejam.googleapis.com/poll?p=e30`
* Get a specific contest's information `https://codejam.googleapis.com/dashboard/0000000000051705/poll?p=e30`
* Get scoreboard data `https://codejam.googleapis.com/scoreboard/%s/poll?p=%s`
* Get a specific handle's data `https://codejam.googleapis.com/scoreboard/%s/find?p=%s`

# Notes

* The query payload is attached to the URL. The format of it is in JSON, encoded in base64.
* The response is also in JSON, encoded in base64