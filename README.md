Discord Stats Bot (Go + discordgo)

What it does
- Posts (and keeps updated) a single message in a predefined text channel, showing:
  - Number of currently connected users in voice across the guild
  - Total members in the guild

Quick start
1) Create a Discord Application + Bot in the Developer Portal and invite it to your server.
2) Enable Gateway Intents for the bot:
   - Guilds
   - GuildMembers (needed for total member count)
   - GuildVoiceStates (needed for connected-in-voice count)
3) Gather IDs:
   - DISCORD_TOKEN: your bot token (without the "Bot " prefix; the app adds it)
   - GUILD_ID: the server ID
   - CHANNEL_ID: the text channel ID where the stats message should live
   - MESSAGE_ID (optional): if you already have a message to update
4) Run
   - go mod tidy
   - go run ./cmd/bot

Configuration via environment variables
- DISCORD_TOKEN
- GUILD_ID
- CHANNEL_ID
- MESSAGE_ID (optional)

Notes
- On the first run (if MESSAGE_ID is not provided), the bot will create the stats message and print its ID in the logs. Persist that ID by setting MESSAGE_ID on subsequent runs to keep updating the same message.
- If you intended "connected" to mean "online presence" (online/idle/dnd), that requires the GuildPresences intent and a different counting method. Ask and I can adapt the code.
