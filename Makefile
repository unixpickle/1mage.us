all: build/server.js build/assets/script.js

build/server.js:
	coffee -c -o build source/server.coffee

build/assets/script.js:
	coffee -c -o build/assets source/script.coffee

clean:
	rm -rf build/*.js
	rm -rf build/assets/script.js
