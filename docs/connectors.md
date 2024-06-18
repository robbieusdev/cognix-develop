### Supported MIME types 
In list below mapping  supported mime types to proto fileType
```json
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       proto.FileType_XLS,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": proto.FileType_DOC,
	"application/pdf": proto.FileType_PDF,
	"application/rtf": proto.FileType_RTF,
	"text/rtf":        proto.FileType_RTF,
	"text/plain":      proto.FileType_TXT,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": proto.FileType_PPT,
```

### Parameter description 
User should prepare json for each connector and copy it to field ConnectorSpecificConfig on page create/edit connector.

#### Web

```json
{
  "url": "",
  "site_map": "",
  "search_for_sitemap": false,
  "url_recursive": false
}
```

#### File 
```json
{
  "file_name": "",
  "mime_type": "" 
}
```
mime type should be determined when user upload file to minio. 

Connector translate this parameter for a chunker service in next format 
```json
    "minio:<bucketname>:<filename>"
```

#### OneDrive 

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

- folder : optional, folder name for scanning
- recursive :  false - scan only given folder , true - scan nested folders
- token : OAuth token for access to ```one drive```

```NOTE```
While UI is not ready 

You can get access using requests 
```js
http://localhost:8080/api/oauth/microsoft/auth_url?redirect_url=http://localhost:8080
```
then go to the link from response 
```json

{
  "status": 200,
  "data": "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=<id>>&scope=offline_access Files.Read.All Sites.ReadWrite.All&response_type=code&redirect_uri=http://localhost:8080/api/oauth/microsoft/callback"
}
```
Sigin using microsoft account and grant permission 

copy token from response 

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
