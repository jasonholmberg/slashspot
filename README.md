# /Spot
This is a slash command based bot for Slack to help increase the visibility of available parking spots in a common location, like a parking gargare where many spot are reserved but unused. A Slack user could leverage `/spot` find an available spot and claim it. A holder of assigned parking spot can also register a spot's availability with `/spot` so that other may use it in the holder's absence. 

## Usage:
`/spot help` will return the help text

`/spot [find or open]` will return a list of spots available today

`/spot [claim or take or reserve] <spot-id>` will take/reserve a spot or tell you if it is taken

`/spot [reg or register or set] <spot-id> [date]` will make a spot available for use for the day. If a data is given, the spot will be made available for that date.

`/spot drop [ <spot-id> | all ]` will drop the registration of a particular spot or `all` will drop all your registrations.  You must have registered a spot to drop its registration.

## How it works

- `/spot` only knows about **registered** spots on any given day and perhaps future date.  It doesn't not know about all spots in a parking garage and their relative status. 

- A spot's availability is dependent on the holder of the spot registering its avialability for the current day or for dates in the future.

- Spots can only be claimed on the current day. You cannot claim a spot for tomorrow, for example.

- It is be possible for a spot to be claimed and then re-registered on the same day.

- `/spot` is **not** smart enough to guard against fraudulant spot registrations, so please play nice and don't make fraudualant registrations.

- `/spot` will track who registers a particular spot and on what date.

- `/spot` will eventually discard all spot registrations set on dates in the past. 

- `/spot` only stores a files-based list of available spots and their respoective dates.  It does not store claimed spots.

## Setting up /Spot

Ensure you have the necessary properties in the `.evn` install next to the compiled artifact.  Specifically, you need to find these in Slack after creating a new Slash App in Slack:

```
export SPOT_SLACK_SIGNING_SECRET=[YOUR_SIGNING_SECRET_HERE]
export SPOT_SLACK_VERIFICATION_TOKEN=[YOUR_VERIFICATION_TOKEN]
```

Deploy this some place after compiling it for the approriate platform. And point you Slack App to the correct location.

