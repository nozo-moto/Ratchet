# Ratchet

## Install

```bash
go get github.com/nozo-moto/Ratchet
```

## Usage

command

```bash
Ratchet --project $PROJECT k --instance $INSTANCE --database $DATABASE --credentials_file $YOUR_CREDENTIALS_FILE queryplan "SELECT * FROM hoge"
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
