default: test

test:
	@scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

test-coverage:
	mkdir -p dist
	go test -coverprofile=dist/coverage.out ./...
	go tool cover -html=dist/coverage.out

test-docker:
	docker-compose --version
	docker-compose up --abort-on-container-exit --exit-code-from=go --force-recreate

.PHONY: default test test-docker
