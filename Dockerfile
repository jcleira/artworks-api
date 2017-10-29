FROM debian
LABEL maintainer "jmc.leira@gmail.com"

# Adding application minimum runtime files
# artworks-api - service binary
# dbconfig.yml - database connection settings

ADD artworks-api .
ADD dbconfig.yml .

# TODO - this Dockefile has the environment hardcoded to preproduction
ENTRYPOINT ["./artworks-api", "-environment", "preproduction"]

# This service port. we hardcode ports with an autoincrement from 3000
EXPOSE 3000
