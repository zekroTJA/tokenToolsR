<div align="center">
     <!-- <img src="https://zekro.de/src/go_chat_logo.png" width="400"/> -->
     <h1>~ tokenTools R ~</h1>
     <strong>Rewrite of <a href="https://github.com/zekroTJA/discordTokenTools" target="_blank">discordTokenTools</a> - fast & fancy</strong><br><br>
     <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
     <!-- <img src="https://forthebadge.com/images/badges/60-percent-of-the-time-works-every-time.svg" height="30" />&nbsp; -->
     <!-- <a href="https://travis-ci.org/zekroTJA/thunder"><img src="https://img.shields.io/travis/zekroTJA/thunder.svg?style=for-the-badge&logo=travis" height="30"></a>&nbsp; -->
     <!-- <a href="https://godoc.org/github.com/zekroTJA/thunder"><img src="https://img.shields.io/badge/docs-godoc-0ee6ea.svg?style=for-the-badge" height="30"></a>&nbsp; -->
     <a href="https://tokentools.zekro.de"><img src="https://forthebadge.com/images/badges/check-it-out.svg" height="30"></a>&nbsp;
     <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
</div>

---

<div align="center">
<a style="font-size: 22px;" href="https://tokentools.zekro.de"><img src="https://img.shields.io/badge/LIVE%20DEMO-CHECK%20IT%20OUT-%234DD0E1.svg?style=for-the-badge"/></a>
</div>

---

# Introduction

`tokenTools R` is the rewrite of the old Discord tokenTools service. Now, directly accessing the discord database without connecting to the gateway, it is way faster than before. Also, the design has changed to a more modern an dynamic style.

With this tool, you can check the validity of Discord **bot** tokens (*not user tokens, because it is **very** dangerous to handle with this tokens publicly!*), get informations about ther accounts and the number and details of the servers the bot account is connected to.

![](https://i.imgur.com/OBdHRbl.png)
![](https://i.imgur.com/Sb5QNEM.png)
![](https://i.imgur.com/XRXtbQq.png)

*Screenshots of build from 2020/04/06 (255a51e42)*

> **DISCLAIMER:**  
> ONLY USE THIS TOOL FOR PRIVATE PURPOSE AND PERSONALLY CREATED TOKENS YOU ARE THE OWNER OF. PLEASE ONLY USE THIS TOOL WITH FOREIGN TOKENS TO CHECK THEIR VALIDITY AND WARN THE OWNER THAT THE TOKEN IS PUBLIC OR TO WARN SERVER OWNERS THAT THEY SHOULD REMOVE THE BOT ACCOUNT!

---

# Building

1. Clone the repo  
```
$ git clone https://github.com/zekroTJA/tokenToolsR
```

2. Configue the build script for the platform, you want to distribute to:  
> `build.sh`
```bash
#!/bin/bash

## CUSTOM BUILD VARS ##
OS=linux
ARCH=amd64
#######################
```

3. Use the script to build:  
```
$ bash build.sh
```

4. Move the build `tokenTools` and the `web/` folder to the location you want to host it from. Keep in mind that the `web/` folder **needs to be in the same location as the build binary!**

5. Run the tool like follwoing:
```
$ ./tokenTools -port 5002
```
If you do not specify any port, `80` will be used! For this, you need to execute the binary with sudo rights (on linux)!

---

# Used 3rd Party Packages and APIs

- [gorilla/websocket](http://www.gorillatoolkit.org/pkg/websocket)
- [Discord API](https://discordapp.com/)

---

© 2018 zekro Development  

[zekro.de](https://zekro.de) | contact[at]zekro.de


