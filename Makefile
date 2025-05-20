# Run templ generation in watch mode to detect all .templ files and
# re-create _templ.txt files on change, then send reload event to browser.
# Default url: http://localhost:7331
templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false

server:
	air \
	--build.cmd "go build -o tmp/bin/main ." \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind-clean:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --clean

tailwind-watch:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

dev:
	make tailwind-clean
	make -j3 tailwind-watch templ server
