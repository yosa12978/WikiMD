
<img src="./static/wikimd.png"/>

A simple wiki based on markdown markup

[![GitHub issues](https://img.shields.io/github/issues/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD/issues)
[![GitHub forks](https://img.shields.io/github/forks/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD/network)
[![GitHub stars](https://img.shields.io/github/stars/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD/stargazers)
[![GitHub license](https://img.shields.io/github/license/yosa12978/WikiMD)](https://github.com/yosa12978/WikiMD)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/yosa12978/WikiMD?include_prereleases)
![GitHub all releases](https://img.shields.io/github/downloads/yosa12978/WikiMD/total)
## Table of contents
- [Table of contents](#table-of-contents)
- [Installation](#installation)
- [Contact](#contact)
- [License](#license)


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

_**I have been released version for windows x64 only**_

_**You can compile source code to any os**_

## Contact

Email: <yusuf_yakubov@hotmail.com>

## License

This repository licensed under GNU GPLv3 license. View [License](LICENSE) file for more details

Copyright © Yusuf Yakubov, 2023
