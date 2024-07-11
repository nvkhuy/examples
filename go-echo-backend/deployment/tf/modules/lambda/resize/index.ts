import { Context, APIGatewayProxyResult, APIGatewayEvent } from "aws-lambda";
import Helper from "./helper";
import Constants from "./constants";
import constants from "./constants";

export const handler = async (
  event: APIGatewayEvent,
  context: Context
): Promise<APIGatewayProxyResult> => {
  console.log(`Event: ${JSON.stringify(event, null, 2)}`);
  console.log(`Context: ${JSON.stringify(context, null, 2)}`);

  const token = event.queryStringParameters?.token;
  const key = event.queryStringParameters?.key;
  const size = event.queryStringParameters?.size;
  const noCache = event.queryStringParameters?.no_cache;
  const checkExists = event.queryStringParameters?.check_exists;
  const sizeKey = `${size}/${key}`;

  console.log("Process", { token, key, size, sizeKey });
  if (!key) {
    return {
      statusCode: 403,
      body: JSON.stringify("Invalid key."),
    };
  }

  const valid = await Helper.isValidJwtToken(token, sizeKey);
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

  const sizeParams = Helper.parseFileSize(size);
  if (!sizeParams) {
    console.log("Invalid size", size)
    return {
      statusCode: 400,
      body: JSON.stringify("Invalid file."),
    };
  }

  if (checkExists && !noCache) {
    try {
      const headObject = await Helper.headObject(constants.env.destBucket, sizeKey)
      if (headObject?.WebsiteRedirectLocation) {
        const url = Constants.env.destCDNUrl + '/' + sizeKey;
        return {
          statusCode: 301,
          headers: {
            Location: url,
          },
          body: "",
        };
      }
    } catch (error) {

    }

  }


  try {
    const data = await Helper.getObject(constants.env.originBucket, key);
    if (!data) {
      return {
        statusCode: 400,
        body: JSON.stringify("Get object failed"),
      };
    }
    console.log("Get object success", data?.length, sizeParams);

    const imageData = await Helper.resizeImage(data, sizeParams);
    if (!imageData) {
      return {
        statusCode: 400,
        body: JSON.stringify("Resize image failed"),
      };
    }
    console.log("Resize image success", data?.length);

    await Helper.putObject(constants.env.destBucket, sizeKey, imageData)

    const url = Constants.env.destCDNUrl + '/' + sizeKey;
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
