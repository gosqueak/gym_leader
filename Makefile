export_json:
	go run . export ./Teamfile ./Teamfile.json

git_push:
	git checkout main
	git add Teamfile Teamfile.json
	git commit -m "updated teamfiles"
	git push

git_push_dev:
	git checkout dev
	git add Teamfile Teamfile.json
	git commit -m "updated teamfiles"
	git push


deploy_teamfile: export_json git_push
deploy_teamfile_dev: export_json git_push_dev