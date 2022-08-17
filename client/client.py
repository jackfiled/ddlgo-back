import user_manager

manager = user_manager.Manager()

manager.header = {
    "Authorization": "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdHVkZW50X2lkIjoiMDAwMDAwMDAwMCIsImNsYXNzbmFtZSI6ImRkZGQiLCJwZXJtaXNzaW9uIjoyLCJleHAiOjE2NjA4MTgxMzIsImlzcyI6IlNxdWlkV2FyZCJ9.dQ81IOrD-nmnNNQB6ZsAjLQAaWzou1hYoUjH3oHRbhE"
}

def upload_excel():
    users = user_manager.parse_excel("permission.xlsx")

    for user in users:
        manager.createUser(user)

def get_users():
    users = manager.readUsers()

    for user in users:
        print(user.toString())
