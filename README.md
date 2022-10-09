# Contest

Mackenzie Math Club's contest permission form system. Written in Go. 

Install dependencies
```
go mod download
```

Build
```
go build
```

If building for Windows, make sure to add the following flags to the build command: `-ldflags -H=windowsgui` 

Despite the fact that Fyne can build to nearly all distributions, only Windows, MacOS, and Linux are supported. 


### Usage

1. Open the executable (do not move the executable out of its folder as it depends on the assets folder)
2. Click on the "Open File" button. Make sure you have a compatible CSV file. 

    A compatible CSV must have the following fields, in that specific order: 

    `email,firstname,lastname,p1teacher,p2teacher`

    The presence of `p2teacher` is optional. However, if there are other additional columns, make sure the `p2teacher` column is empty if it does not contain the Period 2 teacher's name. 

    It is recommended that you verify the table's data before and after importing, as the program does not verify any data itself. 

3. (Optional) Clear the cache. This wipes all existing PDFs generated in the cache folder.
4. Fill in all settings. 

    - Contest Name: the name of the contest as it appears on the permission form
    - Contest Date: the date of the contest as it appears on the permission form
    - Email From: the sending email address. it is recommended that you use `noreply@mackenziemathclub.tk` as it has the appropriate SPF/DKIM records in place (to prevent getting marked as spam/fraud). However, any email address will technically work. 
    - Email Name: the name of the sender as it appears on the email. It is recommended to use `Mackenzie Math Club` for consistency.
    - Email Public Key: this is the public api key taken from mailjet
    - Email Private Key: this is the private/secret api key from mailjet
    - Email Subject & Email Body: these fields are self explanatory. However, this program allows you to use some level of templating by replacing some predetermined strings with data. mail.go#L49 shows a complete list.
