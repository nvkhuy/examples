import AWS from "aws-sdk";
import stream from "stream";
import sharp, { ResizeOptions } from "sharp";
import mime from "mime-types";
import jwt from "jsonwebtoken";
import Constants from "./constants";
import { PromiseResult } from "aws-sdk/lib/request";

const originS3 = new AWS.S3({
  region: Constants.env.originRegion,
});

const destS3 = new AWS.S3({
  region: Constants.env.destRegion,
});

class Helper {
  static async isValidJwtToken(accessToken?: string, audience?: string) {
    if (!accessToken) {
      return false;
    }

    return new Promise((resolve, reject) => {
      jwt.verify(
        accessToken,
        Constants.env.jwtSecret as string,
        {
          algorithms: ["HS256"],
          audience: audience,
        },
        (err: any, user: any) => {
          if (err) {
            console.log("**** audience", audience);
            console.log(err);
            resolve(false);
            return;
          }

          resolve(true);
        }
      );
    });
  }

  static getS3Stream(key: string) {
    if (!Constants.env.originBucket) {
      console.log("Origin bucket is not valid");
      return null;
    }

    return originS3
      .getObject({
        Bucket: Constants.env.originBucket,
        Key: key,
      })
      .createReadStream();
  }

  static putS3Stream(key: string, contentType?: string) {
    if (!Constants.env.destBucket) {
      console.log("Dest bucket is not valid");
      return null;
    }

    const pass = new stream.PassThrough();
    const type = contentType || mime.lookup(key);
    return {
      writeStream: pass,
      success: destS3
        .upload({
          Body: pass,
          Bucket: Constants.env.destBucket,
          Key: key,
          ContentType: type as string,
          ACL: "private",
        })
        .promise(),
    };
  }

  static stream2Sharp(params: ResizeOptions) {
    return sharp()
      .resize(
        Object.assign(params, {
          withoutEnlargement: true,
        })
      )
      .withMetadata();
  }

  static async headFile(key: string) {
    if (!Constants.env.destBucket) {
      return null;
    }

    // check if target key already exists
    let target: PromiseResult<AWS.S3.HeadObjectOutput, AWS.AWSError> | null =
      null;
    let contentType: string | undefined = undefined;
    try {
      await destS3
        .headObject({
          Bucket: Constants.env.destBucket,
          Key: key,
        })
        .promise()
        .then((res) => {
          target = res;
          contentType = res.ContentType;
          return { target, contentType };
        })
        .catch(() => {
          console.log(
            "File %s doesn't exist in %s",
            key,
            Constants.env.destBucket
          );
          return { target, contentType };
        });
    } catch (error) {
      return { target, contentType };
    }
  }

  static parseFileSize(size: string) {
    if (Constants.sizes.indexOf(size) === -1) {
      return null;
    }

    let params: any = {};

    // process size from given string
    if (size.slice(-1) == "w") {
      // extract width only
      params.width = parseInt(size.slice(0, -1), 10);
    } else if (size.slice(-1) == "h") {
      // extract height only
      params.height = parseInt(size.slice(0, -1), 10);
    } else {
      // extract width & height
      var size_components = size.split("x");

      // if there aren't 2 values, stop here
      if (size_components.length != 2) return null;

      params = {
        width: parseInt(size_components[0], 10),
        height: parseInt(size_components[1], 10),
      };

      if (isNaN(params.width) || isNaN(params.height)) return null;
    }

    return params;
  }
}

export default Helper;
