init: clean
	@echo 'Creating new builds...'
	@read -p "Enter version:" version; \
	docker-compose run godo gox --output=build/{{.Dir}}-$$version-{{.OS}}_{{.Arch}}
	
clean: 
	@echo 'Removing old builds...'
	docker-compose run godo rm -rf build
