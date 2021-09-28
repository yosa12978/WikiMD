
<img src="./static/wikimd.png"/>

A simple wiki based on markdown markup

[![GitHub issues](https://img.shields.io/github/issues/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD/issues)
[![GitHub forks](https://img.shields.io/github/forks/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD/network)
[![GitHub stars](https://img.shields.io/github/stars/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD/stargazers)
[![GitHub license](https://img.shields.io/github/license/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD)

## Installation
in config.yaml file

```yaml
Wiki:
  name: WikiMD
  desc: a light markdown wikiengine
  icon: /static/logo.png
Server:
  port: 8989
Mongo:
  conn: mongodb://localhost:27017/
  db: wikimd
```
Change name, desc, icon as you want

Change server port, by default it is 8989

Enter your mongodb server uri eg. ```mongodb://user:pass@localhost:27017/```and mongodb database name eg. ```wikimdDb```

You can change logo by placing it in static folder eg. ```/static/logo.png```

## License

This repository licensed under GNU GPLv3 license. View [License](LICENSE) file for more details