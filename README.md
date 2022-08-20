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
  "app_port": ":4000",
  "jwt_secret": "MakeBUPTGreatAgain",
  "jwgl_out_time": 0,
  "root_config": {
    "username": "root",
    "password": "123456",
    "class_name": "dddd",
    "student_id": "0000000000",
    "permission": 2
  },
  "use_mysql": false,
  "mysql_config": {
    "username": "admin",
    "password": "123456",
    "address": "localhost:3306",
    "database_name": "name"
  }
}
```

配置文件中各个字段的意义如下：

- `app_port` 服务器运行的端口号
- `jwt_secret` 签发JWT令牌时采用的密钥
- `jwgl_out_time` 请求教务系统时，为避免过于频繁的请求而设置的超时时间
- `root_config` 为方便管理 在程序开始运行时会自动在数据库中创建的管理用户账号
- `use_mysql` 程序是否使用`mysql`作为数据库

`root_config`下述字段的意义如下：

- `username`: 用户名/姓名
- `password`: 密码
- `class_name`: 所属班级
- `student_id`: 学号
- `permission`: 权限

`mysql_config`下述字段的意义如下：

- `username`: 数据库用户名
- `password`: 对应用户的密码
- `address`: 数据库所在的地址
- `database_name`: 数据库名称

> 在程序中，班级可被设置为"dddd"以表示大班或者"304"~"309"以表示各小班
> 
> 在程序中，人员的权限可被设置为如下三个层级：
> 
> - 0 普通用户 只能查看而无法修改任何内容
> - 1 班级管理员 可以修改本班的内容
> - 2 根管理员 可以修改所有的内容
> 
> 在`use_mysql`字段被设置为`false`的状态下，`mysql_config`字段可以不填写

在根目录下不存在`config.json`配置文件的时候，程序会采用下述默认配置运行

```json
{
  "app_port": ":8080",
  "jwt_secret": "MakeBUPTGreatAgain",
  "jwgl_out_time": 24,
  "root_config": {
    "username": "root",
    "password": "123456",
    "classname": "dddd",
    "student_id": "0000000000",
    "permission": 2
  },
  "use_mysql": false
}
```

程序会在当前文件夹中自动创建一个`test.db`的数据库文件，采用`sqlite`作为默认的数据库。

## API文档

使用[Apifox](https://www.apifox.cn/)
产生的文档[链接](https://www.apifox.cn/apidoc/shared-5d0ad1be-c569-466d-9c59-3e4686b7e482/api-33104131)



