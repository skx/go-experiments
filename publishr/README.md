
Simple Files Server
-------------------

This program allows files to be uploaded, via HTTP-POST, and
later read back. It can be used to store files on the move.

You can upload a file like so:

    curl --header X-Auth:123456 -X POST -F "file=@/etc/motd" http://localhost:8081/upload

Then get it back again like so:

    curl -v http://localhost:8081/get/vO

(A succesful upload will show you the download-URL.)


Authentication
--------------

Uploads are protected by the "X-Auth" header, which contains a TOTP
value based upon a shared secret.

To generate the secret for the server, and view it, run this:

     publishr init
     publishr secret

It is assumed you can import the secret into a google authenticator,
or use a tool to generate a good response.


Building & Deploying
--------------------

There are a few dependencies which must be installed, and chances are you'll need the `libmagic-dev` package to install one of them:

    apt-get install libmagic-dev

    go get -d ./...
    go build .

As for deployment?  It is assumed you'll be hosting this behind a reverse proxy.
