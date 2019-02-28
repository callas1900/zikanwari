# zikanwari

## Installation

### using repository

```
go ${your_go_path}/src
git clone https://github.com/callas1900/zikanwari.git
cd zikanwari
go install
```
then create configuration file as below
```
mkdir ~/.zikanwari
echo "datajsonname: data.json" ~/.zikanwari/config.yaml
```

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
