# DDLBackend

## 特点

- RESTful风格的DDL管理API

- JWT令牌身份认证

## 安装

> 采用go 1.18 构建

将仓库克隆到本地

```bash
git clone https://github.com/jackfiled/DDLBackend.git
```

切换进入软件目录，使用`go build`编译可执行程序，运行程序。

```bash
cd DDLBackend
go build
./DDLBackend

...
[GIN-debug] Listening and serving HTTP on :8080
```

看见最后一行输出服务器在8080端口上监听说明运行成功。

### 配置文件

程序采用`config.json`文件作为配置文件，该文件的模板如下
```json
{
    // 服务器运行时执行的端口
    "app_port": ":4000",
    // JWT签发密钥时使用的字符串
    "jwt_secret": "MakeBUPTGreatAgain",
    // 请求教务系统API的超时时间
    "jwgl_out_time": 0,
    // 根管理员的设置
    // 该管理员将在程序运行时自动创建
    "root_config": {
        // 用户名
        "username": "root",
        // 密码
        "password": "123456",
        // 所属的班级
        "classname": "dddd",
        // 学号
        "student_id": "0000000000",
        // 权限
        // 0-User 用户 只能够查看信息而不能修改信息
        // 1-Admin 管理员 可以修改自己所在班级的信息
        // 2-Root 根管理员 可以修改所有的信息
        "permission": 2
    }
}
```

在根目录下不存在`config.json`配置文件的时候，程序会采用下述默认配置运行

```json
{
    // 服务器运行时执行的端口
    "app_port": ":8080",
    // JWT签发密钥时使用的字符串
    "jwt_secret": "MakeBUPTGreatAgain",
    // 请求教务系统API的超时时间 单位是小时
    "jwgl_out_time": 24,
    // 根管理员的设置
    // 该管理员将在程序运行时自动创建
    "root_config": {
        // 用户名
        "username": "root",
        // 密码
        "password": "123456",
        // 所属的班级
        "classname": "dddd",
        // 学号
        "student_id": "0000000000",
        // 权限
        // 0-User 用户 只能够查看信息而不能修改信息
        // 1-Admin 管理员 可以修改自己所在班级的信息
        // 2-Root 根管理员 可以修改所有的信息
        "permission": 2
    }
}
```

## API文档

使用[Apifox](https://www.apifox.cn/)产生的文档[链接](https://www.apifox.cn/apidoc/shared-5d0ad1be-c569-466d-9c59-3e4686b7e482/api-33104131)



