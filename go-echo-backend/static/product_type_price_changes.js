function doGet(e) {
    let sheetName = "RWD: ITEM + CMPT Price";
    let readSheetFrom = "A3";
    let readSheetTo = "AA";

    // Get the edited range and its sheet
    let editedRange = e.range;
    let editedSheetName = editedRange.getSheet().getName();

    let responseData = {
        message: "Empty"
    };

    if (editedSheetName === sheetName) {
        // Convert the edited range into a Range object
        let editedRangeAsRange = editedRange;

        // Convert the specified read range into a Range object
        let readRange = e.source.getSheetByName(sheetName).getRange(readSheetFrom + ":" + readSheetTo);

        // Check if the edited range is within the specified read range
        if (isInRange(editedRangeAsRange, readRange)) {
            // Perform your action if the edited range is within the specified range
            let responseData = {
                message: "Edited range is within the specified range."
            };
            patchingProd();
            patchingDev();
        }
    } else {
        // Perform your action if the edited sheet is not "RWD ITEM + CMPT Price"
        let responseData = {
            message: "Edited sheet is not 'RWD ITEM + CMPT Price'."
        };
    }

    // Convert the response data to JSON
    let jsonString = JSON.stringify(responseData);

    // Set the content type to JSON
    let output = ContentService.createTextOutput(jsonString);
    output.setMimeType(ContentService.MimeType.JSON);

    return output;
}

// Function to check if a given range is within another range
function isInRange(testRange, referenceRange) {
    return (
        testRange.getRow() >= referenceRange.getRow() &&
        testRange.getLastRow() <= referenceRange.getLastRow() &&
        testRange.getColumn() >= referenceRange.getColumn() &&
        testRange.getLastColumn() <= referenceRange.getLastColumn()
    );
}

function patchingDev() {
    let url = 'https://dev-api.joininflow.io/api/v1/admin/products/types/price';
    let headers = {
        'authority': 'dev-api.joininflow.io',
        'accept': 'application/json, text/plain, */*',
        'accept-language': 'en-US,en;q=0.9',
        'authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwidHoiOiJBc2lhL1NhaWdvbiIsImNpZCI6IiIsImN0eXBlIjoiYnV5ZXIiLCJhdWQiOiJzdXBlcl9hZG1pbiIsImlzcyI6ImNnNWdwM3VqbnZ2YjBrbnAzOTAwIiwic3ViIjoic3VwZXJfYWRtaW4ifQ.IOez2kQOF9nlBGzbard-eVT6yrjT7Ex7OgtTj8sLmdE',
        'build_date': '13:11:48 13/10/2023',
        'content-type': 'application/json',
        'origin': 'https://dev-admin.joininflow.io',
        'referer': 'https://dev-admin.joininflow.io/',
        'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
        'sec-ch-ua-mobile': '?0',
        'sec-ch-ua-platform': '"macOS"',
        'sec-fetch-dest': 'empty',
        'sec-fetch-mode': 'cors',
        'sec-fetch-site': 'same-site',
        'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36',
    };

    let payload = {
        "spreadsheet_id": "1AnLk432WNOvN5hBVJxgPW8dO-gCda02JY6-c3xdHKp8",
        "sheet_name": "RWD: ITEM + CMPT Price",
        "from": "A3",
        "to": "AA"
    };

    let options = {
        'method': 'patch',
        'headers': headers,
        'payload': JSON.stringify(payload)
    };

    let response = UrlFetchApp.fetch(url, options);
    let content = response.getContentText();

    return ContentService.createTextOutput(content).setMimeType(ContentService.MimeType.JSON);
}

function patchingProd() {
    let url = 'https://api.joininflow.io/api/v1/admin/products/types/price';
    let headers = {
        'authority': 'api.joininflow.io',
        'accept': 'application/json, text/plain, */*',
        'accept-language': 'en-US,en;q=0.9',
        'authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNqY3ZtaDRwcmY4dDA5OHJiZnYwIiwidHoiOiJBc2lhL1NhaWdvbiIsImNpZCI6IiIsImN0eXBlIjoiYnV5ZXIiLCJhdWQiOiJzdXBlcl9hZG1pbiIsImlzcyI6ImNqY3Z0dDRwcmY4dDA5OHJiZnZnIiwic3ViIjoic3VwZXJfYWRtaW4ifQ.8E2Teoqy5PgFlyh5MllYBZD3QC6OIKNr4mc1Xt6ulYY',
        'build_date': '13:11:48 13/10/2023',
        'content-type': 'application/json',
        'origin': 'https://admin.joininflow.io',
        'referer': 'https://admin.joininflow.io/',
        'sec-ch-ua': '"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"',
        'sec-ch-ua-mobile': '?0',
        'sec-ch-ua-platform': '"macOS"',
        'sec-fetch-dest': 'empty',
        'sec-fetch-mode': 'cors',
        'sec-fetch-site': 'same-site',
        'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36',
    };

    let payload = {
        "spreadsheet_id": "1AnLk432WNOvN5hBVJxgPW8dO-gCda02JY6-c3xdHKp8",
        "sheet_name": "RWD: ITEM + CMPT Price",
        "from": "A3",
        "to": "AA"
    };

    let options = {
        'method': 'patch',
        'headers': headers,
        'payload': JSON.stringify(payload)
    };

    let response = UrlFetchApp.fetch(url, options);
    let content = response.getContentText();

    return ContentService.createTextOutput(content).setMimeType(ContentService.MimeType.JSON);
}