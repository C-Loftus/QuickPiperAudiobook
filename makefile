# We limit parallelism with -p 1 and ignore caching with -count=1
# because we mutate files on disk and need to test the config generation
test:
	go test ./... -count=1 -p 1

release:
	git add . 
	git commit -m "release" || true
	git push origin master
	@read -p "Enter tag name: " tag; \
	read -p "Enter tag message: " msg; \
	git tag -a $$tag -m "$$msg"; \
	git push origin $$tag
