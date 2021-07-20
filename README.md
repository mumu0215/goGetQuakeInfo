# goGetQuakeInfo
go program using quake api

**Token should be set to line 21**

## keyWords
|     keyWords      |               Example                |
| :---------------: | :----------------------------------: |
|        ip         |             ip:"1.1.1.1"             |
|      domain       |           domain:"360.cn"            |
| services/service  |     services:"rtsp,https,telnet"     |
| products/products |      products:"BusyBox,Apache"       |
|     response      | response:"220 ProFTPD 1.3.5a Server" |
|        app        |    app: "YApi 可视化接口管理平台"    |
|       title       |             title:admin              |

## Unauthorized Service

| seviceName |             searchWords              |
| :--------: | :----------------------------------: |
|    ftp     |    service.ftp.is_anonymous:true     |
|   rsync    |  service.rsync.authentication:false  |
|  mongodb   | service.mongodb.authentication:false |

### tips

excloude honeypots using keywords like `NOT type:"蜜罐"` when searching.

I have put these keywords in my code. **There is no need to set them when searching**.

### reference

https://quake.360.cn/quake/#/help?id=5eb238f110d2e850d5c6aec8&title=%E6%A3%80%E7%B4%A2%E5%85%B3%E9%94%AE%E8%AF%8D