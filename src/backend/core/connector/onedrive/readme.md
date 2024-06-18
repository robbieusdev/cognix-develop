## Register an application
You can register an application using the Azure Active Directory admin center, or by using the Microsoft Graph PowerShell SDK.

### Azure Active Directory admin center
Open a browser and navigate to the Azure Active Directory admin center and login using a personal account (aka: Microsoft Account) or Work or School Account.

Select Azure Active Directory in the left-hand navigation, then select App registrations under Manage.

Select New registration. Enter a name for your application, for example, Go Graph Tutorial.

Set Supported account types as desired. The options are:

Option	Who can sign in?
```azure
Accounts in this organizational directory only	Only users in your Microsoft 365 organization
Accounts in any organizational directory	Users in any Microsoft 365 organization (work or school accounts)
Accounts in any organizational directory ... and personal Microsoft accounts	Users in any Microsoft 365 organization (work or school accounts) and personal Microsoft accounts
```
Leave Redirect URI empty.

Select Register. On the application's Overview page, copy the value of the Application (client) ID and save it, you will need it in the next step. If you chose Accounts in this organizational directory only for Supported account types, also copy the Directory (tenant) ID and save it.

Select Authentication under Manage.
Locate the Advanced settings section and change the
Allow public client flows toggle to Yes,

then choose Save.

### Create Azure application 

- Go to the Azure App registrations page.  https://aka.ms/AppRegistrations/?referrer=https%3A%2F%2Fdev.onedrive.com 
- When prompted, sign in with your account credentials.
- Find My applications and click Add an app.
- Enter your app's name and click Create application.


#### open Permission  page 
 - Click Application registration 
   - Api permission page 
       - Files.Read.All
       - offline_access
   - Authentication 
       - Add platform 
         -  choose web and define redirect URL https://rag.cognix.ch/api/oauth/microsoft/callback