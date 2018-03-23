package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

var args struct {
	AmpqUrl string
	Queue   string
	Body    string
	Num     int
}

func parseCommandLine() {
	flag.StringVar(&args.AmpqUrl, "url", "amqp://admin:admin@localhost:5672/", "url for amqp")
	flag.StringVar(&args.Queue, "queue", "", "queue / routing key")
	flag.StringVar(&args.Body, "body", "", "JSON message body")
	flag.IntVar(&args.Num, "num", 1, "number of times to be sent")
	flag.Parse()
	if args.Queue == "" {
		fmt.Println("Missing Queue Name")
		os.Exit(1)
	}
	if args.Body == "" {
		args.Body = `"Hello, world!"`
	}
}

func main() {

	parseCommandLine()
	conn, err := amqp.Dial(args.AmpqUrl)
	if err != nil {
		fmt.Println("Error connecting", err)
		os.Exit(1)
	}
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println("Error getting channel", err)
	}
	channel.Confirm(false)
	confirmations := make(chan amqp.Confirmation, 100)
	end := make(chan interface{})
	channel.NotifyPublish(confirmations)
	var last uint64 = 0
	var done = false
	go func() {
		for conf := range confirmations {
			fmt.Printf("Confirmed:%d\n", conf.DeliveryTag)
			if last < conf.DeliveryTag {
				last = conf.DeliveryTag
			}
			if last == uint64(args.Num) {
				done = true
				break
			}
		}
		if done {
			fmt.Println("All messages sent")
		} else {
			fmt.Println("Confirmation channel closed")
		}
		close(end)
	}()

	_, err = channel.QueueDeclare(
		args.Queue,
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		fmt.Println("Error declaring queue", err)
	}
	for i := 0; i < args.Num; i++ {
		err = channel.Publish(
			"",
			args.Queue,
			true,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(args.Body),
			},
		)
		if err != nil {
			fmt.Println("Error publishing message", err)
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println("")
	<-end
	channel.Close()
	conn.Close()
	fmt.Println("Done!", last)
}
