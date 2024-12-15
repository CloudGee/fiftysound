# 五十音学习助手 (FiftySound)

五十音学习助手是一款基于 Go 语言和 Fyne 框架开发的日语学习应用，专注于帮助用户练习平假名和片假名的记忆与书写。

## 功能特点

- **两种学习模式**
  - 模式一：假名 -> 罗马音
  - 模式二：罗马音 -> 假名手写
- **自定义假名范围**：用户可自由选择需要练习的平假名或片假名。
- **动态正确率统计**：实时显示学习过程中的正确率。
- **手写绘图功能**：支持手写输入平假名或片假名。
- **支持多平台**：Windows

---

## 下载与安装

### 直接下载

点击以下链接下载适合您平台的版本：

- [Windows 64位](https://github.com/CloudGee/fiftysound/releases/latest/download/FiftySoundApp_windows_amd64.exe)

下载完成后，按照您的操作系统运行即可。

---

## 使用方式

1. 启动应用程序后，请主动选择学习模式：
   - "模式一: 假名 => 罗马音"
   - "模式二: 罗马音 => 假名手写"

2. 点击"选择假名范围"按钮，弹出选择界面：
   - 可以按行选择五十音图中的某一行（如"あ行"）。
   - 可以逐个勾选具体的假名。注意：必须选择至少两个假名才能开始训练，否则无法进入学习模式。
   - 可以点击"全部随机(包含所有五十音)"，快速选择所有假名。

3. 可以选择单独联系平假名或者片假名，默认是两个都选中，随机出现平假名和片假名的练习题。

4. 选择完成后，返回主界面，点击"开始"按钮进入学习模式。

5. 在学习过程中：
   - 模式一中，输入假名对应的罗马音，点击"判断"按钮查看答案。
   - 模式二中，在绘图区域手写对应的假名，点击"显示答案"查看正确答案进行人工比对。模式二不支持自动判题

6. 学习完成后，可以随时返回主界面调整设置或退出应用。
---

## 开发与构建

### 手动构建

1. 安装依赖：

   ```bash
   go mod tidy
   ```

2. 使用 `fyne-cross` 编译：

   ```bash
   fyne-cross windows -arch amd64 -output FiftySoundApp --app-id com.fiftysound
   fyne-cross darwin -arch amd64 -output FiftySoundApp --app-id com.fiftysound
   fyne-cross linux -arch amd64 -output FiftySoundApp --app-id com.fiftysound
   ```

生成的可执行文件位于 `fyne-cross` 的输出目录下。