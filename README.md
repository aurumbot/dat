# Data directory.

This is one of the key elements to a working bot. As the bot has to
handle a fair bit of data, we went all in and made a universal data
directory. All information the bot will use should be put in here.
The bot will also generate some of its own data from commands and so
on, this will also go in here. However the files that an average
person who just wants a plug-n-play robot should be concerned about
are as follows. Please be aware that the names **must** be kept
unless you have gone through and changed the hard-coded name values
for some reason.

- preferences.json : Contains all the hard coded static data about the
                     bot such as authentication and universal
                     moderation configuration.
