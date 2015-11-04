
format:
	for f in $(find . -type f -name "*.go"); do
		echo "GET $f"
	done


