download_json:
	go run . download https://raw.githubusercontent.com/gosqueak/leader/main/Teamfile.json ./Teamfile.json

export_json:
	go run . export ./Teamfile ./Teamfile.json

git_push:
	git checkout main
	git add .
	git commit -m "updated teamfile"
	git push

deploy: export_json git_push