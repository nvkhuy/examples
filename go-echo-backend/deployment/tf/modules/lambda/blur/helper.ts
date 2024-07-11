import {
  S3Client,
  PutObjectCommand,
  PutObjectCommandInput,
  GetObjectCommandOutput,
  GetObjectCommand,
  GetObjectCommandInput,
  HeadObjectCommand,
  CopyObjectCommand,
  CopyObjectCommandOutput
} from "@aws-sdk/client-s3";
import sharp from "sharp";
import mime from "mime-types";
import jwt from "jsonwebtoken";
import Constants from "./constants";
import constants from "./constants";
import { Readable } from "stream";
import { encode, decode } from "blurhash";
import { getAverageColor } from "./blurhash-utils";

class Helper {
  static s3OriginClient: S3Client = new S3Client({
    region: constants.env.originRegion,
  });
  static s3DestClient: S3Client = new S3Client({
    region: constants.env.destRegion,
  });

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
            console.log(err);
            resolve(false);
            return;
          }

          resolve(true);
        }
      );
    });
  }

  static async streamToBuffer(response: GetObjectCommandOutput) {
    const stream = response.Body as Readable;
    const chunks: Buffer[] = [];
    return new Promise<Buffer>((resolve, reject) => {
      stream.on("data", (chunk) => chunks.push(chunk));
      stream.on("error", (err) => reject(err));
      stream.on("end", () => resolve(Buffer.concat(chunks)));
    });
  }

  static async getObject(bucket: string, key: string) {
    const input: GetObjectCommandInput = {
      Bucket: bucket,
      Key: key,
    };

    const command: GetObjectCommand = new GetObjectCommand(input);

    try {
      const response: GetObjectCommandOutput = await Helper.s3OriginClient.send(
        command
      );

      return Helper.streamToBuffer(response)
    } catch (error) {
      throw error;
    }
  }

  static async putObject(
    bucket: string,
    key: string,
    image: Buffer,
    ct?: string
  ) {
    const contentType = ct || mime.contentType(key);
    const input: PutObjectCommandInput = {
      Bucket: bucket,
      Key: key,
      Body: image,
      ACL: "private",
    };
    if (contentType) {
      input["ContentType"] = contentType;
    }

    const command: PutObjectCommand = new PutObjectCommand(input);

    try {
      await Helper.s3DestClient.send(command);
      return key;
    } catch (error) {
      throw error;
    }
  }

  static async encodeBlurhash(image: Buffer, size: any) {
    try {
      const { data, info } = await sharp(image)
        .resize(
          Object.assign(size, {
            withoutEnlargement: true,
          })
        )
        .ensureAlpha()
        .raw()
        .toBuffer({
          resolveWithObject: true,
        });

      const isPortrait = info.height > info.width
      const blurhashWidth = isPortrait ? 3 : 4
      const blurhashHeight = isPortrait ? 4 : 3

      const encoded = encode(
        new Uint8ClampedArray(data),
        info.width,
        info.height,
        blurhashWidth,
        blurhashHeight
      );


      const blurhashAvg = getAverageColor(encoded)
      const blurhashAvgStr = blurhashAvg ? `${blurhashAvg?.r},${blurhashAvg?.g},${blurhashAvg?.b}` : ''


      return { encoded, blurhashAvg, blurhashAvgStr, info }
    } catch (error) {
      throw error;
    }
  }

  static async convertToBlurhashData(blurhashEncoded: string, size = 32) {
    try {
      const pixels = decode(blurhashEncoded, size, size)
      // Convert the pixel data to a PNG image using sharp
      const pngBuffer = await sharp(Buffer.from(pixels), {
        raw: {
          width: size,
          height: size,
          channels: 4
        }
      })
        .png()
        .toBuffer()

      const blurDataURL = `data:image/png;base64,${pngBuffer.toString('base64')}`
      return blurDataURL
    } catch (error) {
      console.log("**** convertToBlurhashData error", error)
      throw error;
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


  static async getObjectMetadata(
    bucket: string,
    key: string
  ): Promise<Record<string, string> | undefined> {
    try {
      const headObjectCmdInput = new HeadObjectCommand({
        Bucket: bucket,
        Key: key,
      })

      const headObjectCmd = await Helper.s3OriginClient.send(headObjectCmdInput)

      return headObjectCmd.Metadata
    } catch (error) {
      return undefined
    }


  }


  static async updateObjectMetadata(
    bucketName: string,
    objectKey: string,
    metadata: Record<string, string>
  ): Promise<CopyObjectCommandOutput> {
    try {

      const existingMetadata = await Helper.getObjectMetadata(bucketName, objectKey)

      // Update the metadata by adding/modifying key-value pairs
      const updatedMetadata = { ...(existingMetadata || {}), ...(metadata || {}) };

      const copyCmdInput = new CopyObjectCommand({
        Bucket: bucketName,
        CopySource: `${bucketName}/${objectKey}`,
        Key: objectKey,
        Metadata: updatedMetadata || {},
        MetadataDirective: "REPLACE",
      })
      // Update the object's metadata
      const copyCmd = await Helper.s3DestClient.send(copyCmdInput)

      return copyCmd;
    } catch (error) {
      console.error(`Error updating object metadata: ${error}`);
      throw error;
    }
  }


}

export default Helper;
