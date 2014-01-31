all: build/server build/assets

build:
	mkdir build

build/assets: build
	mkdir build/assets
	coffee -c -o build/assets source/assets/*.coffee

build/server: build
	mkdir build/server
	coffee -c -o build/server source/server/*.coffee

clean:
	rm -rf build/
