package push

import (
	"encoding/json"
	"strconv"

	"github.com/lfq7413/tomato/livequery/pubsub"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const unsupportedBadgeKey = "unsupported"

func groupByBadge(installations types.S) map[string]types.S {
	result := map[string]types.S{}
	for _, v := range installations {
		if installation := utils.M(v); installation != nil {
			badge := "0"
			if f, ok := installation["badge"].(float64); ok {
				badge = strconv.Itoa(int(f))
			}
			if utils.S(installation["deviceType"]) != "ios" {
				badge = unsupportedBadgeKey
			}

			list := result[badge]
			if list == nil {
				list = types.S{}
			}
			list = append(list, installation)
			result[badge] = list
		}
	}

	return result
}

type pushWorker struct {
	subscriber pubsub.Subscriber
	adapter    pushAdapter
	channel    string
}

func newPushWorker(adapter pushAdapter, channel string) *pushWorker {
	if channel == "" {
		channel = pushChannel
	}
	subscriber := CreateSubscriber()
	worker := &pushWorker{
		subscriber: subscriber,
		adapter:    adapter,
		channel:    channel,
	}

	subscriber.Subscribe(channel)
	subscriber.On("message", func(args ...string) {
		if len(args) < 2 {
			return
		}
		messageStr := args[1]
		var workItem types.M
		err := json.Unmarshal([]byte(messageStr), &workItem)
		if err != nil {
			return
		}
		worker.run(workItem)
	})

	return worker
}

func (p *pushWorker) unsubscribe() {
	p.subscriber.Unsubscribe(p.channel)
}

func (p *pushWorker) run(workItem types.M) error {
	body := utils.M(workItem["body"])
	query := utils.M(workItem["query"])
	status := utils.M(workItem["pushStatus"])

	auth := rest.Master()
	where := utils.M(query["where"])
	delete(query, "where")

	response, err := rest.Find(auth, "_Installation", where, query, nil)
	if err != nil {
		return err
	}
	if utils.HasResults(response) == false {
		return nil
	}
	results := utils.A(response["results"])

	return p.sendToAdapter(body, results, status)
}

func (p *pushWorker) sendToAdapter(body types.M, installations types.S, status types.M) error {
	pushStatus := newPushStatus(utils.S(status["objectId"]))

	if isPushIncrementing(body) == false {
		results := p.adapter.send(body, installations, pushStatus.objectID)
		return pushStatus.trackSent(results)
	}

	badgeInstallationsMap := groupByBadge(installations)

	for badge, ins := range badgeInstallationsMap {
		payload := utils.CopyMapM(body)
		data := utils.M(payload["data"])
		if data == nil {
			data = types.M{}
		}

		if badge == unsupportedBadgeKey {
			delete(data, "badge")
		} else {
			b, err := strconv.Atoi(badge)
			if err != nil {
				continue
			}
			data["badge"] = b
		}

		payload["data"] = data

		err := p.sendToAdapter(payload, ins, types.M{"objectId": pushStatus.objectID})
		if err != nil {
			return err
		}
	}

	return nil
}
