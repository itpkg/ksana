

format:
	for f in `find . -type f -name "*.go"`; do gofmt -w $$f; done


test:
	for d in utils ioc orm i18n settings cache cmd job logging mux; do cd $$d && go test && cd ..; done



