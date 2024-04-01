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

GO := go
DOCKER := docker

VERSION := 1.0.0
TARGET := ix-device-plugin

COREX_PATH := /usr/local/corex

BUILD_DIR := build

DEPENDS := libixml.so \
           libcuda.so \
           libcuda.so.1 \
           libcudart.so \
           libcudart.so.10.2 \
           libcudart.so.10.2.89 \
           libixthunk.so

IMG_NAME := ix-device-plugin:$(VERSION)

.PHONY: all
all: image

.PHONY: plugin
plugin:
	CGO_CFLAGS=-I$(COREX_PATH)/include \
	$(GO) \
	build \
	-o \
	$(BUILD_DIR)/$(TARGET) \
	cmd/main.go

.PHONY: image
image: plugin
	mkdir -p $(BUILD_DIR)/lib64
	$(foreach lib, $(DEPENDS), cp -P $(COREX_PATH)/lib64/$(lib) $(BUILD_DIR)/lib64;)
	$(DOCKER) build \
		-t $(IMG_NAME) \
		--build-arg EXEC=$(BUILD_DIR)/$(TARGET) \
		--build-arg LIB_DIR=$(BUILD_DIR)/lib64 \
		-f docker/Dockerfile \
	        .

.PHONY: clean
clean:
	rm -fr $(BUILD_DIR)
