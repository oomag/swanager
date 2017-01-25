# API
Base api path is `/api/`

## Resources

Path | Description
----- | ------
[/session](https://github.com/da4nik/swanager/blob/master/docs/session.md) | Login/Logout
[/apps](https://github.com/da4nik/swanager/blob/master/docs/app.md) | Applications
[/services](https://github.com/da4nik/swanager/blob/master/docs/services.md) | Services
[/users](https://github.com/da4nik/swanager/blob/master/docs/users.md) | Users


# Authorization / Authentication

Endpoint | Method | Params | Description
------ | ------ | ---- | ------
[/session]() | POST | email, password | Login
[/session]() | DELETE | | Logout, need to be authenticated

# Manage applications

Endpoint | Method | Params | Description
------ | ------ | ---- | ------
[/apps]() | GET |  | Application list
[/apps]() | POST | | Create application
[/apps/:app_id]() | GET | | Show application spec and status
[/apps/:app_id]() | PUT | | Update application
[/apps/:app_id/start]() | PUT | | Start application
[/apps/:app_id/stop]() | PUT | | Stop application

# Manage services

Endpoint | Method | Params | Description
------ | ------ | ---- | ------
[/services]() <br /> [/apps/:app_id/services]() | GET | app_id | list app's services
[/services/:service_id]() <br /> [/apps/:app_id/services/:service_id]()  | GET | app_id | inspect service
[/services]() <br /> [/apps/:app_id/services]() | POST | app_id | create app service

# Manage users


Endpoint | Method | Params | Description
------ | ------ | ---- | ------
[/users]() | GET |  | User list
[/users/:user_id]() | GET | user_id  | Show user info 
