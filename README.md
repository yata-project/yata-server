# yata-server

## Development

### One Time Setup

#### AWS Config

The yata server needs AWS credentials in order to call various AWS services as
part of normal operation. Create an IAM user in your AWS account with
administrative access and then create a profile called "yata" locally with
`aws configure --profile yata`. See the "Advanced Configuration" section to
customize the profile name.

#### Setting Cognito config

The yata server uses Cognito to authenticate requests.

1. Create a new Cognito user pool in your AWS account.
   1. Under the "General settings" navigation select "App clients" and create an
      app client. Make note of the app client id.
   1. Under the "App integration" navigation
      1. Select "App client settings", locate the app client you created
         earlier:
         1. Check "Cognito User Pool" under "Enabled Identity Providers"
         1. Enter `http://localhost:8888` as the Callback URL(s).
         1. Check "Implicit grant" under "Allowed OAuth Flows".
         1. Check "openid" under "Allowed OAuth Scopes".
         1. Click save; the "Launch Hosted UI" link should now be available,
            click it.
      1. Select "Domain name" and give your user pool a domain prefix. Make note
         of the domain name.
1. Copy the file in `env/SampleConfig.json` to `env/CognitoConfig.json` and
   modify it to suite your setup.

See the "Advanced Configuration" section to customize the filename.

#### Creating a Cognito User

To use meaningfully interact with the server you need to authenticate. To
authenticate you need a user to login in with.

1. Go back to the "App client settings" page of your Cognito user pool.
1. Locate the "Launch Hosted UI" link; click it.
1. Click the "Sign Up" link at the bottom of the UI.
1. Enter a desired username, email, and password.
1. Complete email-based verification.

into Using the domain name you noted earlier go to the following URL to view
your Cognito hosted UI to create a user.

```
https://<your_domain>/login?response_type=token&client_id=<your_app_client_id>&redirect_uri=<your_callback_url>
```

#### Creating DynamoDB Tables

The server uses DynamoDB tables to store items and lists. You can either create
them by deploying https://github.com/TheYeung1/yata-infrastructure or creating
them manually:

1. Create a table called `ListTable`.
   1. With a partition key called `UserID` that's a `String`.
   1. With a sort key called `ListID` that's a `String`.
   1. Uncheck `Use default settings` and change the table to use `On-demand`
      capacity mode. Leave all other settings untouched.
1. Create a table called `ListTable`.
   1. With a partition key called `UserID` that's a `String`.
   1. With a sort key called `ListID-ItemID` that's a `String`.
   1. Uncheck `Use default settings` and change the table to use `On-demand`
      capacity mode. Leave all other settings untouched.

See the "Advanced Configuration" section to customize the table names.

### Everyday

### Getting a JWT token

1. Go back to the "App client settings" page of your Cognito user pool.
1. Locate the "Launch Hosted UI" link; click it.
1. Login in with the user you created earlier.
1. You will be redirected back to the following URL:
   ```
   http://localhost:8888/#id_token=TOKEN&access_token=TOKEN&expires_in=3600&token_type=Bearer
   ```
1. Copy the `TOKEN` to your clipboard.

### Running the Server

1. `go run main.go`
1. The server will start and be available at http://localhost:8888.

### Hitting an API Endpoint

```
curl -H "Authorization: Bearer TOKEN" http://localhost:8888/items
```

Where the `TOKEN` is the same `TOKEN` you retrieved when getting the JWT token
earlier.

### Advanced Configuration

The yata server uses a series of optional command line flags to configure
itself. Run `go run main.go --help` to view these flags and their default
values.
