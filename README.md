# 🔍 zoholab

zoholab is a go libary that can be used to call the zoho API. The motivation to do this project stems from an issue we ran into while developing a way to track data in zoho analytics via their provided go client.
While the client does work there were major improvements that needed to be made for us to be able to use it in our code. This is an improved take of this client with the feartures that we currently need, made with the intent of improving the go client for zoho.

## 💡 Improvements over offical zoho go client

* Added a service struct
* Replace net/http calls with resty calls
* Remove panics from code
* Replace throwing of ServerExceptions with standard error handling
* Removed global variables

## 👷‍♀️ Features

Currently you are able to add a row to a table within zoho analytics. Some new functions might follow if needed.

## 🤯 Usage

Before using: Follow the steps in the [zoho analytics documentation](https://www.zoho.com/analytics/api/#prerequisites) to create a new client and your refresh token. Make sure you create a table that you want to add your data into.

### ✅ You will need
* Client ID
* Client Secret
* Refresh token
* Email address in which the workspace is in
* Workspace name
* Table name

Make sure the client that you have created has the correct permissions to access the table that you want to modify.

### Installation

```
go get github.com/Clarilab/zoholab
```

### 🛳 Import

```go
import "github.com/Clarilab/zoholab"
```

### 😍 Initilise zoholab

```go
func main() {
	authTokenMiddleware := zoholab.middlewares.NewAuthTokenMiddleware(clientid, clientsecret, refreshtoken)

	restyClient := resty.New().OnBeforeRequest(authTokenMiddleware.AddAuthTokenToRequest)

	zohoService := zoholab.NewZohoService(
		restyClient,
		clientid,
		clientsecret,
		refreshtoken,
	)
}
 ```

### 🚣  Code Example for adding a row

```go
func addRow(zohoService *zoholab.ZohoService) error {
	const errMessage = "could not add row"

	url := zoholab.GetUri(email, workspace, tbname)

	columnvalues := map[string]string{
		"your column name": "your column entry",
	}

	addedrows, err := zoholab.AddRow(url, columnvalues)
	if err != nil {
		return errors.Wrap(err, errMessage)
	}

	return nil
}
 ```

* Call this function in your main.go and pass in the zoho instance.

```go
addRow(zohoService)
```

* This is just an example you do not need to make an extra function to call in your main.go, just makes it a little bit cleaner.