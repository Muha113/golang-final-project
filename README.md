# Twitter like web app

## supports GET and POST methods

* ## POST ___/register___
    ## request :
    * username - some unique user name
    * email - user unique email address
    * password - some password
    ## response : 
    * id - primary key
    * username - username which you specified in request
    * email - user email address which you specified in request
* ## POST ___/login___
    ## request :
    * email - user email address
    * password - user password
    ## response :
    * token - json web token
* ## POST ___/subscribe___
    ## request :
    * username - user name for account which you want to subscribe
* ## POST ___/tweets___
    ## request :
    * message - some tweet message
    ## response :
    * id - message id primary key
    * message - tweet message
* ## GET ___/tweets___
    ## response :
    * tweets - all tweets from your subscriptions

## Database : MongoDB
## Run app :
>dep ensure -v

>go run cmd/server/server.go