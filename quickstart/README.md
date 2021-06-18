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

### Step 3: Get your own user token

> Automating user accounts is technically against TOS, use at your own risk.

Press **Ctrl+Shift+I** (⌘⌥I on Mac) on Discord to show developer tools.

Navigate to the **Application** tab.

Select **Local Storage** > **https://discordapp.com** on the left.

Press **Ctrl+R** (⌘R) to reload.

Find **token** at the bottom and copy the value.

> You may also get your own bot token instead.

![](images/3.1.gif)

### Step 4: Configure dscli

**As part of the setup, dscli will delete any remaining channels in the server.**

> If you are using a bot token, add the "Bot " prefix to the discord token.

```
dscli config -t=<YOUR-DISCORD-TOKEN> -i=<YOUR-SERVER-ID> -d
```

For configuring dscli interactively:

```
dscli config
```
