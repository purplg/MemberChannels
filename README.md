# Dynamic Members

### A Discord bot to dynamically create voice channels for users 

# Why?
Often, you see servers with long lists of voice channels labeled:

* Video Game
    * Video Game 1
    * Video Game 2
    * Video Game 3
    * Video Game 4
* Other Game
    * Other Game 1
    * Other Game 2
    * Other Game 3
    * Other Game 4

This project aims to fix that by basing the voice channels on users and only creating the channels when said user wants to.

# How does it work?
When the bot joins your server, it'll create two default channels. A category channel named `Dynamic Channels` and underneat it a placehold voice channel called `[+] Create channel`

When any user joins the placehold channel, the bot:
- Creates a voice channel named `User's channel` with the users' name
- Gives the uesr permission to manage that channel
- Moves the user into the channel

Once all users have left a channel and it has become empty, it'll be deleted until the creator rejoins the placeholder channel.

## Renaming
If any of the channels are renamed, including the category and the placehold voice channel, the name will be remembered.

For example, if a user named `ImAHuman` renames their channel from the default `IAmHuman's Channel` to `This is my channel"`, then from then on when `ImAHuman` joins the placeholder, it'll name their channel `This is my channel`
