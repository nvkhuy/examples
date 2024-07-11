import json
import base64
from io import BytesIO

import requests
from PIL import Image
from yolov8_onnx import YOLOv8

# Initialize YOLOv8 object detector
yolov8_detector = YOLOv8('best.onnx')


def handler(event, context):
    # get payload
    body = json.loads(event['body'])

    # get params
    image_url = body.get('image')
    size = body.get('size', 640)
    conf_threshold = body.get('confidence', 0.3)
    iou_threshold = body.get('overlap', 0.5)

    print('image_url={} size={} conf_threshold={} iou_threshold={}'.format(image_url, size, conf_threshold,
                                                                           iou_threshold))

    # download image
    response = requests.get(image_url)
    img = Image.open(BytesIO(response.content))

    # convert image to base64
    buffered = BytesIO()
    img.save(buffered, format="JPEG")
    img_b64 = base64.b64encode(buffered.getvalue()).decode('ascii')

    # open image
    img = Image.open(BytesIO(base64.b64decode(img_b64.encode('ascii'))))

    # infer result
    detections = yolov8_detector(img, size=size, conf_thres=conf_threshold, iou_thres=iou_threshold)
    predictions = []
    for detect in detections:
        prediction = {}
        if detect:
            if detect['class_id']:
                detect['class_name'] = classes[detect['class_id']]
                yolobox = bbox2yolobox(detect['bbox'][0], detect['bbox'][1], detect['bbox'][2], detect['bbox'][3])
                prediction['x'] = yolobox[0]
                prediction['y'] = yolobox[1]
                prediction['width'] = yolobox[2]
                prediction['height'] = yolobox[3]
                prediction['confidence'] = detect['score']
                prediction['class'] = classes[detect['class_id']]
                prediction['class_id'] = detect['class_id']
                predictions.append(prediction)
    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "image": {
                "width": img.size[0],
                "height": img.size[1],
            },
            "predictions": predictions
        }),
    }


classes = {
    0: 'men-activewear',
    1: 'men-denim',
    2: 'men-outerwears',
    3: 'men-pants',
    4: 'men-shorts',
    5: 'men-sweaters',
    6: 'men-swimwear',
    7: 'men-tops',
    8: 'women-activewear',
    9: 'women-denim',
    10: 'women-dresses',
    11: 'women-outerwears',
    12: 'women-pants',
    13: 'women-shorts',
    14: 'women-skirts',
    15: 'women-sweaters',
    16: 'women-swimwear',
    17: 'women-tops'
}


def bbox2yolobox(x1, y1, x2, y2):
    w = x2 - x1
    h = y2 - y1
    x = x1 + w / 2
    y = y1 + h / 2
    return x, y, w, h
