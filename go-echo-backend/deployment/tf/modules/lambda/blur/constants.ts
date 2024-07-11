import dotenv from "dotenv";

dotenv.config();

const constants = {
  sizes: [
    "48x48",
    "80x80",
    "150x150",
    "200x125",
    "360w",
    "360x270",
    "360h",
    "480h",
    "720h",
    "960h",
    "1080h",
    "1280h",
    "64w",
    "128w",
    "480w",
    "640w",
    "720w",
    "960w",
    "1080w",
    "1280w",
    "2560w",
  ],

  env: {
    destCDNUrl: process.env.AWS_CDN_URL || "",
    storageUrl: process.env.AWS_STORAGE_URL || "",
    originBucket: process.env.AWS_STORAGE_BUCKET || "",
    destBucket: process.env.AWS_CDN_BUCKET || "",
    destRegion: process.env.AWS_REGION || "",
    originRegion: process.env.AWS_REGION || "",
    jwtSecret: process.env.JWT_SECRET || "",
  },
};

console.log(constants);

export default constants;
