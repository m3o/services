# Dumb Init
# FROM alpine:latest as dumb-init
# RUN apk add --no-cache build-base git bash
# RUN git clone https://github.com/Yelp/dumb-init.git
# WORKDIR /dumb-init
# RUN make

# Copy result to new image 
FROM scratch
# COPY --from=dumb-init /dumb-init/dumb-init /bin/dumb-init
COPY . /bin
# ENTRYPOINT ["dumb-init"]
