# Copyright 2017 The Kubernetes Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM haproxy:1.9.11-alpine
RUN apk --no-cache add socat openssl lua5.3 lua-socket binutils dpkg

# NOTE: This needs to be the armhf version of dumb-init
# dumb-init kindly manages SIGCHLD from forked HAProxy processes
#ARG DUMB_INIT_SHA256=81231da1cd074fdc81af62789fead8641ef3f24b6b07366a1c34e5b059faf363
#RUN wget -O/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64\
# && echo "$DUMB_INIT_SHA256  /dumb-init" | sha256sum -c -\
# && chmod +x /dumb-init \

RUN wget http://ftp.us.debian.org/debian/pool/main/d/dumb-init/dumb-init_1.2.0-1_armhf.deb \
    && dpkg -i dumb-init_1.2.0-1_armhf.deb \
    && mkdir -p /ingress-controller /etc/haproxy/maps

COPY . /

STOPSIGNAL SIGTERM
ENTRYPOINT ["/usr/bin/dumb-init", "--", "/start.sh"]
