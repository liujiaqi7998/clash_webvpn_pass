<h1 align="center">
  <img src="https://github.com/Dreamacro/clash/raw/master/docs/logo.png" alt="Clash" width="200">
  <br>Clash<br>
</h1>

<h4 align="center">A rule-based tunnel in Go.</h4>

<p align="center">
  <a href="https://github.com/Dreamacro/clash/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/Dreamacro/clash/release.yml?branch=master&style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/Dreamacro/clash">
    <img src="https://goreportcard.com/badge/github.com/Dreamacro/clash?style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/Dreamacro/clash?style=flat-square">
  <a href="https://github.com/Dreamacro/clash/releases">
    <img src="https://img.shields.io/github/release/Dreamacro/clash/all.svg?style=flat-square">
  </a>
  <a href="https://github.com/Dreamacro/clash/releases/tag/premium">
    <img src="https://img.shields.io/badge/release-Premium-00b4f0?style=flat-square">
  </a>
</p>

## 本分支额外新增功能

本 clash 在通过修改VMess的ws和http混淆的方法实现借助某瑞达科技webvpn穿透园区网络，实现不登录网络账号的情况下访问互联网

一些技术细节：

URL加密参考了 [webvpn-dlut](https://github.com/ESWZY/webvpn-dlut)

配置方法：
```yml
Webvpn:
    enable: true # 使能 Webvpn
    Server: "10.1.1.1" # Webvpn 内网服务器地址
    Host: "webvpn.xxx.edu.cn" # Webvpn 外网接收域名
    Port: 443 # Webvpn 端口
    Tls: true # 访问 Webvpn 是否使用TLS加密（https）
    Cookie: "" # 访问 Webvpn 登录用的Cookie
```
还需要在```proxies:```中 VMess 协议添加```WebVpn: true,```开启该功能

感谢： [3181538941](https://github.com/3181538941) 参与开发

## Features

This is a general overview of the features that comes with Clash.  

- Inbound: HTTP, HTTPS, SOCKS5 server, TUN device
- Outbound: Shadowsocks(R), VMess, Trojan, Snell, SOCKS5, HTTP(S), Wireguard
- Rule-based Routing: dynamic scripting, domain, IP addresses, process name and more
- Fake-IP DNS: minimises impact on DNS pollution and improves network performance
- Transparent Proxy: Redirect TCP and TProxy TCP/UDP with automatic route table/rule management
- Proxy Groups: automatic fallback, load balancing or latency testing
- Remote Providers: load remote proxy lists dynamically
- RESTful API: update configuration in-place via a comprehensive API

*Some of the features may only be available in the [Premium core](https://dreamacro.github.io/clash/premium/introduction.html).*

## Documentation

You can find the latest documentation at [https://dreamacro.github.io/clash/](https://dreamacro.github.io/clash/).

## Credits

- [riobard/go-shadowsocks2](https://github.com/riobard/go-shadowsocks2)
- [v2ray/v2ray-core](https://github.com/v2ray/v2ray-core)
- [WireGuard/wireguard-go](https://github.com/WireGuard/wireguard-go)

## License

This software is released under the GPL-3.0 license.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FDreamacro%2Fclash.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FDreamacro%2Fclash?ref=badge_large)
