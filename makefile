# We limit parallelism to 1 and ignore caching with 1
# because we mutate files on disk to test the config generation
.PHONY: test
test:
	go test ./... -count=1 -p 1