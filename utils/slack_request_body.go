package utils

import "github.com/TRON-US/chaos/network/slack"

func ErrorRequestBody(message, title, level string) slack.RequestBody {
	return slack.RequestBody{
		Attachments: []slack.Attachment{
			{
				Color: "danger",
				Fields: []slack.Field{
					{
						Title: title,
						Value: message,
						Short: false,
					},
					{
						Title: "Server IP",
						Value: GetLocalIpAddress(),
						Short: false,
					},
					{
						Title: "Priority",
						Value: level,
						Short: false,
					},
				},
			},
		},
	}
}

func InfoRequestReportBody(message, title string) slack.RequestBody {
	return slack.RequestBody{
		Attachments: []slack.Attachment{
			{
				Color: "good",
				Fields: []slack.Field{
					{
						Title: title,
						Value: message,
						Short: false,
					},
					{
						Title: "Server IP",
						Value: GetLocalIpAddress(),
						Short: false,
					},
				},
			},
		},
	}
}
