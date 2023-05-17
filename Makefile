export_json:
	go run . export ./Teamfile ./Teamfile.json

git_push:
	git checkout main
	git add Teamfile Teamfile.json
	git commit -m "updated teamfiles"
	git push


deploy_teamfile: export_json git_push