# yata-server

## Getting JWT token
The following domain takes you to the Hosted UI to log in and get a JWT token.
```
https://<your_domain>/login?response_type=token&client_id=<your_app_client_id>&redirect_uri=<your_callback_url>
```

The token will by after the `#idtoken=` in return URL. For example:
```
https://www.example.com/#id_token=123456789tokens123456789&expires_in=3600&token_type=Bearer  
```