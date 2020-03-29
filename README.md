# Clockify

[Clockify](https://clockify.me) is a great tool for time tracking.   
This is just a simple tool for easy and fast reporting, straight from the terminal. 

### Features
- [x] Get Time Total for all projects aggregated by user
- [x] Export in json format
- [ ] Filter by project

### Prerequisites
- Copy `sample.env` to `.env` and add you **API_KEY**

### Examples

##### Normal
```console
$ ./clockify
Workspace: Stefanos Chrs's workspace
	Project: Blog
		User: Stefanos Chrs
			Total: 8h37m27s
```
  
```console
$ ./clockify  -v
Workspace: Stefanos Chrs's workspace
	Project: Blog
		User: Stefanos Chrs
			Total: 8h37m27s
			Entry: PT57M53S - Theme
			Entry: PT2H6M55S - Add video, optimize images
			Entry: PT1H8M39S - Theme
			Entry: PT1H2M31S - CMS
			Entry: PT3H21M29S - Migration
```
  
##### Output: JSON
```console
$ ./clockify -json
[
  {
    "id": "<workspaceId>",
    "name": "Stefanos Chrs's workspace",
    "projects": [
      {
        "id": "<projectId>",
        "name": "Blog",
        "clientId": "<clientId>",
        "clientName": "<clientName>",
        "users": [
          {
            "id": "<userId>",
            "name": "Stefanos Chrs",
            "email": "<userEmail>",
            "profilePicture": "https://s3.eu-central-1.amazonaws.com/clockify-dev/2019-12-02T02%3A25%3A48.029Zme.jpeg",
            "totalTime": 31047000000000,
            "timeEntries": []
          }
        ]
      }
    ]
  }
]
```
  
```console
$ ./clockify -json -v
[
  {
    "id": "<workspaceId>",
    "name": "Stefanos Chrs's workspace",
    "projects": [
      {
        "id": "<projectId>",
        "name": "Blog",
        "clientId": "<clientId>",
        "clientName": "<clientName>",
        "users": [
          {
            "id": "<userId>",
            "name": "Stefanos Chrs",
            "email": "<userEmail>",
            "profilePicture": "https://s3.eu-central-1.amazonaws.com/clockify-dev/2019-12-02T02%3A25%3A48.029Zme.jpeg",
            "totalTime": 31047000000000,
            "timeEntries": [
              {
                "id": "<taskId>",
                "description": "Theme",
                "TimeInterval": {
                  "start": "2020-01-07T08:15:00Z",
                  "end": "2020-01-07T09:12:53Z",
                  "duration": "PT57M53S"
                }
              },
              {
                "id": "<taskId>",
                "description": "Add video, optimize images",
                "TimeInterval": {
                  "start": "2020-01-07T06:00:00Z",
                  "end": "2020-01-07T08:06:55Z",
                  "duration": "PT2H6M55S"
                }
              },
              {
                "id": "<taskId>",
                "description": "Theme",
                "TimeInterval": {
                  "start": "2020-01-06T22:21:05Z",
                  "end": "2020-01-06T23:29:44Z",
                  "duration": "PT1H8M39S"
                }
              },
              {
                "id": "<taskId>",
                "description": "CMS",
                "TimeInterval": {
                  "start": "2020-01-06T13:30:00Z",
                  "end": "2020-01-06T14:32:31Z",
                  "duration": "PT1H2M31S"
                }
              },
              {
                "id": "<taskId>",
                "description": "Migration",
                "TimeInterval": {
                  "start": "2020-01-05T05:38:31Z",
                  "end": "2020-01-05T09:00:00Z",
                  "duration": "PT3H21M29S"
                }
              }
            ]
          }
        ]
      }
    ]
  }
]
```
