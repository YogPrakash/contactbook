## ************************ For Docker-compose setup / Setting up the environment for project ****************************##
1. In machine set the environment variable
#command : export PORT=8080


##***************************************** For testing ***************************************************##
1. install the coverage package  run following commands in terminal
 2.1 - go get golang.org/x/tools/cmd/cover

 #go to the file/directory where your _test.go files are available and run the command given below
 2.2 - go test -v
 #for seeing the test coverage run this command
 2.3 go test -coverprofile  coverage

##***************************************** For authentication **********************************************##
1. Used Basic Auth
-username : username
-password : password

##***************************************** Middlewares Detail ******************************************************##

1.Defined the **type Adapter func(http.Handler) http.Handler**  which  is a function that both takes in and returns an http.Handler.
This is the essence of the wrapper; we will pass in an existing http.Handler, the Adapter will adapt it,
and return a new http.Handler for us to use in its place.

2. To utilise the adaptor functionality we have a  Adapt function which takes the handler you want to adapt, and a list of our Adapter types.
This Adapt function will simply iterate over all adapters, calling them one by one (in reverse order) in a chained manner,
 returning the result of the first adapter.

3.In diffrent API's we are going to do query to mongo db . So each time when any route is going to hit ,each time will gonna create a
mongo session for that. To avoid the mongodb session creation each time i am adapting the mongodb session to middleware which is one time creation and
sending the copy of session to each handler using **WithDB** adaptor.

4. Like db session , for authentication doing the same to avoid the API call to reach the handler function .All authentication will gonna happen at
 middleware and if fails it will return from there only. to achieve this i have **BasicAuth** adaptor.
 4.1 There is some api's which doesn't need authentication . As per our middleware architecture all api's call will first go to BasicAuth
 adaptor and if we don't take care it will fail the api by doing the auth check. So to avoid that i have made the check in basicAuth adaptor if particular
 url matches only then go for authentication .

5. Now we have adap function , withDB ,BasicAuth . To define the handler function for all api endpoints i have **MakeHandler** which
will define the all api' end point and type .

6. Will add this  **MakeHandler** to adap function along with **WithDB** and **BasicAuth** which will do the basic authentication and set the
copy of db session to requests context which we can use further for all other routes .And doing this in the end it will return the handler.

7. For routing used NewServeMux ,which allocates and returns a new ServeMux. ServeMux Handle registers the handler for the given pattern.