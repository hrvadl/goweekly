# Go weekly 🤓

## Quick start 🚀

### Run app 😬

Make sure you have docker & docker compose tools installed. Then, copy contents of `.env.example` to `.env` and populate it with values. Finally, in the root of the directory run:

```bash
docker compose up
```

## Architecture

App consists of 4 microservices:

1. Crawler - is responsible for obtaining articles and publishing them to the RabbitMQ. Should be run as a cron/scheduled job (to be done).
2. Core - core service, responsible for listening to RabbitMQ queue and calling other services.
3. Translator - service, used for translation. Uses lingva under the hood.
4. Sender - service, used to distribute translated messages to Telegram channel.

## Local software development 👷🏻

### Code quality 💅🏻

Make sure you have [pre-commit](https://pre-commit.com/) installed and added to your `$PATH`. Then run following command from the root of the repo:

```bash
pre-commit install
```

Now, you must install [golangci-lint](https://golangci-lint.run/) and include it to your path as well.

Congrats, now pre-commit will run before your commits.
Also, in case you're using VSCode linting on save should be integrated into your editor.

## What's the purpose? 🥸

I **like** _Go_ and I **like** reading tech articles.
I appreciate work from admins of [go-weekly website](https://golangweekly.com/).
But, _I don't like_ to be notified about fresh weekly digest by mail.
I'd like receive weekly digest directly in the telegram, because ... Because, why not?
Also, I'd like to contribute to the Ukrainian golang community and make summaries translated to ukrainian language. 🇺🇦

## Contributors 🏋🏻

- [Alex 🤓](https://github.com/oleksandrcherevkov)
- [Vadym 💅🏻](https://github.com/hrvadl)
