# Member Channels

### A Discord bot to enable users to create their own voice channel

# Why?
Often, you see servers with excessive lists of voice channels labeled:

![Excessive channel lists](static/channellist.png "Excessive channel lists")

This project aims to fix that by basing the voice channels on users and only creating the channels when said user wants to.

# Features

- Automagically create and delete voice channels for users
- The owner of the MemberChannel has permission to rename the channel (and maybe more in future)
- The name of the MemberChannel is save per user per server
- If the owner of the MemberChannel leaves, it is handed off to the user that's been there the longest (FIFO)

# How does it work?
When the bot joins your server, it'll create two default channels. A category channel named `Member Channels` and underneath it a placeholder voice channel called `[ + New ]`

When any user joins the placehold channel, the bot will:
- Create a voice channel named `User's channel` with the users' name
- Gives the user permission to edit that channel
- Moves the user into the channel

Once all users have left a channel the channel will be removed.

![Joining a channel](static/joining.gif "Joining a channel")

## Renaming channels
If any of the channels are renamed, including the category and the placeholder voice channel, the name will be remembered.

For example, say a user named `ImAHuman` renames their channel from the default `IAmHuman's Channel` to `This is my channel`. From then on when `ImAHuman` joins the placeholder, it'll name their channel `This is my channel` instead.

# Building

## Prerequisites

Not yet written.

# Joining the bot to your server

## Your own bot

Not yet written.

## Public bot

Not yet written
