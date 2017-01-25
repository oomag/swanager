# Manage services

Endpoint | Method | Params | Description
------ | ------ | ---- | ------
[/services]() <br /> [/apps/:app_id/services]() | GET | app_id | list app's services
[/services/:service_id]() <br /> [/apps/:app_id/services/:service_id]()  | GET | app_id | inspect service
[/services]() <br /> [/apps/:app_id/services]() | POST | app_id | create app service


## List 
Get services list for application

* **URL:**
 
  `/services` - in this case `app_id` should be provided as url param <br />
  `/apps/:app_id/services`
  
* **Method:**
 
  `GET`
  
* **URL Params** 

   **Required:**
   
  `app_id=[string]` - Service's application id

* **Success Response:**
    * **Code:** 200 <br />
      **Content:** `{}`

* **Error Response:**
    * **Code:** 401 Unauthorized
    * **Code:** 404 Not found
 
* **Sample Call:**
  ```javascript
    $.ajax({
      url: "/api/services", // or "/api/apps/<app_id>/services"
      dataType: "json",
      type : "GET",
      data: { app_id: "app_id" },
      success : function(r) {
        console.log(r);
      }
    });
  ```

## Inspect
Inspect service and it's current status

* **URL:**
 
  `/services/:service_id` - in this case `app_id` should be provided as url param <br />
  `/apps/:app_id/services/:service_id`
  
* **Method:**
 
  `GET`
  
* **URL Params** 

   **Required:**
   
  `app_id=[string]` - Service's application id

* **Success Response:**
    * **Code:** 200 <br />
      **Content:** 
      ```json
      {
        "service": {},
        "status": []
      }
      ```

* **Error Response:**
    * **Code:** 401 Unauthorized
    * **Code:** 404 Not found
 
* **Sample Call:**
  ```javascript
    $.ajax({
      url: "/api/services/<service_id>", // or "/api/apps/<app_id>/services"
      dataType: "json",
      type : "GET",
      data: { app_id: "app_id" },
      success : function(r) {
        console.log(r);
      }
    });
  ```
