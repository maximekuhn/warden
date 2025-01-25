# migrations
This directory contains all SQL files to apply to the database to make it ready for `warden`.  

> [!IMPORTANT]  
> If `warden` is already running in production, changing any existing file won't have any effect
> on the database, as the server will only apply required migrations and skip old ones.
