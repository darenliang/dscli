# Dscli Quickstart

### Step 1: Enable developer mode

Go to the Discord client or https://discord.com/app.

Click on the gear icon.

![](images/1.1.png)

Go to "Appearance" and toggle "Developer Mode" to on.

![](images/1.2.png)

### Step 2: Create a Discord server

Go to the Discord client or https://discord.com/app.

Click on the plus icon.

![](images/2.1.png)

Click on "Create My Own".

![](images/2.2.png)

Enter anything as the server name and click on "Create".

![](images/2.3.png)

Right-click on the server icon and click on "Copy ID".

![](images/2.4.png)

**Save the server ID.**

### Step 3: Create a Discord bot

Go to https://discord.com/developers/applications.

Click on "New Application".

![](images/3.1.png)

Enter anything as the Discord bot name and click on "Create".

![](images/3.2.png)

Go to "Bot" and click on "Add Bot".

![](images/3.3.png)

Click on "Yes, do it!".

![](images/3.4.png)

Click on "Copy".

![](images/3.5.png)

**Save the Discord bot token.**

### Step 4: Invite the Discord bot to the server

Go to https://discord.com/developers/applications.

Select the application, go to "OAuth2" and click on "bot" under "SCOPES".

![](images/4.1.png)

Under "SCOPES" and in "BOT PERMISSIONS", click on "Administrator".

Click on "Copy".

![](images/4.2.png)

Go to the link copied and invite the bot.

![](images/4.3.png)

### Step 5: Configure dscli

**As part of the setup, dscli will delete any remaining channels in the server.**

```
dscli -t=<YOUR-DISCORD-BOT-TOKEN> -i=<YOUR-SERVER-ID> -d
```
