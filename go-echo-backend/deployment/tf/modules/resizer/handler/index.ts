import { Context, APIGatewayProxyResult, APIGatewayEvent } from "aws-lambda";
import Helper from "./helper";
import Constants from "./constants";

export const handler = async (
  event: APIGatewayEvent,
  context: Context
): Promise<APIGatewayProxyResult> => {
  console.log(`Event: ${JSON.stringify(event, null, 2)}`);
  console.log(`Context: ${JSON.stringify(context, null, 2)}`);
  const token = event.queryStringParameters?.token;
  const key = event.queryStringParameters?.key;
  const size = event.queryStringParameters?.size;

  if (!key) {
    return {
      statusCode: 403,
      body: JSON.stringify("Invalid key."),
    };
  }

  const valid = await Helper.isValidJwtToken(token, `${size}/${key}`);
  if (!valid) {
    return {
      statusCode: 403,
      body: JSON.stringify("Invalid token."),
    };
  }

  if (!size) {
    const url = Constants.env.storageUrl + "/" + key;
    return {
      statusCode: 301,
      headers: {
        Location: url,
      },
      body: "",
    };
  }

  const params = Helper.parseFileSize(size);
  if (!params) {
    return {
      statusCode: 400,
      body: JSON.stringify("Invalid file."),
    };
  }

  try {
    const readStream = Helper.getS3Stream(key);
    if (!readStream) {
      return {
        statusCode: 400,
        body: JSON.stringify("Get stream failed"),
      };
    }
    const resizeStream = Helper.stream2Sharp(params);
    const putSteamResp = Helper.putS3Stream(key);
    if (!putSteamResp) {
      return {
        statusCode: 400,
        body: JSON.stringify("Put stream failed"),
      };
    }

    // trigger stream
    readStream.pipe(resizeStream).pipe(putSteamResp?.writeStream);

    // wait for the stream
    await putSteamResp?.success;

    const url = Constants.env.destCDNUrl + "/" + key;

    // 301 redirect to new image
    return {
      statusCode: 301,
      headers: {
        Location: url,
      },
      body: "",
    };
  } catch (err: any) {
    console.log("err", err);
    return {
      statusCode: 500,
      body: err?.message,
    };
  }
};
