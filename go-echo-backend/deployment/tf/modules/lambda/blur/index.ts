import { Context, APIGatewayProxyResult, APIGatewayEvent } from "aws-lambda";
import Helper from "./helper";
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
  const sizeKey = `blur/${size}/${key}`;

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
    return {
      statusCode: 403,
      body: JSON.stringify("Invalid size."),
    };
  }

  const sizeParams = Helper.parseFileSize(size);
  if (!sizeParams) {
    return {
      statusCode: 400,
      body: JSON.stringify("Invalid size."),
    };
  }

  // if (checkExists && !noCache) {
  //   const metadata = await Helper.getObjectMetadata(constants.env.originBucket, key)
  //   if (typeof metadata?.blurhash === 'string' && metadata?.blurhash !== '' && typeof metadata?.blurhash_data === 'string' && metadata?.blurhash_data !== '') {
  //     return {
  //       statusCode: 200,
  //       body: JSON.stringify({
  //         blurhash: metadata?.blurhash,
  //         blurhash_data_url: metadata?.blurhash_data,
  //         blurhash_avg: metadata?.blurhash_avg,
  //         blurhash_thumbnail: metadata?.blurhash_thumbnail,
  //         blurhash_image_width: metadata?.blurhash_image_width,
  //         blurhash_image_height: metadata?.blurhash_image_height,
  //       }),
  //     };
  //   }
  // }


  try {
    const data = await Helper.getObject(constants.env.originBucket, key);
    if (!data) {
      return {
        statusCode: 400,
        body: JSON.stringify("Get object failed"),
      };
    }
    console.log("Get object success", data?.length);

    const { encoded, blurhashAvgStr, info } = await Helper.encodeBlurhash(data, sizeParams);
    if (!encoded) {
      return {
        statusCode: 400,
        body: JSON.stringify("Transform image failed"),
      };
    }

    const blurData = await Helper.convertToBlurhashData(encoded, 32);

    const metadata = {
      blurhash: encoded,
      blurhash_data: blurData,
      blurhash_avg: blurhashAvgStr,
      blurhash_thumbnail: size,
      blurhash_image_width: `${info?.width}`,
      blurhash_image_height: `${info?.height}`,
    }

    // await Helper.updateObjectMetadata(constants.env.originBucket, key, metadata)


    return {
      statusCode: 200,
      body: JSON.stringify(metadata),
    };
  } catch (err: any) {
    console.log("err", err);
    return {
      statusCode: 500,
      body: err?.message,
    };
  }
};
