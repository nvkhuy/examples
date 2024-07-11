# Build

```
docker build --platform linux/amd64 -t classiy-docker-image:test .
```

# Run

```
docker run --platform linux/amd64 -p 9000:8080 classiy-docker-image:test
```

# Deploy

```
task tf:apply ENV=dev -- -target=module.lambda_classify
```

# Watch Log

```
task app:lambda:watch ENV=dev SERVICE=classify  
```

# Test Local

```
curl --location 'http://localhost:9000/2015-03-31/functions/function/invocations' \
--header 'Content-Type: application/json' \
--data '{
    "body": {
        "image": "https://lucky.a.bigcontent.io/v1/static/DT-BOOT-ALLSIZES",
        "size": 640,
        "confidence": 0.3,
        "overlap": 0
    }
}'
```

# Inflow CURL Sample

```curl
curl --location 'https://dev-classify.joininflow.io/classify' \
--header 'Content-Type: application/json' \
--data '{
    "image_url": "https://lucky.a.bigcontent.io/v1/static/DT-BOOT-ALLSIZES",
    "size": 640,
    "conf_thres": 0.2,
    "iou_thres": 0.1
}'
```

# Roboflow CURL Sample

```curl
curl --location --request POST 'https://detect.roboflow.com/fashion-product-matching-wxvi2/1?api_key=FVziBTEtV3jNZTNOXAZJ&confidence=40&overlap=30&format=json&image=https%3A%2F%2Flucky.a.bigcontent.io%2Fv1%2Fstatic%2FDT-BOOT-ALLSIZES' \
--header 'accept: */*' \
--header 'accept-language: en' \
--header 'content-length: 0' \
--header 'cookie: _gcl_au=1.1.2035857451.1712546726; cookie_utms={%22host%22:%22roboflow.com%22%2C%22path%22:%22/%22%2C%22referrer%22:%22https://www.google.com/%22}; ajs_anonymous_id=439e0cd9-fe40-4d4c-9175-1a75e125ed33; amplitude_idundefinedroboflow.com=eyJvcHRPdXQiOmZhbHNlLCJzZXNzaW9uSWQiOm51bGwsImxhc3RFdmVudFRpbWUiOm51bGwsImV2ZW50SWQiOjAsImlkZW50aWZ5SWQiOjAsInNlcXVlbmNlTnVtYmVyIjowfQ==; _fbp=fb.1.1712546728153.1861297625; ajs_user_id=oN0VbCGmrOdNPfrX82wLMRJHBLh1; _cioid=oN0VbCGmrOdNPfrX82wLMRJHBLh1; __session=eyJhbGciOiJSUzI1NiIsImtpZCI6Il9PQzZaZyJ9.eyJpc3MiOiJodHRwczovL3Nlc3Npb24uZmlyZWJhc2UuZ29vZ2xlLmNvbS9yb2JvZmxvdy1wbGF0Zm9ybSIsIm5hbWUiOiJIdXkgTmd1eWVuIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0w3dGdacHd5VzVRMzktZVVPZnc0cGNFb21OYVliLWdIUEZWSjNFNmVSSVx1MDAzZHM5Ni1jIiwid29ya3NwYWNlcyI6eyJvTjBWYkNHbXJPZE5QZnJYODJ3TE1SSkhCTGgxIjoib3duZXIiLCJBN1pZRlAybTJFWjg2UFpCM0w4UzNRMkgzY2kyIjoib3duZXIiLCJOWU5qUjN5TTZNVFNrZkxjTHBXcVljd0Z6MWwyIjoib3duZXIiLCJWelgyaXI1UjAwbHE5dGpub1dEZCI6Im93bmVyIiwiVlVoZ2Fxd2xGQU90VnpSYTA3cHQiOiJvd25lciIsIjZvUndvU1FObVdhd2NlZGo1TVVHIjoib3duZXIiLCJtc2FSeDEzNnROUzhnQks4QzByS1NDNXIycWIyIjoib3duZXIifSwiYXVkIjoicm9ib2Zsb3ctcGxhdGZvcm0iLCJhdXRoX3RpbWUiOjE3MTI1NDY3NTUsInVzZXJfaWQiOiJvTjBWYkNHbXJPZE5QZnJYODJ3TE1SSkhCTGgxIiwic3ViIjoib04wVmJDR21yT2ROUGZyWDgyd0xNUkpIQkxoMSIsImlhdCI6MTcxMjU0NzE0NSwiZXhwIjoxNzEyOTc5MTQ1LCJlbWFpbCI6Imh1eW5ndXllbkBqb2luaW5mbG93LmlvIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZ29vZ2xlLmNvbSI6WyIxMTEyMDExOTcxNTgxNzI5MTg1NjkiXSwiZW1haWwiOlsiaHV5bmd1eWVuQGpvaW5pbmZsb3cuaW8iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.lxgQZsFYwYll8VyoPvEuTdR286Lpfy7EAIkr5uKpTTLbyUBUDHH51pBZLmNgCdgbsOF4wRPwBNtZMcpVuCiGY2-ACg0jm9CzWxXOnp23sAHNpvDqGA2auRCWid7CCs_ZxaUhrpIRMOuU2b2_xdHJQdpfxwR9Dj0qbRIRHRj9QuOcD8RS8mmz0ychTI399DjVV-ZmEEAs4HaR-6f9elHOUHPZGR7aDm3QBK8McMLKsmMKp2PY4aJ_NGxdt2kmcw_8g9aMVhkRwIwaHqGDTGJ87SjMaGyguIvvj8YLvlRi3jOzDLVpr33k9d5bcj0Snt7T0GHj7K8krhUzUw9AnsvP4Q; _gid=GA1.2.2081226656.1712547168; _ga=GA1.2.2026003725.1712546725; _ga_7RNES0270G=GS1.1.1712547168.1.1.1712547176.0.0.0; amplitude_id_11ee28f1673d40b5f704a83b880a5ddbroboflow.com=eyJkZXZpY2VJZCI6IjlmMWNiYWZiLWFhNTUtNDQzZC1hZTlhLThhNWZmMTZlZWY5M1IiLCJ1c2VySWQiOiJvTjBWYkNHbXJPZE5QZnJYODJ3TE1SSkhCTGgxIiwib3B0T3V0IjpmYWxzZSwic2Vzc2lvbklkIjoxNzEyNTQ2NzI3NzYwLCJsYXN0RXZlbnRUaW1lIjoxNzEyNTQ3MTkwOTA3LCJldmVudElkIjo3LCJpZGVudGlmeUlkIjo3LCJzZXF1ZW5jZU51bWJlciI6MTR9; _ga_SEKT4K1EWR=GS1.1.1712546725.1.1.1712547210.0.0.0; crisp-client%2Fsession%2Fd5d3c29f-9108-4cd7-8296-580b989bc9bc=session_e3e75f98-f47b-4df0-adf0-b0e42512ea70' \
--header 'origin: https://detect.roboflow.com' \
--header 'referer: https://detect.roboflow.com/?model=fashion-product-matching-wxvi2&version=1&api_key=FVziBTEtV3jNZTNOXAZJ' \
--header 'sec-ch-ua: "Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"' \
--header 'sec-ch-ua-mobile: ?0' \
--header 'sec-ch-ua-platform: "macOS"' \
--header 'sec-fetch-dest: empty' \
--header 'sec-fetch-mode: cors' \
--header 'sec-fetch-site: same-origin' \
--header 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36' \
--header 'x-requested-with: XMLHttpRequest'
```

# Response Sample

```
{
    "time": 0.12093428700018194,
    "image": {
        "width": 640,
        "height": 842
    },
    "predictions": [
        {
            "x": 318.5,
            "y": 540.5,
            "width": 301,
            "height": 489,
            "confidence": 0.9080273509025574,
            "class": "women-denim",
            "class_id": 9,
            "detection_id": "d45feed6-4bb0-4581-8321-765f8e04f885"
        }
    ]
}
```