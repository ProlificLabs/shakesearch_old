build:
	npm run build --prefix web/reactjs
	rm -Rf public/**
	cp -a web/reactjs/build/. ./public/
