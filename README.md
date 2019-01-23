# zikanwari

## Installation
```
git clone https://github.com/callas1900/zikanwari.git
cd zikanwari
go build
mv zikanwari ${AS_YOU_WANT}
```
then create configuration file as below
```
cat ~/.zikanwari.yaml
DataJsonPath: ./data.json
```


### requirements
* golang

## How to use

something like below
```
zikanwari init
zikanwari set LUNCH 12:00-13:00
zikanwari set MTG 15:30-16:00
zikanwari show
```

## usage lib
[cobra](https://github.com/spf13/cobra)
