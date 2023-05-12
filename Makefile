export_json:
	go run .

git_push:
	git checkout main
	git add .
	git commit -m "updated team json"
	git push

deploy: export_json git_push