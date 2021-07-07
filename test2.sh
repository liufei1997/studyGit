#!/bin/bash



# 下面三个正确
#sed -i '' 's/git.in.codoon.com\/third\/gorm.*/git.in.codoon.com\/third\/gorm feature\/gorm_add_new_feature/g' cs.csv
#sed -i '' 's/git.in.codoon.com\/backend\/common.*/git.in.codoon.com\/backend\/common feature\/gorm_add_new_feature/g' cs.csv
#sed -i '' '/.*git\.in\.codoon\.com\/third\/go\-sql\-driver\/mysql.*/d' cs.csv

#sed -i '' 's/git.in.codoon.com\/backend\/framework.*/git.in.codoon.com\/backend\/framework master/g' cs.csv
#sed -i '' 's/git.in.codoon.com\/backend\/serverapi[ ].*/git.in.codoon.com\/backend\/serverapi master/g' cs.csv


#sed -i '' '/(^\s*)git/d' cs.csv

                    #git.in.codoon.com\/third\/go-sql-driver/mysql111
#sed  -i '' '/[ \t]*git\.in\.codoon\.com\/third\/go\-sql\-driver\/mysql.*/d' cs.csv
#sed  -i '' '/^[ \t]*git.in.codoon.com\/third\/go-sql-driver\/mysql.*/d' testshell
#sed -i '' '/*git.in.codoon.com\/third\/go-sql-driver\/mysql.*/d’ cs.csv

# git.in.codoon.com\/third\/go-sql-driver\/mysql.*

#sed  -i '' '/^[ \t]*a\/b\.-a.*/d' testshell
#sed  -i '' '/^[\t ]*gitlab.*/d' testshell

#sed -i '' '/[ \t]*git\.in\.codoon\.com\/third\/go\-sql\-driver\/mysql.*/d' testshell
#sed -i '' '/.*git\.in\.codoon\.com\/third\/go\-sql\-driver\/mysql.*/d' testshell

sed  -i '' 's/x$//' testshell
#sed -i '' 's/git.in.codoon.com\/backend\/serverapi[\/]*/' testall


function git.branch {
  br=`git branch | grep "*"`
  echo ${br/* /}
}

git.branch
