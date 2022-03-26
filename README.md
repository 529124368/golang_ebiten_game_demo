# 通过golang的ebiten引擎做的类传奇小demo  powerd by ebiten引擎

！转载请注明来源

```
下载源码
1. go mod vendor 
2. go run main.go   运行游戏
```
## 控制说明
  1.目前demo 支持 上 下 左 右 按键移动 和 鼠标右键点击移动。
  
### 更新履历
  1. 左上角信息提示 更新！！
  2. 实现鼠标移动
  3. 攻击动作追加和UI追加
  4. 游戏分辨率提升＋ 游戏素材二进制打包（为了支持浏览器运行）使用embed
  5. 优化重复代码，提高代码健壮性和可维护性
  6. 增加武器和技能显示
  
## 运行效果如下
   
  web版   体验地址 http://www.zimuge.tk/index.html ※服务器比较垃圾，初次加载需要1分钟左右
   
  ![2](https://user-images.githubusercontent.com/22612129/160224243-73f635a5-976d-4098-9e1b-a3940831ae79.png)

  pc版  
   
  ![1](https://user-images.githubusercontent.com/22612129/160224182-a6908e4d-fa3e-406d-a276-67c09648d729.png)

   



