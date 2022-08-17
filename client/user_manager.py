import openpyxl
import requests
import json


class Manager:
    rootUrl: str = "http://localhost:4000"
    header: map

    def __init__(self, root=None) -> None:
        if root == None:
            print("Use the default url {}".format(self.rootUrl))
        else:
            self.rootUrl = root

    def getRootUserToken(self, student_id: str, password: str) -> bool:
        url = self.rootUrl + "/login"
        response = requests.post(url, json={
            "student_id": student_id,
            "password": password
        })
        if response.status_code == 200:
            jsonMap = json.loads(response.text)
            self.header = {
                "Authorization": "Bearer " + jsonMap["token"]
            }
            print(jsonMap["token"])
            return True
        else:
            print(response.text)
            return False
    
    def readUsers(self):
        url = self.rootUrl + "/users"

        response = requests.get(url, headers=self.header)

        if response.status_code == 200:
            lst = json.loads(response.text)
            users = []
            for item in lst:
                user = UserModel()
                user.fromJson(item)
                users.append(user)
            return users
        else:
            print(response.text)
            return
    
    def readUser(self, id):
        url = self.rootUrl + "/users/" + str(id)
        response = requests.get(url, headers=self.header)

        if response.status_code == 200:
            jsonMap = json.loads(response.text)
            user = UserModel()
            user.fromJson(jsonMap)
            return user
        else:
            print(response.text)
        
    def createUser(self, user):
        url = self.rootUrl + "/users"

        response = requests.post(url, json=user.toJson(), headers=self.header)

        if response.status_code == 201:
            jsonMap = json.loads(response.text)
            user = UserModel()
            user.fromJson(jsonMap)
            print(user.toString())
        else:
            print(response.text)
    
    def updateUser(self, user):
        url = self.rootUrl + "/users/" + str(user.id)

        jsonMap = user.toJson()
        jsonMap["id"] = user.id

        response = requests.put(url, json=jsonMap, headers=self.header)

        if response.status_code != 204:
            print(response.text)
    
    def deleteUser(self, id: int):
        url = self.rootUrl + "/users/" + str(id)

        response = requests.delete(url, headers=self.header)

        if response.status_code != 204:
            print(response.text)


class UserModel:
    id: int
    username: str
    password: str
    student_id: str
    permission: int
    class_name: str

    def fromJson(self, json: map):
        self.id = json["id"]
        self.username = json["username"]
        self.password = json["password"]
        self.student_id = json["student_id"]
        self.permission = json["permission"]
        self.class_name = json["classname"]
    
    def toJson(self):
        result = {}
        result["username"] = self.username
        result["student_id"] = self.student_id
        result["permission"] = self.permission
        result["classname"] = self.class_name
        result["password"] = self.password
        return result
    
    def toString(self):
        return """username: {}
student_id: {}
classname: {}
permission: {}""".format(self.username, self.student_id, self.class_name, self.permission)

    def checkSame(self, other: "UserModel"):
        return self.username == other.username and self.student_id == other.student_id and self.class_name == other.class_name and self.permission == other.permission


def parse_excel(path: str)-> list:
    work_book = openpyxl.load_workbook(path)
    work_sheet = work_book.active

    users = []

    for row in work_sheet.values:
        user = UserModel()
        user.class_name = str(row[0] - 2021211000)
        user.student_id = str(row[1])
        user.username = row[2]
        user.permission = row[3]
        user.password = "123456"

        users.append(user)
        
    return users