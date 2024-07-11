function doGet(e) {
    let sheetENName = "EN";
    let sheetVIName = "VI";
    let readSheetFrom = "A1";
    let readSheetTo = "B";

    // Get the edited range and its sheet
    let editedRange = e.range;
    let editedSheetName = editedRange.getSheet().getName();

    let responseData = {
        message: "Empty"
    };

    if (editedSheetName === sheetENName || editedSheetName === sheetVIName) {
        // Convert the edited range into a Range object
        let editedRangeAsRange = editedRange;

        // Convert the specified read range into a Range object
        let readRange = e.source.getSheetByName(editedSheetName).getRange(readSheetFrom + ":" + readSheetTo);

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
        let responseData = {
            message: "Edited sheet is not 'EN' or 'VI'."
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
    const url = 'https://dev-api.joininflow.io/api/v1/admin/settings/seo/translations';
    const token = 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwidHoiOiJBc2lhL1NhaWdvbiIsImNpZCI6IiIsImN0eXBlIjoiYnV5ZXIiLCJhdWQiOiJzdXBlcl9hZG1pbiIsImlzcyI6ImNnNWdwM3VqbnZ2YjBrbnAzOTAwIiwic3ViIjoic3VwZXJfYWRtaW4ifQ.IOez2kQOF9nlBGzbard-eVT6yrjT7Ex7OgtTj8sLmdE';

    const payload = {
        "domain": "website"
    };

    const headers = {
        'Authorization': token,
        'Content-Type': 'application/json',
    };

    const options = {
        method: 'patch',
        headers: headers,
        payload: JSON.stringify(payload),
    };

    const response = UrlFetchApp.fetch(url, options);
    const content = response.getContentText();

    return ContentService.createTextOutput(content).setMimeType(ContentService.MimeType.JSON);
}

function patchingProd() {
    const url = 'https://api.joininflow.io/api/v1/admin/settings/seo/translations';
    const token = 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNqY3ZtaDRwcmY4dDA5OHJiZnYwIiwidHoiOiJBc2lhL1NhaWdvbiIsImNpZCI6IiIsImN0eXBlIjoiYnV5ZXIiLCJhdWQiOiJzdXBlcl9hZG1pbiIsImlzcyI6ImNqY3Z0dDRwcmY4dDA5OHJiZnZnIiwic3ViIjoic3VwZXJfYWRtaW4ifQ.8E2Teoqy5PgFlyh5MllYBZD3QC6OIKNr4mc1Xt6ulYY';

    const payload = {
        "domain": "website"
    };

    const headers = {
        'Authorization': token,
        'Content-Type': 'application/json',
    };

    const options = {
        method: 'patch',
        headers: headers,
        payload: JSON.stringify(payload),
    };

    const response = UrlFetchApp.fetch(url, options);
    const content = response.getContentText();

    return ContentService.createTextOutput(content).setMimeType(ContentService.MimeType.JSON);
}
