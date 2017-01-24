## Authorization / Authentication
---
Authenticaton is processed via `Authorization` header, which should contain session token. Authorization token is got from **POST** `/session` call.


**Authorize**
---
Authorizes user with email and password and returns session token

* **URL:**
  `/session`
  
* **Method:**
  `POST`
  
* **URL Params**  
  None

* **Data json Params**
  ```json
    {
      "email": "your email",
      "password": "your password"
    }
  ```
* **Success Response:**
    * **Code:** 200 <br />
      **Content:** `{"token": "auth_token_here"}`

* **Error Response:**
    * **Code:** 401 Unauthorized
    * **Code:** 400 Bad request
 
* **Sample Call:**
  ```javascript
    $.ajax({
      url: "/api/session",
      dataType: "json",
      type : "POST",
      data: JSON.stringify({
        "email": "aaa@example.com",
        "password": "12345"
      }),
      success : function(r) {
        console.log(r);
      }
    });
  ```

**Cancel session**
---
Deauthenticates current session

* **URL:**
  `/session`
  
* **Method:**
  `DELETE`
  
* **URL Params**  
  None

* **Data json Params**
  None
  
* **Success Response:**
    * **Code:** 200      

* **Error Response:**
    * **Code:** 422 Unprocessable entity
 
* **Sample Call:**
  ```javascript
    $.ajax({
      url: "/api/session",
      dataType: "json",
      type : "DELETE",
      headers: {"Authorization": "session_token"},
      success : function(r) {
        console.log(r);
      }
    });
  ```
