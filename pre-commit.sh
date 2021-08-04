#! /bin/bash
basepath=$(pwd)
echo "now at ${basepath}"
for file in "${basepath}"/*;do
if [[ -d "$file" ]];then
	if [[ "${file##/*/}"x != "vendor"x ]]
	then
	  goimports -w -v "${file##/*/}"
	fi
fi
done

for file in  *.go
do
    goimports -w -v "${file##/*/}"
done
golangci-lint run -c ./.golangci.yml

go mod tidy
