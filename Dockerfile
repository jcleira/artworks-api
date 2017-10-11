FROM scratch
LABEL maintainer "jmc.leira@gmail.com"

# Adding application minimum runtime files
# artworks-api - service binary
# dbconfig.yml - database connection settings
WORKDIR /app

ADD artworks-api /app/
ADD dbconfig.yml /app/

# TODO - this Dockefile has the environment hardcoded to preproduction
ENTRYPOINT ["./artworks-api", "-environment", "preproduction"]

# This service port. we hardcode ports with an autoincrement from 3000
EXPOSE 3000
