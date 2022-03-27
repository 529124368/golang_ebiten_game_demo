# 通过golang的ebiten引擎做的类传奇小demo

！转载请注明来源

```
1. 下载源码
2. go mod vendor 
3. go run main.go  或   go build
```
## 控制说明
  1.目前demo 支持 上 下 左 右 按键移动 和 鼠标右键点击移动。
  
### 更新履历
  1. 左上角信息提示追加
  2. 实现鼠标移动
  3. 攻击动作追加和UI追加
  4. 游戏分辨率提升＋ 游戏素材二进制打包（为了支持浏览器运行）使用embed
  5. 优化重复代码，提高代码健壮性和可维护性
  6. 增加武器和技能显示
  7. 优化代码+减少map使用+使用goroutine协程加载素材
  8. 添加玩家和怪物类定义 (持续更新代码结构中........ last date 2022.03.27)
  
## 运行效果如下
   
  web版   体验地址 http://www.zimuge.tk/index.html ※服务器比较垃圾，初次加载需要1分钟左右
   
  ![2](https://user-images.githubusercontent.com/22612129/160224243-73f635a5-976d-4098-9e1b-a3940831ae79.png)

  pc版  
   
  ![1](https://user-images.githubusercontent.com/22612129/160224182-a6908e4d-fa3e-406d-a276-67c09648d729.png)
  
  MacOs
  
  <img width="1180" alt="スクリーンショット 2022-03-26 23 13 42" src="https://user-images.githubusercontent.com/22612129/160243441-cd8d3de7-e7fc-46ef-b607-00ee1a62414d.png">


   



