#!/bin/bash

#项目绝对路径
servicePath="$1"
goModPath="$1/go.mod"

isSuccess() {
  if [ $? != 0 ]; then
    git checkout "$goModPath"
    echo "执行失败已退出"
    exit
  fi
}

replaceAllBranch() {
  array=$(grep "$1" "$3" | awk '{print $2}')
  for i in ${array[*]}; do
    sed -i '' "s/$i/$2/" "$3"
  done
}

delPackage() {
  sed -i '' "/$1/d" "$2"
}

updateGoMod(){
    replaceAllBranch "backend\/serverapi" "master" "$goModPath"
    isSuccess
    replaceAllBranch "backend\/common" "feature\/gorm\_add\_new\_feature" "$goModPath"
    isSuccess
    replaceAllBranch "third\/gorm" "feature\/gorm\_add\_new\_feature" "$goModPath"
    isSuccess
    replaceAllBranch "backend\/framework" "master" "$goModPath"
    isSuccess
    delPackage "third\/go\-sql\-driver\/mysql" "$goModPath"
    isSuccess
}

checkGoModExist() {
  if [[ -f "$goModPath" ]]; then
    echo "go.mod文件存在"
    cd "$servicePath" || exit
    git checkout -b feature/gorm_update_v1_max
    isSuccess
    updateGoMod
    go build
    isSuccess
    echo "更新go.mod文件并编译成功"
  else
    echo "go.mod文件不存在"
    exit
  fi
}

checkGoModExist "$1"
