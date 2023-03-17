# X-UI

X-UI is a webUI panel based on Xray-core which supports multi protocols and multi users  
This project is a fork of [vaxilu&#39;s project](https://github.com/vaxilu/x-ui),and it is a experiental project which used by myself for learning golang   
If you need more language options ,please open a issue and let me know that


# basics

- support system status info check
- support multi protocols and multi users
- support protocols：vmess、vless、trojan、shadowsocks、dokodemo-door、socks、http
- support many transport method including tcp、udp、ws、kcp etc
- traffic counting,traffic restrict and time restrcit
- support custom configuration template
- support https access fot WebUI
- support SSL cert issue by Acme
- support telegram bot notify and control
- more functions in control menu


# installation
Make sure your system `bash` and `curl` and `network` are ready,here we go

```
bash <(curl -Ls https://raw.githubusercontent.com/mazafard/x-ui/master/install.sh)
```  


## shortcut
After Installation，you can input `x-ui`to enter control menu，current menu details：
```
 
  x-ui control menu
  0. exit
————————————————
  1. install   x-ui
  2. update    x-ui
  3. uninstall x-ui
————————————————
  4. reset username
  5. reset panel
  6. reset panel port
  7. check panel info
————————————————
  8. start x-ui
  9. stop  x-ui
  10. restart x-ui
  11. check x-ui status
  12. check x-ui logs
————————————————
  13. enable  x-ui on sysyem startup
  14. disabel x-ui on sysyem startup
————————————————
  15. enable bbr 
  16. issuse certs
 
x-ui status: running
enable on system startup: yes
xray status: running

please input a legal number[0-16]: 
```

## Suggested system as follows:
- CentOS 7+
- Ubuntu 16+
- Debian 8+



# credits

- [vaxilu/x-ui](https://github.com/vaxilu/x-ui)
- [FranzKafkaYu/x-ui](https://github.com/FranzKafkaYu/x-ui/)
- [XTLS/Xray-core](https://github.com/XTLS/Xray-core)
- [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)




## Stargazers over time

[![Stargazers over time](https://starchart.cc/mazafard/x-ui.svg)](https://starchart.cc/mazafard/x-ui)
