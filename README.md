# Discord-Attendance
Internal tool that allows users to clock in and out of the workplace using Discord

## How to setup
Create a tokens.env file in the root directory and add the following:
- BOT_TOKEN
- APP_ID
- AUTHORIZED
Where authorized is a comma separated list of administrative users. Alternatively, you can add these to your environment variables.

To run the bot, use the following command:
```go run src/main.go```

### Creating a [Discord Bot](https://discordapp.com/developers/applications)
1. Give the app a friendly name and click the **Create App** button
2. Put **APP ID** in your tokens.env file
3. Scroll down to the **Bot** section
4. Click the **Create a Bot User** button
6. Put **TOKEN** in your tokens.env file
7. Generate a link to add the bot to your server by clicking the **Generate OAuth2 URL** button
8. Select the **bot** scope and the **Administrator** permission, as well as the **applications.commands** permission
9. Copy the generated link and paste it into your browser
10. Select the server you want to add the bot to and click **Authorize**

## How to use
Right click -> Apps to clock-in, clock-out, and new-period. The following slash commands are also available:

- /clockin
- /clockout
- /addhours [hours] [user]
- /removehours [hours] [user]
- /export [user]

Unauthorized users can only clock-in and clock-out. Authorized users can use all commands. /addhours and /remove hours are used to manually add or remove hours from an user.

At the end of the month, the administrator can use the /export command to export a user's hours to a CSV file. This file can then be viewed to calculate each user's salary. The administrator can then use new-period to reset the hours for the next month.