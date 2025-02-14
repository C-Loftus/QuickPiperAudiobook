# We limit parallelism with -p 1 and ignore caching with -count=1
# because we mutate files on disk and need to test the config generation
.PHONY: test
test:
	go test ./... -count=1 -p 1