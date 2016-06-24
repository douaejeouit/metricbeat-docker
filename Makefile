BEATNAME=metricbeat-docker
BEAT_DIR=github.com/ingensi/dockerbeat-docker
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS=${GOPATH}/src/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/metricbeat/Makefile

# Initial beat setup
.PHONY: setup
setup:
	make create-metricset
	make collect

.PHONY: update-deps
update-deps:
	glide update --no-recursive --strip-vcs

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

# Create binary packages for the beat
pack: create-packer
	cd dev-tools/packer; make deps images metricbeat-docker
