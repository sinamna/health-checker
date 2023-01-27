# health-checker
You can check health of your endpoints with this beautiful async health-checker :)
## Endpoints

#### user endpoints
```
POST /user/signup
POST /user/login # returns jwt token
```

#### endpoints 
WARNING: there endpoints requires authentication
```
POST /endpoint/ # registers endpoint for user who invoke this endpoint
GET /endpoint/ # returns endpoints of user with their ID
GET /endpoint/:endpoint_id # returns the status of endpoint (number of failed and success checks)
```
#### alerts 
WARNING: there endpoints requires authentication
```
GET /alerts/:endpoint_id # returns the alerts of endpoint
```
