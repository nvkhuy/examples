const platforms = [
  { id: "dej340x3ljg03", name: "Inflow Brand" },
  { id: "d2pby2jugh76ta", name: "Inflow Seller" },
  { id: "d2nq6kkfvoa84i", name: "Inflow Website Legacy" },
  { id: "dicl44iirx47s", name: "Inflow Website" },
  { id: "d1vvfs0iyron7b", name: "Inflow Admin" },
];

function getPlatform(msg) {
  return platforms?.find((item) => msg.includes(item?.id));
}

function getEnvPrefix(msg) {
  if (msg?.includes("/dev")) {
    return "[Dev] ";
  }

  if (msg?.includes("/beta")) {
    return "[Beta] ";
  }

  if (msg?.includes("/prod")) {
    return "[Prod] ";
  }

  return "";
}
export const handler = async (event) => {
  const webhookURL = process.env.SLACK_WEBHOOK_URL;

  console.log(`Webhook URL: ${webhookURL}`);
  console.log(`Event: ${JSON.stringify(event, null, 2)}`);

  let sns = event.Records[0].Sns.Message || "";
  sns = sns.replace(/['"]+/g, "");

  const platform = getPlatform(sns);
  const envPrefix = getEnvPrefix(sns);

  let color = "";
  let title = "";
  if (sns.includes("build status is FAILED")) {
    color = "#E52E59";
    title = `${envPrefix}Release to ${platform?.name} failed`;
  } else if (sns.includes("build status is SUCCEED")) {
    color = "#21E27C";
    title = `${envPrefix}Release to ${platform?.name} succeeded.`;
  } else if (sns.includes(`build status is STARTED`)) {
    color = "#3788DD";
    title = `${envPrefix}Release to ${platform?.name} has started...`;
  }

  const data = JSON.stringify({
    attachments: [
      {
        mrkdwn_in: ["text"],
        fallback: sns,
        color,
        title,
        text: sns,
      },
    ],
  });

  const rawResponse = await fetch(webhookURL, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: data,
  });

  const content = await rawResponse.json();

  return content;
};
