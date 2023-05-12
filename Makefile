export_json:
	go run . export ./Teamfile ./Teamfile.json

git_push:
	git checkout main
	git add .
	git commit -m "updated teamfile"
	git push

deploy: export_json git_push