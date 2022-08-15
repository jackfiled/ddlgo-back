import user_manager

manager = user_manager.Manager()

manager.header = {
    "Authorization": "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvb3QiLCJjbGFzc25hbWUiOiJkZGRkIiwicGVybWlzc2lvbiI6MiwiZXhwIjoxNjYwNjM2NzM1LCJpc3MiOiJTcXVpZFdhcmQifQ.irxMPA7eVP8hm_1DReQyXAxfvt8gebxXwxYouSFDD0Y"
}

def upload_excel():
    users = user_manager.parse_excel("permission.xlsx")

    for user in users:
        manager.createUser(user)

def get_users():
    users = manager.readUsers()

    for user in users:
        print(user.toString())
