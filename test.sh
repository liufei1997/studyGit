#!/bin/bash
#framework,serverapi系列采用master分支
# 删除 git.in.codoon.com/third/go-sql-driver/mysql
# common，gorm指定用feature/gorm_add_new_feature分支
# 批量替换gorm.RecordNotFound为gorm.ErrRecordNotFound
#awk '/^$/ {print "Blank line"}' go.mod
#awk '/^$/ {print "Blank line"}' go.mod
# 	git.in.codoon.com/backend/framework v0.0.0-20210705055424-a9a890125198
#awk 'sub("git.in.codoon.com/backend/framework","git.in.codoon.com/backend/framework master")' go.mod
#grep "git.in.codoon.com/third/go-sql-driver/mysql" go.mod
#echo ${#master}
#awk '{ sub(/AAAA/,"BBBB"); print $0 }' a
#sed -i "s/AAAA/6666/g” grep "AAAA" -rl ./a
#sed 's/^#.*//'  a

#LC_CTYPE=C  sed -i '' 's/[0-9]*//g' cs.csv
#
#newFramework="git.in.codoon.com/backend/framework master"

#LC_CTYPE=C sed -i '' "/^6/cQ  cs.csv

#variable="v1"
#sed -i '' '/^git.in.codoon.com\/backend\/framework*/c qqq/' cs.csv
#v1=ewe
#sed -i '' 's/^git.in.codoon.com\/backend\/framework.*/git.in.codoon.com\/backend\/framework master/' cs.csv
#sed -i '' 's/^git.in.codoon.com\/backend\/common.*/git.in.codoon.com\/backend\/common master/' cs.csv

#sed -i '' '/aaa/p' cs.csv
# sed -i "s/原文本/替换后的文本/g" `grep -rl 原文本 ./`
#sed -i "s/git.in.codoon.com\/backend\/framework/git.in.codoon.com\/backend\/framework master/g" `grep  "git.in.codoon.com/backend/framework" cs.csv`
#grep "git.in.codoon.com/backend/framework" cs.csv | sed ''

sed -i '' 's/git.in.codoon.com\/backend\/framework.*/git.in.codoon.com\/backend\/framework master/g' cs.csv
sed -i '' 's/git.in.codoon.com\/backend\/serverapi[ ].*/git.in.codoon.com\/backend\/serverapi master/g' cs.csv
sed -i '' 's/git.in.codoon.com\/backend\/common.*/git.in.codoon.com\/backend\/common feature\/gorm_add_new_feature/g' cs.csv
sed -i '' 's/git.in.codoon.com\/third\/gorm.*/git.in.codoon.com\/third\/gorm feature\/gorm_add_new_feature/g' cs.csv