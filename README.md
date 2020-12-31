# yata-server

## AWS Config
In the AWS credentials for your system, make a profile named "yata" to use for the server.

## Setting Cognito config
In the env folder, create a file named `CognitoConfig.json`. Take a peek at `SampleConfig.json` to see the values needed to identify the cognito pool to use. 

## Getting JWT token
The following domain takes you to the Hosted UI to log in and get a JWT token.
```
https://<your_domain>/login?response_type=token&client_id=<your_app_client_id>&redirect_uri=<your_callback_url>
```

The token will by after the `#idtoken=` in return URL. For example:
```
https://www.example.com/#id_token=123456789tokens123456789&expires_in=3600&token_type=Bearer  
```