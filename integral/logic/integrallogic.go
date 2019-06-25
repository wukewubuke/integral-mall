package logic

import (
	"context"
	"github.com/streadway/amqp"
	"github.com/yakaa/log4g"
	"my-integral-mall/integral/model"
	"my-integral-mall/integral/protos"
)

type (
	IntegralLogic struct {
		dialHost      string
		queueName     string //队列名称
		rabbitMqConn  *amqp.Connection
		integralModel *model.IntegralModel
		channel     *amqp.Channel
	}
)

func NewIntegralLogic(dialHost, queueName string, integralModel *model.IntegralModel) (*IntegralLogic, error) {
	l := &IntegralLogic{
		dialHost:      dialHost,
		queueName:     queueName,
		integralModel: integralModel,
	}

	if err := l.createDial(); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *IntegralLogic) createDial() error {
	conn, err := amqp.Dial(l.dialHost)
	if err != nil {
		return err
	}
	l.rabbitMqConn = conn
	l.channel, err = l.rabbitMqConn.Channel()
	if err != nil {
		return err
	}
	return nil
}

//关闭rabbit mq链接
func (l *IntegralLogic) CloseRabbitMqConn() {
	if err := l.rabbitMqConn.Close(); err != nil {
		log4g.ErrorFormat("CloseRabbitMqConn error :%+v", err)
	}
	if l.channel != nil {
		if err := l.channel.Close(); err != nil {
			log4g.ErrorFormat("close rabbit mq consume channel error:%+v", err)
		}
	}

}

//向rabbit mq发送消息 (生产消息)
func (l *IntegralLogic) PushMessage(message string) {


	q, err := l.QueueDeclare(l.channel)
	if err != nil {
		log4g.ErrorFormat("channel QueueDeclare error :%+v", err)
		return
	}

	err = l.channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body:        []byte(message),
		})

	if err != nil {
		log4g.ErrorFormat("channel publish error:%+v", err)
	}

}

//消费消息
func (l *IntegralLogic) ConsumeMessage() {
	q, err := l.QueueDeclare(l.channel)
	if err != nil {
		log4g.ErrorFormat("consume Message err %+v", err)
		return
	}

	msgs, err := l.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	go func() {
		for d := range msgs {
			msg := string(d.Body)
			//如果执行失败，把sql语句在放入队列中
			if err := l.integralModel.ExecSql(msg); err != nil {
				l.PushMessage(msg)
			} else {
				log4g.InfoFormat("Consume message %s success!!!", msg)
			}

		}
	}()

}

func (l *IntegralLogic) QueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(l.queueName,
		false, false,
		false, false, nil)

}

//添加积分
func (l *IntegralLogic) AddIntegral(_ context.Context, r *protos.AddIntegralRequest) (*protos.IntegralResponse, error) {
	l.PushMessage(l.integralModel.InsertIntegralSql(r.UserId, r.Integral))
	resp := &protos.IntegralResponse{
		UserId:   r.UserId,
		Integral: r.Integral,
	}
	return resp, nil
}

//消费积分
func (l *IntegralLogic) ConsumerIntegral(_ context.Context, r *protos.ConsumerIntegralRequest) (*protos.IntegralResponse, error) {
	l.PushMessage(l.integralModel.UpdateIntegralByUserIdSql(r.UserId, r.ConsumerIntegral))
	return new(protos.IntegralResponse), nil
}


func (l *IntegralLogic) FindOneByUserId(_ context.Context,r *protos.FindOneByUserIdRequest) (*protos.IntegralResponse, error){
	integral, err := l.integralModel.FindByUserId(r.UserId)
	if err != nil {
		return nil, err
	}

	return &protos.IntegralResponse{
		UserId: integral.UserId,
		Integral: integral.Integral,
	}, nil
}
