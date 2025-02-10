# We limit parallelism to 1 because we mutate the filesystem to 
# test the config generation
.PHONY: test
test:
	go test ./... -count=1 -p 1