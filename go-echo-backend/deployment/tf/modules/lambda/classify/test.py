import json

# Assuming the handler function is in the same file
from lambda_function import handler


def test_handler_with_url(image_url, size=640, conf_thres=0.3, iou_thres=0.1):
    # Prepare the event
    event = {
        'body': json.dumps({
            'image': image_url,
            'size': size,
            'confidence': conf_thres,
            'overlap': iou_thres,
        })
    }

    # Call the handler
    result = handler(event, None)
    print(result)
    return result


# Test the function
image_url = 'https://lucky.a.bigcontent.io/v1/static/DT-BOOT-ALLSIZES'
test_handler_with_url(image_url)
