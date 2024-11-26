COVERAGE_DIR ?= .coverage

# cp from: https://github.com/yyle88/syncmap/blob/5ad6d0d6a16cbef3f4b5c5d5e6b56b46cb55a16a/Makefile#L4
test:
	@-rm -r $(COVERAGE_DIR)
	@mkdir $(COVERAGE_DIR)
	make test-with-flags TEST_FLAGS='-v -race -covermode atomic -coverprofile $$(COVERAGE_DIR)/combined.txt -bench=. -benchmem -timeout 20m'

test-with-flags:
	@go test $(TEST_FLAGS) ./...
