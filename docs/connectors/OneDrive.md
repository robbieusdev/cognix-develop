#### OneDrive 

From CogniX UI navigate to connectors and create a new connector
Choose OneDrive
At step 2:
- Choose a name, it's just a description
- Fill the "Connector Specific Configration" with the json below filled with the corect data
- Refresh frequency in seconds is the delta of time that CogniX will use to start a new scan on your connected data source
- Connector credential, fill with a random number, it's not used

```json
{
  "folder": "",
  "recursive": false,
  "token": {
    "access_token": "",
    "expiry": "",
    "refresh_token": "",
    "token_type": ""
  }
}


```

json properties: </br>
**folder** </br>
string, optional if not set CogniX will analyze the whole drive<br/>
example: older/chapter1
<br?><br> 
**recursive**  <br/>
bool, (default false), you can omit. It indicates if CogniX shall analyze all the subfolder of the given path or not <br>
**token**  <br/> 
The OAuth token you generate from OneDrive for CogniX to have access to the resource. Below a detailed description on how to get it

Since the UI is still under costruction folling some manual steps to get the OneDrive toke.
This process will be automated with the UI evelution

paste in your browser the following link if you are running CogniX on your private Docker deployment
```js
    http://localhost:8080/api/oauth/microsoft/auth_url?redirect_url=http://localhost:8080
```

If you are using CogniX from [rag.cognix.ch](https://rag.cognix.ch)
```js
    https://rag.cognix.ch:8080/api/oauth/microsoft/auth_url?redirect_url=http://rag.cognix.ch:8080
```

once you paste the link above in the browser you will get a json. copy only the content of data without " taking care not to forget any char at the beginning or at the end of the text <br/>
you will get something similar to the json below:<br/>

```json
{
  "status": 200,
  "data": "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=<id>>&scope=offline_access Files.Read.All Sites.ReadWrite.All&response_type=code&redirect_uri=http://localhost:8080/api/oauth/microsoft/callback"
}
```

paste the link as described above in a new browser window <br/>. 
Sign in using you microsoft account and grant permission to CogniX<br/>
There's a checkbox you need to mark "Consent on behalf of your company"<br/>
Click Accept <br/>

You will be prompted with another json similar to the one described above<br/>
Copy the token from the response<br/>
The property named "access_token", "expiry": refresh_token": "", "token_type" and paste in the json provided above <br/>
It might be a bit complex because access token and refresh token are very long string
Make sure to copy them properly
The token that you will receive will look like the sample below

```json
{
  "status": 200,
  "data": {
    "id": "",
    "email": "",
    "name": "",
    "given_name": "",
    "family_name": "",
    "access_token": "",
    "refresh_token": "",
    "token": {
      "access_token": "",
      "token_type": "",
      "refresh_token": "",
      "expiry": ""
    }
  }
}
```

Now you have a json filled with all the values CogniX needs.<br/>
Paste it into the connector specific configuration <br/>

**Refresh frequency** is in second it tells CogniX every each seconds it need to rescan the source.
Make it at least 86400 (one day in seconds) <br/>
**connector credentials**
not used add a number