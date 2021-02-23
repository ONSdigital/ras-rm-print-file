# ras-rm-print-file

This service is design to replace functionality in the [action-exporter](https://github.com/ONSdigital/rm-actionexporter-service)

Built in Go as a micro-service this application is takes the entire print file contents as in an json input. It then
applies a Go template to transform the JSON input into a csv file. This CSV file is then upload to both GCS and SFTP.

## Payload:
An example of the JSON payload expected by this service can be found [here](example.json)

## Output
An example of the print file output can be found [here](example.csv)

## Service Design
Please refer to the [design](https://github.com/ONSdigital/ras-rm-documentation/tree/main/service-design/ras-rm-print-file)