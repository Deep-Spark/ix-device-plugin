# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# All Rights Reserved.
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


FROM ubuntu:20.04
RUN mkdir /usr/local/corex

ARG EXEC
ARG LIB_DIR
ARG BIN_DIR
COPY $EXEC /usr/bin
COPY $LIB_DIR /usr/local/corex/lib64
COPY $BIN_DIR /usr/local/corex/bin

ENV LD_LIBRARY_PATH="/usr/local/corex/lib64"
ENV PATH="/usr/local/corex/bin:${PATH}"

LABEL io.k8s.display-name="IX Device Plugin"
LABEL name="IX Device Plugin"
LABEL vendor="Iluvatar Corex"
ARG PLUGIN_VERSION="N/A"
LABEL version=${PLUGIN_VERSION}
LABEL release="N/A"
LABEL summary="IluvatarCorex device plugin for Kubernetes"
LABEL description="See summary"

ENTRYPOINT ["ix-device-plugin"]
