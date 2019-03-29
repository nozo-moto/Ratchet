# Spanner_QueryPlan

# install 
```
go get github.com/nozo-moto/spanner_queryplan
```

# run sample

command

``` bash
spanner_queryplan --project $PROJECT k --instance $INSTANCE --database $DATABASE --credentials_file $YOUR_CREDENTIALS_FILE queryplan "SELECT * FROM hoge"
```

result 
```
Run QueryPlan
Query:  SELECT * FROM hoge
----------
 Distributed Union
 Local Distributed Union
  Serialize Result:
   TableScan:hoge UserId:=UserId Created:=Created UUID:=UUID Data:=Data
```
