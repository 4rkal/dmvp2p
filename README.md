# dmvp2p

DMVP2P: Donate Monero Via P2Pool

DMVP2P allows you to donate monero to your favorite projects/organisations/individuals using your cpu and the magic of [P2Pool](https://p2pool.io)

## Installing/Downloads

Builds for Linux and Windows are available in the [releases page](https://github.com/4rkal/dmvp2p/releases)

You can also build the project from source by follow these steps:

(You will have to have [go](https://go.dev) installed)

`git clone https://github.com/4rkal/dmvp2p`

`cd dmvp2p`

`

## How can I add myself to get donations?

Edit [users.json](https://github.com/4rkal/dmvp2p/blob/main/helpers/users.json) fill in your details, and create a PR. (you can also contact me and I'll add you [contact](https://4rkal.com/about/#get-in-touch))

Here is an example of a user:
```
    {
    "name": "4rkal",
    "github": "4rkal",
    "x": "4rkal_",
    "website": "https://4rkal.com",
    "address": "425evPnXucaJcAvgYbcasrYJ5k1qUEB7gAB4DSS7nM73hcxMgf3L9fzQCRA441tZbAcSXhRR4DDxT5B3oxBKbnns9A5Z4mi",
    "description": "Programmer, blogger, crypto user."
    }
```

**NOTE**: because of the way that P2Pool works the address **must** be a primary address! (starting in 4 and not 8)