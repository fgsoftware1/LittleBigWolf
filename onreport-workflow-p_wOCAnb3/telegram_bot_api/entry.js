import { axios } from "@pipedream/platform";

export default defineComponent({
  props: {
    telegram_bot_api: {
      type: "app",
      app: "telegram_bot_api",
    },
    twitter: {
      type: "app",
      app: "twitter",
    },
    command: {
      type: "string",
      label: "Command",
      description: "The command to listen for in Telegram messages",
    },
  },
  async run({ steps, $ }) {
    const updates = await axios($, {
      url: `https://api.telegram.org/bot${this.telegram_bot_api.$auth.token}/getUpdates`,
    });

    const commandMessage = updates.result.find(
      (update) => update.message.text.startsWith(`/${this.command}`)
    );

    if (commandMessage) {
      const mention = commandMessage.message.text.split(" ")[1];
      if (mention.startsWith("@")) {
        await axios($, {
          method: "POST",
          url: `https://api.twitter.com/1.1/users/report_spam.json`,
          headers: {
            Authorization: `Bearer ${this.twitter.$auth.oauth_access_token}`,
          },
          data: {
            screen_name: mention.slice(1),
          },
        });
      }
    }
  },
});