package hook

import (
	"context"
	"errors"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	log "github.com/sirupsen/logrus"
)

var (
	instance     Queue
	instanceOnce sync.Once
	instanceErr  error
)

const eventsQueueName = string(messaging.QueueNameEvents)

// This implementation is to make the semantics clear that this is the lifecycle of a hook message:
// eventHookMessage --> added to queue --> received --> executed handler (events) --> publish the message using Publisher.Publish()
// queue for adding hook messages to be published
type Queue interface {
	// adds a message to the hooks queue
	Add(ctx context.Context, message *events.EventHookMessage) error
}

type queue struct {
	queueName string
	sender    messaging.MessageSender
}

// Add implements Queue
func (q *queue) Add(ctx context.Context, message *events.EventHookMessage) error {
	results, err := q.sender.Send(ctx, q.queueName, message)
	if err != nil {
		log.Error("Error attempting toadd event message to queue [%s]: %v", q.queueName, err)
		return err
	} else {
		log.Debug("EventHook message sent successfully")
	}
	if len(results) > 0 {
		for _, result := range results {
			if result.Error != nil {
				log.Error("Error sending event message: %v", result.Error)
				return result.Error
			}
		}
	}
	return nil
}

// enqueues a message to the event hooks service
func Add(message *events.EventHookMessage) error {
	if instance == nil {
		return errors.New("hook queue not configured. call Configure() first")
	}
	ctx := context.TODO()

	if message != nil {
		if message.Id == uuid.Nil {
			message.Id = uuid.New()
		}
	}

	return instance.Add(ctx, message)
}

func Configure(appConfig *config.AppConfig) error {
	instanceOnce.Do(func() {
		credential, err := getAzureCredential()
		if err != nil {
			log.Error("Error creating Azure credential for hook.Queue: %v", err)
		}

		sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
			SubscriptionId:          appConfig.Azure.SubscriptionId,
			Location:                appConfig.Azure.Location,
			ResourceGroupName:       appConfig.Azure.ResourceGroupName,
			FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
		})
		if err != nil {
			log.Error("Error creating message sender for hook.Queue: %v", err)
			instanceErr = err
			return
		}

		instance = NewEventHookQueue(sender)
	})
	return instanceErr
}

func NewEventHookQueue(messageSender messaging.MessageSender) Queue {
	return &queue{
		queueName: eventsQueueName,
		sender:    messageSender,
	}
}

func getAzureCredential() (azcore.TokenCredential, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	return credential, err
}