# Run templ generation in watch mode to detect all .templ files and
# re-create _templ.txt files on change, then send reload event to browser.
# Default url: http://localhost:7331
templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false

server:
	air \
	--build.cmd "go build -o /tmp/airmain ." \
	--build.bin "/tmp/airmain" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

clean-regenerate:
	tailwindcss -i ./ui/assets/css/input.css -o ./ui/assets/css/output.css --clean
	sqlc generate

tailwind-watch:
	tailwindcss -i ./ui/assets/css/input.css -o ./ui/assets/css/output.css --watch

dev:
	make clean-regenerate
	make -j3 tailwind-watch templ server
