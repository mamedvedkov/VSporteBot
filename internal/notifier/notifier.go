package notifier

import "fmt"

type Sender interface {
	Send(id int64, msg string)
}

type NotifyData struct {
	Id     int64
	Entity string
	Inn    string
	Sum    string
	Month  string
}

type Lister interface {
	ListIds() []NotifyData
}

func Notify(sender Sender, lister Lister) func() {
	go fetchAndSend(sender, lister)

	return func() {
		fetchAndSend(sender, lister)
	}
}

const notification = "От вас необходимо выставить чек.\nОплата в размере %s за %s поступила от\nЮрлицо:\t%s\nИнн:\t%s"

func fetchAndSend(sender Sender, lister Lister) {
	infos := lister.ListIds()

	for idx := range infos {
		sender.Send(infos[idx].Id,
			fmt.Sprintf(
				notification,
				infos[idx].Sum,
				infos[idx].Month,
				infos[idx].Entity,
				infos[idx].Inn,
			),
		)
	}
}
